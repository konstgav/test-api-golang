package mailserver

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/smtp"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	UnimplementedMailerServer
}

type mess struct {
	Email string
	Text  []byte
}

func (s *server) SendPass(ctx context.Context, in *MsgRequest) (*MsgReply, error) {
	log.Println("Server recieved message: ", in.ProductName, in.Email)
	m := mess{Text: []byte("Hello, you have just bought the " + in.ProductName + "!"), Email: in.Email}
	select {
	case queue <- m:
	default:
		return &MsgReply{IsSent: false}, nil
	}
	return &MsgReply{IsSent: true}, nil
}

func getSMTPClient() *smtp.Client {
	host, _, err := net.SplitHostPort(cnf.smtphost)
	if err != nil {
		log.Println("SplitHostPort", err)
		return nil
	}
	tlsconfig := &tls.Config{
		ServerName: host,
	}
	conn, err := net.Dial("tcp", cnf.smtphost)
	if err != nil {
		log.Println("tls.dial", err)
		return nil
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Println("new client", err)
		return nil
	}
	if err = client.StartTLS(tlsconfig); err != nil {
		log.Println("starttls", err)
		return nil
	}
	auth := smtp.PlainAuth("", cnf.from, cnf.pass, host)
	if err = client.Auth(auth); err != nil {
		log.Println("auth", err)
		return nil
	}
	log.Println("Connection to smtp-server is created")
	return client
}

func MessageLoop() {
	Init()
	queue = make(chan mess, 10)
	client := getSMTPClient()
	if client == nil {
		return
	}
	var err error
	defer func() {
		err := client.Quit()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	for m := range queue {
		err = client.Noop()
		if err != nil {
			log.Println("reestablish connection", err)
			return
		}
		if err = client.Mail(cnf.from); err != nil {
			log.Println(err)
			return
		}
		if err = client.Rcpt(m.Email); err != nil {
			log.Println(err)
			return
		}
		writecloser, err := client.Data()
		if err != nil {
			log.Println(err)
			return
		}
		_, err = writecloser.Write(m.Text)
		if err != nil {
			log.Println(err)
			return
		}
		err = writecloser.Close()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Successfully send message")
	}
}

func StartMailer() {
	listener, err := net.Listen("tcp", ":20100")
	if err != nil {
		log.Fatal("failed to listen", err)
	}
	log.Printf("start listening for emails at port %s", ":20100")
	rpcserv := grpc.NewServer()
	RegisterMailerServer(rpcserv, &server{})
	reflection.Register(rpcserv)
	err = rpcserv.Serve(listener)
	if err != nil {
		log.Fatal("failed to serve", err)
	}
}

type conf struct {
	smtphost, from, pass string
}

var cnf conf

var queue chan mess

func Init() {
	if os.Getenv("MAILER_REMOTE_HOST") == "" || os.Getenv("MAILER_FROM") == "" || os.Getenv("MAILER_PASSWORD") == "" {
		panic("Environmental variable do not set")
	}
	cnf = conf{
		os.Getenv("MAILER_REMOTE_HOST"),
		os.Getenv("MAILER_FROM"),
		os.Getenv("MAILER_PASSWORD"),
	}
}
