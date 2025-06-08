# go-avito-shop
[![ru](https://img.shields.io/badge/lang-ru-blue?style=flat)](https://github.com/englandrecoil/go-avito-shop/blob/master/README.md)

**Avito Shop API** is a RESTful web service written in Go. This internal merchandise store allows users to:
- Purchase items using coins  
- Transfer coins to other employees  
- View transaction history (incoming and outgoing transfers) and a list of purchased items  

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/englandrecoil/go-avito-shop.git
```
2. Navigate into the cloned repository:
```bash
cd go-avito-shop
```
3.Run the service with Docker:
```bash
docker compose up --build 
```
Once started, the server will be available at localhost:8080

You can find the full endpoint documentation here [API](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025/schema.json)

## Additions
### 1. User Registration Endpoint
The original documentation did not provide a registration method. Since JWT is used for authentication, I added an `/api/reg` endpoint to create users. It returns an access token upon successful registration, similar to `/api/auth`.
- Method: `POST`
- URL: `/api/reg`
- Accepts JSON with `username` and `password`
- Creates a new user and returns an access token
- Responds with `201 Created` on success

Example request:
```json
{
  "username": "username",
  "password": "password"
}
```
Example response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI..."
}
```

### 2. Database Schema
Since the service uses PostgreSQL, I decided to include the database structure:
<img src="https://i.ibb.co/VpPp2CfZ/schemadb.jpg">

### 3. Testing
Integration and some unit tests are in handlers_test.go. They use a test database without mocks and are not isolated, so they may affect each other. Covered scenarios include:
- Buying merchandise
- Transferring coins
- Fetching transfer history and purchased items

Also included are unit tests (json_test.go, auth_test.go) that check JSON handling, password hashing, JWT generation/validation, etc.

### 4. Test Results
Code coverage from `go test ./... -cover`:
<img src="https://i.ibb.co/Nd66Twh4/testinfo.jpg">
- Total coverage: 47.3%
- internal/auth module: 83.9%
- internal/database module: 0.0%
The internal/database module is generated with sqlc from SQL queries, hence the 0% coverage. Actual business logic coverage is higher.

### 5. .env File
The .env file is intentionally included in the repo. Although it's not suitable for production, this simplifies review and local setup. It does not contain sensitive data.
