package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	db "github.com/liquiddev99/dropbyte-backend/db/sqlc"
	"github.com/liquiddev99/dropbyte-backend/request"
	"github.com/liquiddev99/dropbyte-backend/token"
)

type deteleFileRequest struct {
	FileId   string `json:"file_id"   binding:"required"`
	FileName string `json:"file_name" binding:"required"`
}

type downloadFileRequest struct {
	FileId string `form:"file_id" binding:"required"`
}

func (server *Server) getFiles(ctx *gin.Context) {
	authPayload := ctx.MustGet("payload").(*token.Payload)

	arg := db.ListFilesParams{
		Owner:  authPayload.UserId,
		Limit:  50,
		Offset: 0,
	}

	files, err := server.db.ListFiles(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	ctx.JSON(http.StatusOK, files)
}

func (server *Server) downloadFileById(ctx *gin.Context) {
	var req downloadFileRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	authResponse, err := request.AuthorizeAccount(
		server.config.B2ApplicationKeyId,
		server.config.B2ApplicationKey,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	res, err := request.DownloadFileById(req.FileId, authResponse.AuthorizationToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	ctx.Data(http.StatusOK, "application/octet-stream", res)
}

func (server *Server) deleteFileById(ctx *gin.Context) {
	var req deteleFileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	authResponse, err := request.AuthorizeAccount(
		server.config.B2ApplicationKeyId,
		server.config.B2ApplicationKey,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	res, err := request.DeleteFileById(req.FileId, req.FileName, authResponse.AuthorizationToken)
	if err != nil {
		ctx.JSON(res.StatusCode, responseError(err))
		return
	}

	err = server.db.DeleteFile(ctx, req.FileId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}
