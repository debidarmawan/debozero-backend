package database_test

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	database "github.com/debidarmawan/debozero/database/sqlc"
	"github.com/debidarmawan/debozero/utils"
	"github.com/stretchr/testify/assert"
)

func cleanup() {
	err := testQuery.DeleteAllUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func createRandomUser(t *testing.T) database.User {
	password, err := utils.GenerateHashPassword(utils.RandomString(8))

	if err != nil {
		log.Fatal("Unale to generate hashed password", err)
	}

	arg := database.CreateUserParams{
		Email:    utils.RandomEmail(),
		Password: password,
	}

	user, err := testQuery.CreateUser(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, user)
	assert.Equal(t, user.Email, arg.Email)
	assert.Equal(t, user.Password, arg.Password)
	assert.WithinDuration(t, user.CreatedAt, time.Now(), 2*time.Second)
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 2*time.Second)

	return user
}

func TestCreateUser(t *testing.T) {
	defer cleanup()

	user1 := createRandomUser(t)
	user2, err := testQuery.CreateUser(context.Background(), database.CreateUserParams{
		Email:    user1.Email,
		Password: user1.Password,
	})
	assert.Error(t, err)
	assert.Empty(t, user2)
}

func TestUpdateUser(t *testing.T) {
	defer cleanup()

	user := createRandomUser(t)
	newPassword, err := utils.GenerateHashPassword(utils.RandomString(8))
	if err != nil {
		log.Fatal("Unale to generate hashed password", err)
	}

	arg := database.UpdateUserPasswordParams{
		Password:  newPassword,
		ID:        user.ID,
		UpdatedAt: time.Now(),
	}

	newUser, err := testQuery.UpdateUserPassword(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, newUser)
	assert.Equal(t, newUser.Password, arg.Password)
	assert.Equal(t, user.Email, newUser.Email)
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 2*time.Second)
}

func TestGetUserByID(t *testing.T) {
	defer cleanup()

	user := createRandomUser(t)
	getUser, err := testQuery.GetUserByID(context.Background(), user.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, getUser)
	assert.Equal(t, getUser.Password, user.Password)
	assert.Equal(t, user.Email, getUser.Email)
}

func TestGetUserByEmail(t *testing.T) {
	defer cleanup()

	user := createRandomUser(t)
	getUser, err := testQuery.GetUserByEmail(context.Background(), user.Email)

	assert.NoError(t, err)
	assert.NotEmpty(t, getUser)
	assert.Equal(t, getUser.Password, user.Password)
	assert.Equal(t, user.Email, getUser.Email)
}

func TestDeleteUser(t *testing.T) {
	defer cleanup()

	user := createRandomUser(t)
	err := testQuery.DeleteUser(context.Background(), user.ID)

	assert.NoError(t, err)

	getUser, err := testQuery.GetUserByID(context.Background(), user.ID)

	assert.Error(t, err)
	assert.Empty(t, getUser)
}

func TestListUsers(t *testing.T) {
	defer cleanup()

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createRandomUser(t)
		}()
	}

	wg.Wait()

	arg := database.GetUserListsParams{
		Offset: 0,
		Limit:  10,
	}

	users, err := testQuery.GetUserLists(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, users)
	assert.Equal(t, len(users), 10)
}
