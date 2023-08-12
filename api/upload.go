package api

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	db "github.com/liquiddev99/dropbyte-backend/db/sqlc"
	"github.com/liquiddev99/dropbyte-backend/request"
	"github.com/liquiddev99/dropbyte-backend/token"
)

type responseFile struct {
	FileID   string `json:"fileId"`
	BucketID string `json:"bucketId"`
	FileName string `json:"fileName"`
	Size     uint   `json:"contentLength"`
	FileType string `json:"contentType"`
}

func (server *Server) guestUploadFile(ctx *gin.Context) {
	// Get file information
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}
	defer openedFile.Close()

	fileContent := &bytes.Buffer{}
	io.Copy(fileContent, openedFile)

	_, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	openedHeader, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}
	defer openedHeader.Close()

	// Get Authorization Token
	authResponse, err := request.AuthorizeAccount(
		server.config.AccountId,
		server.config.ApplicationKey,
	)
	if err != nil {
		ctx.JSON(authResponse.StatusCode, responseError(err))
		return
	}

	// Get Upload Url
	urlResponse, err := request.GetUploadUrl(
		server.config.BucketId,
		authResponse.AuthorizationToken,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	// Upload file
	uploadRequest, err := http.NewRequest(
		"POST",
		urlResponse.UploadUrl,
		fileContent,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	sha1Hash := sha1.New()
	if _, err := io.Copy(sha1Hash, openedHeader); err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	contentSHA1 := hex.EncodeToString(sha1Hash.Sum(nil))

	uploadRequest.Header.Set(
		"Authorization",
		urlResponse.AuthorizationToken,
	)
	uploadRequest.Header.Set("X-Bz-File-Name", url.QueryEscape(file.Filename))
	uploadRequest.Header.Set("Content-Length", fmt.Sprintf("%d", fileContent.Len()))
	uploadRequest.Header.Set("Content-Type", "b2/x-auto")
	uploadRequest.Header.Set("X-Bz-Content-Sha1", contentSHA1)

	// Upload the file to specific dst.
	client := &http.Client{}

	uploadResp, err := client.Do(uploadRequest)
	if uploadResp.StatusCode != 200 {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}
	defer uploadResp.Body.Close()

	var responseData responseFile
	err = json.NewDecoder(uploadResp.Body).Decode(&responseData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	fileArg := db.CreateFileParams{
		FileID:   responseData.FileID,
		BucketID: responseData.BucketID,
		Size:     fmt.Sprintf("%d", responseData.Size),
		Name:     responseData.FileName,
		FileType: responseData.FileType,
	}

	_, err = server.db.CreateFile(ctx, fileArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	ctx.JSON(http.StatusOK, responseData)
}

func (server *Server) userUploadFile(ctx *gin.Context) {
	authPayload := ctx.MustGet("payload").(*token.Payload)

	// Get file information
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}
	defer openedFile.Close()

	fileContent := &bytes.Buffer{}
	io.Copy(fileContent, openedFile)

	_, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	openedHeader, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}
	defer openedHeader.Close()

	// Get Authorization Token
	authResponse, err := request.AuthorizeAccount(
		server.config.AccountId,
		server.config.ApplicationKey,
	)
	if err != nil {
		ctx.JSON(authResponse.StatusCode, responseError(err))
		return
	}

	// Get Upload Url
	urlResponse, err := request.GetUploadUrl(
		server.config.BucketId,
		authResponse.AuthorizationToken,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	// Upload file
	uploadRequest, err := http.NewRequest(
		"POST",
		urlResponse.UploadUrl,
		fileContent,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	sha1Hash := sha1.New()
	if _, err := io.Copy(sha1Hash, openedHeader); err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	contentSHA1 := hex.EncodeToString(sha1Hash.Sum(nil))

	uploadRequest.Header.Set(
		"Authorization",
		urlResponse.AuthorizationToken,
	)
	uploadRequest.Header.Set("X-Bz-File-Name", url.QueryEscape(file.Filename))
	uploadRequest.Header.Set("Content-Length", fmt.Sprintf("%d", fileContent.Len()))
	uploadRequest.Header.Set("Content-Type", "b2/x-auto")
	uploadRequest.Header.Set("X-Bz-Content-Sha1", contentSHA1)

	// Upload the file to specific dst.
	client := &http.Client{}

	uploadResp, err := client.Do(uploadRequest)
	if uploadResp.StatusCode != 200 {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}
	defer uploadResp.Body.Close()

	var responseData responseFile
	err = json.NewDecoder(uploadResp.Body).Decode(&responseData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	fileArg := db.CreateFileParams{
		FileID:   responseData.FileID,
		BucketID: responseData.BucketID,
		Owner:    authPayload.UserId,
		Size:     fmt.Sprintf("%d", responseData.Size),
		Name:     responseData.FileName,
		FileType: responseData.FileType,
	}

	_, err = server.db.CreateFile(ctx, fileArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	ctx.JSON(http.StatusOK, responseData)
}
