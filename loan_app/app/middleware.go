package app

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/app/store"
	"github.com/goProjects/loan_app/lib"
	"github.com/goProjects/loan_app/lib/db"
	"github.com/goProjects/loan_app/lib/logger"
	"github.com/goProjects/loan_app/lib/utils"
	"github.com/goProjects/loan_app/lib/web"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

const (
	SERVICE_ID        = "Loan-Service-Id"
	SERVICE_NONCE     = "Loan-Service-Nonce"
	SERVICE_SIGNATURE = "Loan-Service-Signature"
	ACCESS_TOKEN      = "Access-Token"

	adminRole role = "ADMIN"
)

type (
	handlerType = func(r *web.Request) *web.APIResponse
	role        string
)

func ServeEndpoint(nextHandler handlerType) httprouter.Handle {
	return baseMiddleware(authenticateWithHmacDigestMiddleware(nextHandler))
}

func authenticateWithHmacDigestMiddleware(nextHandler handlerType) handlerType {
	return func(r *web.Request) *web.APIResponse {
		serviceID := r.GetReqHeader(SERVICE_ID)
		nonce := r.GetReqHeader(SERVICE_NONCE)
		serviceSignature := r.GetReqHeader(SERVICE_SIGNATURE)
		if !validateHmacDigest(serviceID, nonce, serviceSignature, getServiceAuthConfig()) {
			return web.ErrUnauthorizedRequest("unauthorized")
		}
		return nextHandler(r)
	}
}

func validateAccessTokenMiddleware(roles []role, nextHandler handlerType) handlerType {
	return func(r *web.Request) *web.APIResponse {
		accessToken := r.GetReqHeader(ACCESS_TOKEN)
		if accessToken == "" {
			return web.ErrUnauthorizedRequest("unauthorized")
		}

		userStore := store.NewUserStore()
		user, err := userStore.GetUserByAccessToken(r.Context(), db.Get(), accessToken)
		if err != nil {
			logger.E(r.Context(), "auth failed", zap.Error(err))
			return web.ErrUnauthorizedRequest("unauthorized")
		}
		r.Push(utils.CurrentUserIDStoreKey, user.ID)

		if ok := isRoleValid(user.UserRoles, roles); !ok {
			return web.ErrUnauthorizedRequest("unauthorized")
		}

		return nextHandler(r)
	}
}

func baseMiddleware(nextHandler handlerType) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		startTime := time.Now()

		webReq := web.NewRequest(req)
		for i := range ps {
			webReq.SetPathParam(ps[i].Key, ps[i].Value)
		}

		reqID := webReq.GetRequestIDFromRequestHeader()
		ctx := utils.SetRequestID(webReq.Context(), reqID)
		w.Header().Set(utils.RequestIDHeader, reqID)
		webReq = webReq.WithContext(ctx)

		defer lib.PanicHandler(req.Context(), nil, func(_ error) {
			apiResponse := web.ErrInternalServerError
			writeJsonAPIResponse(w, apiResponse)
			logRequest(webReq, apiResponse, startTime)
		})

		apiResponse := nextHandler(webReq)
		writeJsonAPIResponse(w, apiResponse)
		logRequest(webReq, apiResponse, startTime)
	}
}

func writeJsonAPIResponse(w http.ResponseWriter, resp *web.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.HTTPStatusCode())
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func logRequest(req *web.Request, resp *web.APIResponse, startTime time.Time) {
	logger.D(req.Context(), "response details", zap.Any("response", resp))

	logger.I(req.Context(), "Request processed",
		zap.Int("status", resp.HTTPStatusCode()),
		zap.String("route", req.GetRoute()),
		zap.Any("path_params", req.GetPathParams()),
		zap.Any("query_params", req.QueryParams()),
		zap.Float64("duration_ms", float64(time.Since(startTime).Nanoseconds())/1e6),
		zap.String("method", req.Method),
	)
}

func validateHmacDigest(serviceID, nonce, serviceSignature string, serviceConfig map[string]string) bool {
	if serviceID == "" || nonce == "" || serviceSignature == "" {
		return false
	}
	message := strings.Join([]string{nonce, serviceID}, "-")
	digest := calculateHmacDigest(message, serviceConfig[serviceID])
	return digest == serviceSignature
}

func calculateHmacDigest(data, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func getServiceAuthConfig() map[string]string {
	config := os.Getenv("SERVICE_AUTH_CONFIG")
	res := make(map[string]string)

	for _, confString := range strings.Split(config, "|") {
		confPair := strings.Split(confString, ":")
		if len(confPair) < 2 {
			continue
		}
		res[confPair[0]] = confPair[1]
	}
	return res
}

func isRoleValid(userRoles []*model.UserRole, requiredRoles []role) bool {
	for _, rolea := range requiredRoles {
		for _, roleb := range userRoles {
			if !strings.EqualFold(string(rolea), roleb.Role.Name) {
				return false
			}
		}
	}

	return true
}
