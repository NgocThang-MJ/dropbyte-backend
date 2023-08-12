package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	db "github.com/liquiddev99/dropbyte-backend/db/sqlc"
	"github.com/liquiddev99/dropbyte-backend/util"
)

type createUserRequest struct {
	Password string `json:"password"  binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email"     binding:"required,email"`
}

type userResponse struct {
	FullName  string    `json:"full_name"    binding:"required"`
	Email     string    `json:"email"        binding:"required,email"`
	Token     string    `json:"access_token"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User, token string) userResponse {
	return userResponse{
		FullName:  user.FullName,
		Email:     user.Email,
		Token:     token,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	arg := db.CreateUserParams{
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.db.CreateUser(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			ctx.JSON(http.StatusBadRequest, responseError(pgErr))
			return
		}
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}
	token, err := server.token.CreateToken(user.ID, server.config.AccessTokenDuration)

	ctx.SetCookie("access_token", token, 86400, "/", "localhost", true, true)

	userResponse := newUserResponse(user, token)

	ctx.JSON(http.StatusOK, userResponse)
}

type loginUserRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	user, err := server.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, responseError(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	token, err := server.token.CreateToken(user.ID, server.config.AccessTokenDuration)

	ctx.SetCookie("access_token", token, 86400, "/", "localhost", true, true)

	userResponse := newUserResponse(user, token)

	ctx.JSON(http.StatusOK, userResponse)
}
