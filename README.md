# go-avito-shop
**Avito Shop API** - RESTful веб-сервис, написанный на языке Go. Данный сервис представляет из себя внутренний магазин мерча, который позволяет:
- Приобретать товары за монеты
- Передавать монеты другим сотрудникам
- Смотреть историю транзакций (исходящие и входящие переводы) и список купленных товаров

## Инструкция для запуска
1. Склонировать репозиторий с помощью:
``` bash
git clone https://github.com/englandrecoil/go-avito-shop.git
```
2. Перейти в склонированный репозиторий:
``` bash
cd go-avito-shop
```
3. Запустить сервис с Docker:
```bash
docker compose up --build 
```
После успешного запуска сервер будет доступен снаружи как localhost:8080

Подробное описание эндпоинтов можно найти по ссылке [API](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025/schema.json)

## Дополнения
### 1. Добавлен эндпоинт для регистрации новых пользователей
В изначальной документации не было предусмотрено отдельного метода регистрации пользователей. Однако для аутентификации используется JWT, поэтому мне понадобился способ создания нового пользователя. Поэтому был добавлен эндпоинт `/api/reg`. Он также возвращает токен доступа при успешной регистрации, как и эндпоинт `/api/auth` при успешной аутентификации.
Вот его структура:
- Метод: `POST`
- URL: `/api/reg`
- Принимает JSON с `username` и `password`
- Создает нового пользователя и выдает токен доступа
- Возвращает `201 Created` при успешной регистрации

Пример запроса:
```json
{
  "username": "username",
  "password": "password"
}
```
Пример ответа:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI..."
}
```

### 2. Схема БД
Так как сервис использует PostgreSQL для хранения данных, я посчитал, что мне стоит показать структуру схемы. Ниже представлена схема базы данных:
<img src="https://i.ibb.co/VpPp2CfZ/schemadb.jpg">

### 3. Тестирование
Интеграционные тесты и некоторые юнит-тесты находятся в файле `handlers_test.go`. Они используют тестовую БД без моков, операции не изолированы, то есть тесты могут влиять друг на друга. Проверяются сценарии из нескольких шагов(регистрации, авторизации, переводов и покупок). 

Сценарии:
- Покупка мерча
- Передача монет
- Получение информации об истории переводов/купленных предметов

Также были реализованы юнит-тесты (json_test.go, auth_test.go), которые проверяют корректность обработки JSON ответов, хеширования, генерацию и валидацию JWT и т.д.

### 4. Результаты тестов
Ниже приведены проценты покрытия кодом, которые были получены при использовании `go test ./... -cover`. 
<img src="https://i.ibb.co/Nd66Twh4/testinfo.jpg">

Общее покрытие кода составляет: 47.3%
Покрытие модуля аутентификации (internal/auth): 83.9%
Покрытие модуля базы данных (internal/database): 0.0%
Нулевое покрытие для internal/database объясняется тем, что этот модуль был автоматически сгенерирован утилитой sqlc на основе SQL-запросов.
По итогу фактическое покрытие бизнес-логики выше, чем отображается в общем показателе.

### 5. .env файл
Файл .env намеренно был оставлен в общем доступе. Хоть я и понимаю, что в проде .env точно не должен храниться так, однако мне показалось, что это упростит проверку и не придется настраивать переменные среды, а также там не находятся конфиденциальных данных.


