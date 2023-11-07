package api

import (
	"context"
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
