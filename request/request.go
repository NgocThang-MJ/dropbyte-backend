package request

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type authResponse struct {
	AccountId          string `json:"accountId"`
	AuthorizationToken string `json:"authorizationToken"`
	StatusCode         int    `json:"statusCode"`
}

type urlResponse struct {
	UploadUrl          string `json:"uploadUrl"`
	AuthorizationToken string `json:"authorizationToken"`
	BucketId           string `json:"bucketId"`
	StatusCode         int    `json:"statusCode"`
}

type deleteFileResponse struct {
	FileId     string `json:"fileId"`
	FileName   string `json:"fileName"`
	StatusCode int    `json:"statusCode"`
}

type responseBodyOnError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func AuthorizeAccount(accountId string, applicationKey string) (response authResponse, err error) {
	request, err := http.NewRequest(
		http.MethodGet,
		"https://api.backblazeb2.com/b2api/v3/b2_authorize_account",
		nil,
	)
	if err != nil {
		return
	}

	basicAuth := base64.StdEncoding.EncodeToString(
		[]byte(accountId + ":" + applicationKey),
	)
	request.Header.Set("Authorization", "Basic "+basicAuth)

	res, err := http.DefaultClient.Do(request)
	response.StatusCode = res.StatusCode
	if res.StatusCode != 200 {
		return response, errors.New("Failed to authorize in Backblaze")
	}
	if err != nil {
		return
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(resBody), &response)

	return response, nil
}

func GetUploadUrl(bucketId string, authToken string) (response urlResponse, err error) {
	request, err := http.NewRequest(
		http.MethodGet,
		"https://api005.backblazeb2.com/b2api/v3/b2_get_upload_url?bucketId="+bucketId,
		nil,
	)
	if err != nil {
		return
	}

	request.Header.Set("Authorization", authToken)

	res, err := http.DefaultClient.Do(request)
	response.StatusCode = res.StatusCode
	if res.StatusCode != 200 {
		return response, errors.New("Failed to get upload url")
	}
	if err != nil {
		return
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(resBody), &response)

	return response, nil
}

func DeleteFileById(
	fileId string,
	fileName string,
	authToken string,
) (response deleteFileResponse, err error) {
	jsonBody := []byte(fmt.Sprintf(`{"fileName" : "%s", "fileId": "%s"}`, fileName, fileId))
	bodyReader := bytes.NewReader(jsonBody)
	request, err := http.NewRequest(
		http.MethodPost,
		"https://api005.backblazeb2.com/b2api/v2/b2_delete_file_version",
		bodyReader,
	)
	if err != nil {
		return
	}

	request.Header.Set("Authorization", authToken)

	res, err := http.DefaultClient.Do(request)
	response.StatusCode = res.StatusCode

	if res.StatusCode != 200 {
		var errorResponse responseBodyOnError
		err = json.NewDecoder(res.Body).Decode(&errorResponse)
		if err != nil {
			fmt.Println("Error reading response body:", err)
		}
		return response, errors.New(errorResponse.Message)
	}
	if err != nil {
		return
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(resBody), &response)
	if err != nil {
		return
	}

	return response, nil
}

func DownloadFileById(fileId string, authToken string) (resBody []byte, err error) {
	request, err := http.NewRequest(
		http.MethodGet,
		"https://api005.backblazeb2.com/b2api/v2/b2_download_file_by_id?fileId="+fileId,
		nil,
	)
	if err != nil {
		return
	}

	request.Header.Set("Authorization", authToken)

	res, err := http.DefaultClient.Do(request)
	if res.StatusCode != 200 {
		var errorResponse responseBodyOnError
		err = json.NewDecoder(res.Body).Decode(&errorResponse)
		if err != nil {
			fmt.Println("Error reading response body:", err)
		}
		return resBody, errors.New(errorResponse.Message)
	}
	if err != nil {
		return
	}
	resBody, err = ioutil.ReadAll(res.Body)

	fmt.Println("Response Body:", string(resBody))

	if err != nil {
		return
	}
	return resBody, nil
}
