package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/englandrecoil/go-avito-shop/internal/auth"
	"github.com/englandrecoil/go-avito-shop/internal/database"
	"github.com/google/uuid"
)

type sendCoinsRequestParams struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func (cfg *apiConfig) handlerSendCoins(w http.ResponseWriter, r *http.Request) {
	transactionInfo := sendCoinsRequestParams{}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find bearer token", err)
		return
	}

	jwtID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate bearer token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&transactionInfo); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}
	if transactionInfo.Amount < 0 {
		respondWithError(w, http.StatusBadRequest, "Wrong amount of coins provided", errors.New("invalid amount parameter"))
	}

	dbUser, err := cfg.db.GetUserByUsername(r.Context(), transactionInfo.ToUser)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find user", err)
		return
	}

	// user can't send coins to himself
	if dbUser.ID == jwtID {
		respondWithError(w, http.StatusBadRequest, "You can't send coins to yourself", err)
		return
	}

	if err = cfg.transferCoins(r.Context(), jwtID, dbUser.ID, transactionInfo.Amount); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't sent coins", err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

func (cfg *apiConfig) transferCoins(ctx context.Context, senderID, receiverID uuid.UUID, amount int) error {
	tx, err := cfg.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	qtx := cfg.db.WithTx(tx)

	err = qtx.DeductBalance(ctx, database.DeductBalanceParams{
		ID:      senderID,
		Balance: int32(amount),
	})
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deducting balance: %w", err)
	}

	err = qtx.AddBalance(ctx, database.AddBalanceParams{
		ID:      receiverID,
		Balance: int32(amount),
	})
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error adding balance: %w", err)
	}

	err = qtx.InsertTransaction(ctx, database.InsertTransactionParams{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     int32(amount),
	})
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error writing transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
