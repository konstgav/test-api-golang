# test-api-golang

Xsolla Backend School 2021. Test API for e-commerce game developer's system.

## TODO

- Рефакторинг: выделить сущности контроллер, сервис, модель для GraphQL и gRPC
- Добавить обработку ошибок в response
- Добавить авторизацию
- Добавить тесты

## Краткое описание

**Задача**: реализация системы управления товарами для площадки электронной коммерции (продажа таких товаров, как игры, мерч, виртуальная валюта и др.). Реализовать методы API для операций CRUD по управлению товарами. Товар определяется уникальным идентификатором и обязательно должен иметь [SKU](https://ru.wikipedia.org/wiki/SKU), имя, тип, стоимость.

REST API сервер доступен по ссылке `http://localhost:8080/graphql`

1. **Создание товара**. Метод генерирует и возвращает уникальный идентификатор товара.

```bash
POST /product
```

В теле запроса должны быть перечислены поля объекта в формате JSON без дополнительного заворачивания.

2. **Редактирование товара**. Метод изменяет все данные о товаре по его идентификатору.

```bash
PUT /product/:id
```

В теле запроса должны быть перечислены все поля объекта в формате JSON.

3. **Удаление товара по его идентификатору**.

```bash
DELETE /product/:id
```

4. **Получение информации о товаре по его идентификатору**.

```bash
GET /product/:id
```

5. **Получение каталога товаров**. Метод возвращает список всех добавленных товаров по частям. В запросе нужно указать номер запрашиваемой страницы `page` и максимальное количество записей на одну страницу `limit_per_page`.

```bash
GET /product?page=...&limit_per_page=...
```

Простой случай: в случае успеха сервер возвращает 200 OK с массивом объектов в формате JSON в теле ответа (т.е. ответ начинается с [ и заканчивается ]).

Если массив получился пустой, всё равно вовзращается 200 OK с пустым масивом [] в теле ответа.

## Развертывание и тестирование

Для развертывания приложения с помощью `docker-compose` необходимо выполнить команды:  

```bash
git clone https://github.com/konstgav/test-api-golang.git
cd test-api-golang
docker-compose up 
```

Генерация графа зависимостей пакетов

```bash
~/go/bin/godepgraph -novendor -s  -p github.com,go.mongodb.org,golang.org,google.golang.org . | dot -Tpng -o godepgraph.png
```

![Зависимость пакетов приложения](godepgraph.png?raw=true "Dependencies graph")

## GraphQL

Запускается graphql-сервер, который позвояет реализовать базовые CRUD операции. Среда для тестирования GraphiQL доступна по ссылке `http://localhost:5000/graphql`.

1. `list: [Product]` - возвращает список товаров.

2. `product(id: Int): Product` - возвращает товар по идентификатору.

3. `create(id: Int!name: String!sku: String!type: String!price: Int!): Product` - создает новый товар.

4. `delete(id: Int!): Product` - удаляет товар по идентификатору.

5. `update(id: Int!name: String!sku: String!type: String!price: Int!): Product` - обновляет информацию о товаре.

Пример использования:

```(bash)
curl http://localhost:5000/graphql?query=%7Blist%7Bname%7D%7D
```

## gRPC

Добавлен gRPC сервис, позволяющий отправить письмо при совершении пользователем покупки. В API добавлен endpoint `POST http://localhost:8080/sendmail`, сервер получает идентификатор товара (`_id`) и email пользователя (`email`) в теле запроса. API-сервер обращается в gRPC-сервис, чтобы отправить письмо.

Для генерации кода клиента и сервера используется команда:

```(bash)
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative mailserver/mailer.proto
```

Пример использования:

```(bash)
curl -X POST -H "Content-Type: application/json" -d '{"_id": 2,"email":"gavrilov.k@psu.ru"}' "http://localhost:8080/sendmail" 
```

Для работы приложения необходимо создать файл `.env` в корневом каталоге, содержащий значения переменных окружения:

```
MAILER_REMOTE_HOST="smtp-server:port"
MAILER_FROM="user@example.com"
MAILER_PASSWORD="passwd"
```

## OAuth 2.0 Xsolla Login

Доступ для создания, удаления и редактирования товаров предоставляется только авторизированным клиентам. Авторизация backend-клиентов производится при помощи API сервиса [Xsolla Login](https://developers.xsolla.com/doc/login/features/connecting-oauth2/). Для тестирования необходимо выполнить следующие шаги:

1. Создать publisher-аккаунт на сервисе [Xsolla Login](https://xsolla.com/products/login).

2. Настроить проект в publisher-аккаунте в соотвествии с [описанием](https://developers.xsolla.com/doc/login/integration-guide/create-project/).

3. Для получения учетных данных `projectId`, `client_id`, `client_secret`, `secret key` требуется создать [login-проект](https://developers.xsolla.com/doc/login/integration-guide/set-up-login-project/).

4. Зарегистрировать нового пользователя при помощи POST запроса на [этой странице](https://developers.xsolla.com/login-api/auth/jwt/jwt-register-new-user/).

5. На вкладке `Users` login-проекта должен отобразиться новый пользователь. Нужно создать группу пользователей `readers`, зайти в профиль нового пользователя и добавить его в группу `readers`. Это необходимо для предоставления доступа на создание, удаление и редактирование товаров.

6. Локально запусить веб-приложение командой `docker-compose up`.

7. В браузере перейти по ссылке для тестирования POST запроса `http://localhost:8080/test-post`, который добавляет новый товар в приложение. При попытке напрямую без авторизации добавить товар, приложение вернет ошибку `401 Unauthorized`. Вызывается функция, которая обращается к сервису Xsolla Login за JWT-токеном и отправляет запрос, содержащий в заголовке полученный токен.

8. В папке `oauth` необходимо иметь следующие файлы с учетными данными:

    8.1. Файл `xsolla-login-account-credentials.json` содержит `client_id` и `client_secret` для login-проект.

    8.2. Файл `xsolla-login-user-credentials.json` содержит `password` и `username` для пользователя.

    8.3. Файл `secret.pem` содержит секретный ключ для валидации JWT-токена.

## OAuth 2.0 with Google account

Доступ для просмотра одного товара или списка товаров предоставляется только авторизированным клиентам. Авторизация frontend-клиентов производится при помощи [Google API](https://console.cloud.google.com). Не удалось обменть google access token на JWT токен, не реализована валидатция токена. Для тестирования необходимо выполнить следующие шаги:

1. Создать проект в `https://console.cloud.google.com`.

2. Во вкладке `Credentials` проекта сгенерировать `OAuth 2.0 Client IDs`. В качестве `redirect_URL` указать `http://localhost:8080/oauth2callback`.

3. Локально запусить веб-приложение командой `docker-compose up`.

4. В браузере перейти по ссылке для авторизации с помощью google-аккаунта [`http://localhost:8080/authorize`](http://localhost:8080/authorize), потребуется ввести логин/пароль.

5. В браузере перейти по ссылке [`http://localhost:8080/product`](http://localhost:8080/product) для просмотра списка товаров. При попытке напрямую без авторизации просмотреть товары, приложение вернет ошибку `401 Unauthorized`.

6. В папке `oauth` необходимо хранить файл `google-account-credentials.json` с учетными данными `client_id`, `client_secret`, `redirect_URL`, `scopes`.

## RabbitMQ

После обращения к эндпоинту `POST /rabbitmq` в поле `link` json-запроса передается URL-строка c путем к лендингу игры. Эта строка отправляется другому приложению с помощью RabbitMQ на обработку (консьюмеру).  Реализован механизм проверки доступности данного лендинга и переотправки (повторной обработки данного сообщения) запроса в случае, если попытка была неудачной. Мы считаем, что попытка удачна, если при обращении к лендингу возвращается статус код 20*. Неудачной попыткой мы считаем получение любого отличного от 20* статус кода в ответ на запрос.

При неудачной попытке сообщение должно поступать на повторную обработку консьюмеру не сразу, а через заданные `resendTime` секунд. При этом сообщение возвращется обратно в очередь, все сообщения в которой "живут" `messageTTL` секунд. Для проверки доступности URL запускаются несколько воркеров в гоурутинах, количество которых определятся константой `MaxOutstanding`.

`rabbitmq` стартует вместе с приложением командой `docker-compose up`. Для тестирования приложения выполнить команды:

```(bash)
curl -X POST http://localhost:8080/rabbitmq -d '{"link":"https://www.google.com1"}'
curl -X POST http://localhost:8080/rabbitmq -d '{"link":"https://www.google.com"}'
curl -X POST http://localhost:8080/rabbitmq -d '{"link":"https://www.xsolla.com1"}'
```

В терминале приложение сообщает о рабочей (`Correct link`) или нерабочей (`Incorrect link`) URL. Можно отследить, как сообщения возвращаются в очередь спустя заданное время.
Для интерактивного запуска docker-контейнера `rabbitmq` c доступом к [веб-интефейсу](http://localhost:15672) используется команда:

```(bash)
docker run --rm -it -p 15672:15672 -p 5672:5672 rabbitmq:3-management
```
