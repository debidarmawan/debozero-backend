package api

import (
	"database/sql"
	"fmt"
	"net/http"

	database "github.com/debidarmawan/debozero/database/sqlc"
	"github.com/debidarmawan/debozero/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	queries *database.Queries
	router  *gin.Engine
	config  *utils.Config
}

var tokenController *utils.JWTToken

func NewServer(envPath string) *Server {
	config, err := utils.LoadConfig(envPath)
	if err != nil {
		panic(fmt.Sprintf("Could not load env config: %v", err))
	}

	conn, err := sql.Open(config.DBDriver, config.DBSourceLive)
	if err != nil {
		panic(fmt.Sprintf("Could not connect to database: %v", err))
	}

	tokenController = utils.NewJWTToken(config)

	q := database.New(conn)

	g := gin.Default()

	g.Use(cors.Default())

	return &Server{
		queries: q,
		router:  g,
		config:  config,
	}
}

func (s *Server) Start(port int) {
	s.router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to DeboZero"})
	})

	User{}.router(s)
	Auth{}.router(s)

	s.router.Run(fmt.Sprintf(":%d", port))
}
