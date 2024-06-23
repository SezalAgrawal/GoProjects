package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/goProjects/loan_app/lib/logger"
	"github.com/goProjects/loan_app/lib/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"errors"
)

var (
	headersToBeMasked = []string{"Authorization"}
)

type Request struct {
	*http.Request

	route      string
	pathParams map[string]string
	params     map[string]string
	store      map[string]interface{}
}

func NewRequest(r *http.Request) *Request {
	webRequest := &Request{
		Request: r,
		route:   getPath(r),
	}
	return webRequest
}

func (r *Request) WithContext(ctx context.Context) *Request {
	if ctx == nil {
		panic("nil context")
	}
	httpRequest := r.Request.WithContext(ctx)
	r2 := new(Request)
	*r2 = *r
	r2.Request = httpRequest

	return r2
}

func (r *Request) GetRoute() string {
	return r.route
}

func (r *Request) GetPathParams() map[string]string {
	return r.pathParams
}

func (r *Request) GetPathParam(key string) string {
	if value, ok := r.pathParams[key]; ok {
		return value
	}
	return ""
}

func (r *Request) SetPathParam(key, value string) {
	if r.pathParams == nil {
		r.pathParams = make(map[string]string)
	}
	r.route = strings.Replace(r.route, value, fmt.Sprintf(":%s", key), 1)
	r.pathParams[key] = value
}

func (r *Request) QueryParam(key string) string {
	return r.URL.Query().Get(key)
}

func (r *Request) QueryParams() map[string]string {
	if r.params != nil {
		return r.params
	}
	r.params = map[string]string{}
	for key, val := range r.URL.Query() {
		r.params[key] = strings.Join(val, " | ")
	}
	return r.params
}

func (r *Request) QueryParamExists(keys ...string) bool {
	for _, key := range keys {
		if r.URL.Query().Get(key) == "" {
			return false
		}
	}
	return true
}

// Headers shouldn't be used for logging
func (r *Request) Headers() map[string]interface{} {
	headers := map[string]interface{}{}
	for key, value := range r.Header {
		if strings.ToLower(key) == "content-type" {
			continue
		}
		if strings.ToLower(key) == "accept" {
			continue
		}
		headers[key] = value
	}
	return headers
}

// MaskedHeaders returns a request headers with masked values read from an array
func (r *Request) MaskedHeaders() http.Header {
	headers := r.Header.Clone()
	for _, key := range headersToBeMasked {
		k := textproto.CanonicalMIMEHeaderKey(key)
		_, ok := headers[k]
		if ok {
			headers.Set(key, "*******")
		}
	}
	return headers
}

func (r *Request) ReadBody() (map[string]interface{}, error) {
	bodyMap := make(map[string]interface{})

	if r.ContentLength == 0 {
		err := errors.New("empty body")
		return bodyMap, err
	}

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		logger.W(r.Context(), "Error reading request body", zap.String("error", err.Error()))
		return bodyMap, err
	}

	bodyMap, err = unmarshalRequestBody(bodyByte)
	if err != nil {
		logger.W(r.Context(), "Error decoding request json",
			zap.String("error", err.Error()), zap.Any("headers", r.MaskedHeaders()),
			zap.String("url", r.URL.String()), zap.String("body", string(bodyByte)))
	}
	return bodyMap, err
}

func (r *Request) Bind(v interface{}) error {
	body := r.Request.Body
	defer body.Close()

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	return d.Decode(v)
}

func (r *Request) GetRequestIP() string {
	fIps := r.Header["X-Forwarded-For"]
	if len(fIps) < 1 {
		if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			return ip
		}

		return net.ParseIP(r.RemoteAddr).String()
	}
	return strings.TrimSpace(strings.Split(fIps[0], ",")[0])
}

func (r *Request) GetReqHeader(key string) string {
	if len(r.Header[key]) < 1 {
		return ""
	}
	return r.Header[key][0]
}

func (r *Request) GetRequestIDFromRequestHeader() string {
	requestID := r.GetReqHeader(string(utils.RequestIDHeader))
	if requestID == "" {
		requestID = uuid.NewString()
	}
	return requestID
}

func (r *Request) ValidateBodyToStruct(s interface{}, structValidations ...validator.StructLevelFunc) error {
	err := r.Bind(s)
	if err == io.EOF {
		reqBody, e := r.ReadBody()
		if e == nil {
			jsonString, _ := json.Marshal(reqBody)
			err = json.Unmarshal(jsonString, &s)
		}
	}

	if err != nil {
		return utils.HandleValidationErrors(err)
	}

	return utils.ValidateStruct(s, structValidations...)
}

func (r *Request) ValidateParamsToStruct(s interface{}, structValidations ...validator.StructLevelFunc) error {
	jsonString, _ := json.Marshal(r.QueryParams())
	err := json.Unmarshal(jsonString, &s)
	if err != nil {
		return utils.HandleValidationErrors(err)
	}
	return utils.ValidateStruct(s, structValidations...)
}

func unmarshalRequestBody(body []byte) (map[string]interface{}, error) {
	bodyMap := make(map[string]interface{})
	b := bytes.NewBuffer(body)
	decoder := json.NewDecoder(b)
	decoder.UseNumber()
	err := decoder.Decode(&bodyMap)
	return bodyMap, err
}

func getPath(r *http.Request) string {
	path := "undefined"
	if r != nil && r.URL != nil && r.URL.Path != "" {
		path = r.URL.Path
	}
	return path
}

func (r *Request) Push(key string, value interface{}) {
	if r.store == nil {
		r.store = map[string]interface{}{}
	}
	r.store[key] = value
}

func (r *Request) Value(key string) interface{} {
	return r.store[key]
}
