package main

import (
	"database/sql"
	"net/http"

	"github.com/englandrecoil/go-avito-shop/internal/auth"
)

func (cfg *apiConfig) handlerInfo(w http.ResponseWriter, r *http.Request) {
	type infoRequestParams struct {
		Coins     int `json:"coins"`
		Inventory []struct {
			Type     string `json:"type"`
			Quantity int    `json:"quantity"`
		} `json:"inventory"`
		CoinHistory struct {
			Received []struct {
				FromUser string `json:"fromUser"`
				Amount   int    `json:"amount"`
			} `json:"received"`
			Sent []struct {
				ToUser string `json:"toUser"`
				Amount int    `json:"amount"`
			} `json:"sent"`
		} `json:"coinHistory"`
	}
	infoUser := infoRequestParams{}

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

	dbUser, err := cfg.db.GetUserByID(r.Context(), jwtID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find user", err)
		return
	}

	receivedHistory, err := cfg.db.GetReceivedHistory(r.Context(), jwtID)
	if err != nil {
		if err != sql.ErrNoRows {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get receive history", err)
			return
		}
	}

	sentHistory, err := cfg.db.GetSentHistory(r.Context(), jwtID)
	if err != nil {
		if err != sql.ErrNoRows {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get receive history", err)
			return
		}
	}

	inventory, err := cfg.db.GetInventory(r.Context(), jwtID)
	if err != nil {
		if err != sql.ErrNoRows {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get receive history", err)
			return
		}
	}

	// aggregating data
	infoUser.Coins = int(dbUser.Balance)
	for _, value := range inventory {
		infoUser.Inventory = append(infoUser.Inventory, struct {
			Type     string `json:"type"`
			Quantity int    `json:"quantity"`
		}{
			Type:     value.ItemName,
			Quantity: int(value.Quantity),
		})
	}
	for _, value := range receivedHistory {
		infoUser.CoinHistory.Received = append(infoUser.CoinHistory.Received, struct {
			FromUser string `json:"fromUser"`
			Amount   int    `json:"amount"`
		}{
			FromUser: value.SenderName,
			Amount:   int(value.Received),
		})
	}
	for _, value := range sentHistory {
		infoUser.CoinHistory.Sent = append(infoUser.CoinHistory.Sent, struct {
			ToUser string `json:"toUser"`
			Amount int    `json:"amount"`
		}{
			ToUser: value.ReceiverName,
			Amount: int(value.Sent),
		})
	}
	respondWithJSON(w, http.StatusOK, infoUser)
}
