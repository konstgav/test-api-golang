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

Доступ для создания, удаления и редактирования товаров предоставляется только авторизированным клиентам. Авторизация производится при помощи сервиса [Xsolla Login](https://developers.xsolla.com/doc/login/features/connecting-oauth2/). Для тестирования необходимо выполнить следующие шаги:

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
