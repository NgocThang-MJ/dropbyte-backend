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
	"github.com/kurin/blazer/b2"

	db "github.com/liquiddev99/dropbyte-backend/db/sqlc"
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

	// Upload file
	uploadRequest, err := http.NewRequest(
		"POST",
		server.b2UploadUrl,
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
		server.b2UrlAuthToken,
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

	// Upload file
	uploadRequest, err := http.NewRequest(
		"POST",
		server.b2UploadUrl,
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
		server.b2UrlAuthToken,
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

func (server *Server) guestUploadFileB2(ctx *gin.Context) {
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

	//	asdf, err := b2.NewClient(ctx, server.config.B2ApplicationKeyId, server.config.B2ApplicationKey)
	b2, err := b2.NewClient(ctx, server.config.B2ApplicationKeyId, server.config.B2ApplicationKey)
	bucket, err := b2.Bucket(ctx, "liquiddev99")

	obj := bucket.Object(file.Filename)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, openedFile); err != nil {
		w.Close()
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}
	w.Close()

	ctx.String(http.StatusOK, "Ok")
}
