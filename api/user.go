package api

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	database "github.com/debidarmawan/debozero-backend/database/sqlc"
	"github.com/debidarmawan/debozero-backend/utils"
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
	serverGroup.PATCH("username", u.updateUsername)
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
	userID, err := utils.GetActiveUser(c)
	if err != nil {
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

type UpdateUsernameType struct {
	Username string `json:"username" binding:"required"`
}

func (u *User) updateUsername(c *gin.Context) {
	userID, err := utils.GetActiveUser(c)
	if err != nil {
		return
	}

	var userInfo UpdateUsernameType

	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := database.UpdateUsernameParams{
		ID: userID,
		Username: sql.NullString{
			String: userInfo.Username,
			Valid:  true,
		},
		UpdatedAt: time.Now(),
	}

	user, err := u.server.queries.UpdateUsername(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{}.toUserResponse(&user))
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	UserName  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u UserResponse) toUserResponse(user *database.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		UserName:  user.Username.String,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
