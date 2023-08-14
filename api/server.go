package api

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	db "github.com/liquiddev99/dropbyte-backend/db/sqlc"
	"github.com/liquiddev99/dropbyte-backend/token"
	"github.com/liquiddev99/dropbyte-backend/util"
)

type Server struct {
	config util.Config
	db     *db.Queries
	router *gin.Engine
	token  token.Token
}

func NewServer(config util.Config, db *db.Queries) (*Server, error) {
	token, err := token.NewMaker(config.SymmetricKey)
	if err != nil {
		log.Fatal("Cannot create token maker")
	}
	server := &Server{config: config, db: db, token: token}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	corsConf := cors.DefaultConfig()
	corsConf.AllowOrigins = []string{server.config.OriginAllowed}
	corsConf.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}
	corsConf.AllowCredentials = true
	corsConf.AllowHeaders = []string{
		"Content-Type",
		"Authorization",
		"Accept",
		"X-Requested-With",
		"Origin",
		"Access-Control-Request-Headers",
	}

	router.Use(cors.New(corsConf))

	authRoutes := router.Group("/").Use(authMiddleware(server.token))
	authRoutes.Use(cors.New(corsConf))

	router.POST("/upload", server.guestUploadFile)
	router.POST("/signup", server.createUser)
	router.POST("/login", server.loginUser)

	authRoutes.POST("/user/upload", server.userUploadFile)
	authRoutes.GET("/user/files", server.getFiles)
	authRoutes.POST("/user/file/delete", server.deleteFileById)
	authRoutes.GET("/user/file/download", server.downloadFileById)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func responseError(err error) gin.H {
	return gin.H{"error": err.Error()}
}
