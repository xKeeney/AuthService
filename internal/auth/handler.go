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

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	_, err := httpData.ReadRequestBody(r, &req)
	if err != nil {
		h.appLogger.Errorf("REGISTER ERROR: read_request_body error: %v", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request!"))
		return
	}

	if req.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}
	if req.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	result, err := h.authService.RegisterUser(req.Email, req.Password)
	if err != nil {
		h.appLogger.Errorf("REGISTER ERROR: %v", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error!"))
		return
	}

	if !result {
		h.appLogger.Errorf("REGISTER ERROR: %v", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("User with email already exist!"))
		return
	}

	w.WriteHeader(http.StatusOK)
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
