package request

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/liquiddev99/dropbyte-backend/util"
)

func TestAuthorizeAccount(t *testing.T) {
	config, err := util.LoadConfig("../")
	require.NoError(t, err)

	testCases := []struct {
		name          string
		accountId     string
		appKey        string
		checkResponse func(t *testing.T, response authResponse, err error)
	}{
		{
			name:      "OK",
			accountId: config.AccountId,
			appKey:    config.ApplicationKey,
			checkResponse: func(t *testing.T, response authResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, response.AccountId, config.AccountId)
			},
		},
		{
			name:      "Unauthorized",
			accountId: config.AccountId,
			appKey:    "Invalid",
			checkResponse: func(t *testing.T, response authResponse, err error) {
				require.Error(t, err)
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			response, err := AuthorizeAccount(testCase.accountId, testCase.appKey)

			testCase.checkResponse(t, response, err)
		})

	}
}

func TestGetUploadUrl(t *testing.T) {
	config, err := util.LoadConfig("../")
	require.NoError(t, err)

	testCases := []struct {
		name          string
		bucketId      string
		checkResponse func(t *testing.T, response urlResponse, err error)
	}{
		{
			name:     "OK",
			bucketId: config.BucketId,
			checkResponse: func(t *testing.T, response urlResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, response.BucketId, config.BucketId)
			},
		},
		{
			name:     "Unauthorized",
			bucketId: "Invalid",
			checkResponse: func(t *testing.T, response urlResponse, err error) {
				require.Error(t, err)
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			authResponse, err := AuthorizeAccount(config.AccountId, config.ApplicationKey)
			require.NoError(t, err)

			urlResponse, err := GetUploadUrl(testCase.bucketId, authResponse.AuthorizationToken)

			testCase.checkResponse(t, urlResponse, err)
		})

	}
}
