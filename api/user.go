package api

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	database "github.com/debidarmawan/debozero-backend/database/sqlc"
	"github.com/gin-gonic/gin"
)

type User struct {
	server *Server
}

func (u User) router(server *Server) {
	u.server = server

	serverGroup := server.router.Group("/users", AuthenticatedMiddleware())

	serverGroup.GET("", u.listUsers)
	serverGroup.GET("me", u.getLoggedInUser)
}

type UserParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (u *User) listUsers(c *gin.Context) {
	arg := database.GetUserListsParams{
		Offset: 0,
		Limit:  10,
	}

	users, err := u.server.queries.GetUserLists(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newUsers := []UserResponse{}

	for _, v := range users {
		n := UserResponse{}.toUserResponse(&v)
		newUsers = append(newUsers, *n)
	}

	c.JSON(http.StatusOK, newUsers)
}

func (u *User) getLoggedInUser(c *gin.Context) {
	value, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	userID, ok := value.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unexpected error"})
		return
	}

	user, err := u.server.queries.GetUserByID(context.Background(), userID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{}.toUserResponse(&user))
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u UserResponse) toUserResponse(user *database.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
