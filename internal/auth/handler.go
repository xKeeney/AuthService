package auth

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/xKeeney/httpForge/httpData"
	"github.com/xKeeney/httpForge/httpLogger"
)

type authHandler struct {
	authService *authService
	appLogger   *httpLogger.HttpLogger
}

func InitAuthHandler(authService *authService, appLogger *httpLogger.HttpLogger) *authHandler {
	return &authHandler{
		authService: authService,
		appLogger:   appLogger,
	}
}

func (h *authHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email        string `json:"email"`
		PasswordHash string `json:"password_hash"`
		Status       string `json:"status"`
	}

	type ResponseBody struct {
		Status string `json:"status"`
	}

	var requestBody RequestBody

	_, err := httpData.ReadRequestBody(r, &requestBody)
	if err != nil {
		h.appLogger.Errorf("CreateUser read request body error: %v", err)
		return
	}

	uuidstr := uuid.NewString()
	if err := h.authService.CreateUser(
		uuidstr,
		requestBody.Email,
		requestBody.PasswordHash,
		requestBody.Status,
	); err != nil {
		h.appLogger.Errorf("CreateUser error: %v", err)
		return
	}

	responseBody := ResponseBody{
		Status: "success",
	}

	httpData.ResponseJSON(w, responseBody, 200)
}
