package controller_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	"github.com/goProjects/loan_app/app"
	"github.com/goProjects/loan_app/lib/utils"
)

func sendRequest(t *testing.T, httpMethod, path string, body map[string]interface{}, queryParam, headers map[string]string) (responseCode int, responseBody map[string]interface{}) {
	data, err := json.Marshal(body)
	assert.Nil(t, err)

	serviceId, serviceKey := uuid.New().String(), uuid.New().String()
	t.Setenv("SERVICE_AUTH_CONFIG", serviceId+":"+serviceKey)

	var reqBody io.Reader = http.NoBody
	if data != nil {
		reqBody = bytes.NewBuffer(data)
	}
	req := httptest.NewRequest(httpMethod, path, reqBody)
	setAuthHeaders(req, serviceKey, serviceId)
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	q := make(url.Values)

	for key, value := range queryParam {
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()

	router := httprouter.New()
	app.InitRoutes(router)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	responseBody = make(map[string]interface{})

	err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	return recorder.Code, responseBody
}

func setAuthHeaders(req *http.Request, serviceKey, serviceId string) {
	nonce := utils.NewUUID()

	key := fmt.Sprintf("%s-%s", nonce, serviceId)
	signature := hmac.New(sha1.New, []byte(serviceKey))
	signature.Write([]byte(key))
	serviceSignature := hex.EncodeToString(signature.Sum(nil))

	req.Header.Set("LOAN-SERVICE-ID", serviceId)
	req.Header.Set("LOAN-SERVICE-NONCE", nonce)
	req.Header.Set("LOAN-SERVICE-SIGNATURE", serviceSignature)
}
