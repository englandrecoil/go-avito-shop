package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/englandrecoil/go-avito-shop/internal/auth"
	"github.com/englandrecoil/go-avito-shop/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB
var testQueries *database.Queries

func runMigrations(db *sql.DB) {
	migrationsDir := "sql/schema"
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func TestMain(m *testing.M) {
	_ = godotenv.Load()

	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		log.Fatal("TEST_DB_URL must be set")
	}

	var err error
	testDB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	testQueries = database.New(testDB)

	runMigrations(testDB)

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

func setupTestDB(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(`
		TRUNCATE TABLE users CASCADE;
	`)
	if err != nil {
		t.Fatalf("Failed to clean test DB: %v", err)
	}
}

func TestHandlerRegister(t *testing.T) {
	setupTestDB(t)
	cfg := apiConfig{db: testQueries, secret: "test-secret"}

	user := CredentialsRequestParams{
		Username: "testuser",
		Password: "password123",
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/reg", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	cfg.handlerRegister(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", resp["username"])
	assert.Contains(t, resp, "token")
}

func TestHandlerAuth(t *testing.T) {
	setupTestDB(t)

	cfg := apiConfig{db: testQueries, secret: "test-secret"}

	hashedPassword, _ := auth.HashPassword("password123")
	_, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		Username:       "testuser",
		HashedPassword: hashedPassword,
	})

	login := CredentialsRequestParams{
		Username: "testuser",
		Password: "password123",
	}

	body, _ := json.Marshal(login)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	cfg.handlerAuth(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	responseStruct := authResponseParams{}

	err = json.Unmarshal(rec.Body.Bytes(), &responseStruct)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseStruct.Token)
}

func TestHandlerBuyItem(t *testing.T) {
	setupTestDB(t)

	cfg := apiConfig{db: testQueries, secret: "test-secret"}

	hashedPassword, _ := auth.HashPassword("password123")
	_, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		Username:       "testuser",
		HashedPassword: hashedPassword,
	})
	assert.NoError(t, err)

	login := CredentialsRequestParams{
		Username: "testuser",
		Password: "password123",
	}

	body, _ := json.Marshal(login)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	cfg.handlerAuth(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	responseStruct := authResponseParams{}
	err = json.Unmarshal(rec.Body.Bytes(), &responseStruct)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseStruct.Token)

	req = httptest.NewRequest(http.MethodGet, "/api/buy", nil)
	req.SetPathValue("item", "book")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+responseStruct.Token)
	rec = httptest.NewRecorder()
	cfg.handlerBuyItem(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	dbUser, err := cfg.db.GetUserByUsername(req.Context(), login.Username)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	assert.Equal(t, 950, int(dbUser.Balance))

	items, err := cfg.db.GetInventory(req.Context(), dbUser.ID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	assert.NotEmpty(t, items)
}

func TestHandlerInfo(t *testing.T) {
	cfg := apiConfig{db: testQueries, secret: "test-secret"}

	login := CredentialsRequestParams{
		Username: "testuser",
		Password: "password123",
	}

	body, _ := json.Marshal(login)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	cfg.handlerAuth(rec, req)
	responseStruct := authResponseParams{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseStruct)

	assert.NoError(t, err)
	assert.NotEmpty(t, responseStruct.Token)
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+responseStruct.Token)
	rec = httptest.NewRecorder()

	cfg.handlerInfo(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandlerSendCoins(t *testing.T) {
	setupTestDB(t)

	cfg := apiConfig{db: testQueries, secret: "test-secret", conn: testDB}
	// register user1
	hashedPassword, _ := auth.HashPassword("password123")
	_, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		Username:       "testuser1",
		HashedPassword: hashedPassword,
	})
	assert.NoError(t, err)

	login := CredentialsRequestParams{
		Username: "testuser1",
		Password: "password123",
	}

	body, _ := json.Marshal(login)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	cfg.handlerAuth(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// get access token for user1
	responseStruct1 := authResponseParams{}
	err = json.Unmarshal(rec.Body.Bytes(), &responseStruct1)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseStruct1.Token)

	// reguster user2
	hashedPassword, _ = auth.HashPassword("password123")
	_, err = cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		Username:       "testuser2",
		HashedPassword: hashedPassword,
	})
	assert.NoError(t, err)

	login = CredentialsRequestParams{
		Username: "testuser2",
		Password: "password123",
	}

	body, _ = json.Marshal(login)
	req = httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()

	cfg.handlerAuth(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// get access token for user 2
	responseStruct2 := authResponseParams{}
	err = json.Unmarshal(rec.Body.Bytes(), &responseStruct2)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseStruct2.Token)

	// user1 sending coins to user2
	transactionInfo := sendCoinsRequestParams{
		ToUser: "testuser2",
		Amount: 521,
	}
	body, _ = json.Marshal(transactionInfo)

	req = httptest.NewRequest(http.MethodPost, "/api/sendCoins", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+responseStruct1.Token)
	rec = httptest.NewRecorder()

	cfg.handlerSendCoins(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// check balance for both users
	dbUser1, err := cfg.db.GetUserByUsername(req.Context(), "testuser1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	dbUser2, err := cfg.db.GetUserByUsername(req.Context(), "testuser2")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assert.Equal(t, 479, int(dbUser1.Balance))
	assert.Equal(t, 1521, int(dbUser2.Balance))
}
