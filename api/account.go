package api

import (
	"context"
	"database/sql"
	"net/http"

	database "github.com/debidarmawan/debozero-backend/database/sqlc"
	"github.com/debidarmawan/debozero-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Account struct {
	server *Server
}

func (a Account) router(server *Server) {
	a.server = server

	serverGroup := server.router.Group("/account", AuthenticatedMiddleware())

	serverGroup.POST("create", a.createAccount)
	serverGroup.GET("", a.getUserAccount)
	serverGroup.POST("transfer", a.transfer)
}

type AccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (a Account) createAccount(c *gin.Context) {
	userID, err := utils.GetActiveUser(c)
	if err != nil {
		return
	}

	acc := new(AccountRequest)

	if err := c.ShouldBindJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := database.CreateAccountParams{
		UserID:   int32(userID),
		Currency: acc.Currency,
		Balance:  0,
	}

	account, err := a.server.queries.CreateAccount(context.Background(), arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "you already have account with this currency"})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (a Account) getUserAccount(c *gin.Context) {
	userID, err := utils.GetActiveUser(c)
	if err != nil {
		return
	}

	accounts, err := a.server.queries.GetAccountByUserID(context.Background(), int32(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

type TransferRequest struct {
	ToAccountID   int32   `json:"to_account_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	FromAccountID int32   `json:"from_account_id" binding:"required"`
}

func (a *Account) transfer(c *gin.Context) {
	userID, err := utils.GetActiveUser(c)
	if err != nil {
		return
	}

	tr := new(TransferRequest)

	if err := c.ShouldBindJSON(&tr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fromAccount, err := a.server.queries.GetAccountByID(context.Background(), int64(tr.FromAccountID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't get account"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if fromAccount.UserID != int32(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't get account"})
		return
	}

	toAccount, err := a.server.queries.GetAccountByID(context.Background(), int64(tr.ToAccountID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "destination account is not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if toAccount.Currency != fromAccount.Currency {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Currency is not match"})
		return
	}

	if fromAccount.Balance < tr.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insuficient balance"})
		return
	}

	txArg := database.CreateTransferParams{
		FromAccountID: tr.FromAccountID,
		ToAccountID:   tr.ToAccountID,
		Amount:        tr.Amount,
	}

	tx, err := a.server.queries.TransferTx(context.Background(), txArg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error is occured"})
		return
	}

	c.JSON(http.StatusCreated, tx)
}
