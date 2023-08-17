package api

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	db "github.com/liquiddev99/dropbyte-backend/db/sqlc"
	"github.com/liquiddev99/dropbyte-backend/request"
	"github.com/liquiddev99/dropbyte-backend/token"
	"github.com/liquiddev99/dropbyte-backend/util"
)

type Server struct {
	config         util.Config
	db             *db.Queries
	router         *gin.Engine
	token          token.Token
	b2UploadUrl    string
	b2UrlAuthToken string
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

func (server *Server) scheduledTask() {
	authResponse, err := request.AuthorizeAccount(
		server.config.B2ApplicationKeyId,
		server.config.B2ApplicationKey,
	)
	if err != nil {
		log.Fatal("Failed to authorize b2 account")
		return
	}

	urlResponse, err := request.GetUploadUrl(
		server.config.BucketId,
		authResponse.AuthorizationToken,
	)
	if err != nil {
		log.Fatal("Failed to get upload url b2")
		return
	}
	server.b2UploadUrl = urlResponse.UploadUrl
	server.b2UrlAuthToken = urlResponse.AuthorizationToken
}

func (server *Server) startScheduledTask() {
	interval := 12 * time.Hour
	ticker := time.Tick(interval)

	server.scheduledTask()

	go func() {
		for {
			<-ticker
			server.scheduledTask()
		}
	}()
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
	router.MaxMultipartMemory = 300 << 20

	authRoutes := router.Group("/").Use(authMiddleware(server.token))

	router.POST("/upload", server.guestUploadFile)
	router.POST("/signup", server.createUser)
	router.POST("/login", server.loginUser)

	authRoutes.POST("/user/upload", server.userUploadFile)
	authRoutes.GET("/user/files", server.getFiles)
	authRoutes.POST("/user/file/delete", server.deleteFileById)
	authRoutes.GET("/user/file/download", server.downloadFileById)
	authRoutes.POST("/user/logout", server.logout)
	server.router = router
}

func (server *Server) Start(address string) error {
	server.startScheduledTask()
	return server.router.Run(address)
}

func responseError(err error) gin.H {
	return gin.H{"error": err.Error()}
}
