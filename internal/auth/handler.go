package auth

import (
	"errors"
	"net/http"
	"time"

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

	// Read request body
	var req RegisterRequest
	_, err := httpData.ReadRequestBody(r, &req)
	if err != nil {
		h.appLogger.Errorf("REGISTER ERROR: read_request_body error: %v", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request!"))
		return
	}

	// Validate request body params
	if req.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}
	if req.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	// Register with errors
	if err := h.authService.RegisterUser(req.Email, req.Password); err != nil {
		h.appLogger.Errorf("REGISTER ERROR: %v", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		switch {
		case errors.Is(err, ErrUserExist):
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Email already registered!"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error!"))
			return
		}
	}

	// Response OK
	w.WriteHeader(http.StatusOK)
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Read request body
	var req LoginRequest
	_, err := httpData.ReadRequestBody(r, &req)
	if err != nil {
		h.appLogger.Errorf("LOGIN ERROR: read_request_body error: %v", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request!"))
		return
	}

	// Validate request body params
	if req.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}
	if req.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	// login with errors
	accessToken, refreshToken, err := h.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			resp := ErrorResponse{
				Error: "Email or password incorrect",
			}
			httpData.ResponseJSON(w, resp, http.StatusUnauthorized)
			return
		case errors.Is(err, ErrWrongPassword):
			resp := ErrorResponse{
				Error: "Email or password incorrect",
			}
			httpData.ResponseJSON(w, resp, http.StatusUnauthorized)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error!"))
			return
		}
	}

	// set cookie
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(360 * time.Hour), // срок действия
		HttpOnly: true,                            // запрет доступа из JS
		Secure:   false,                           // передавать только по HTTPS (рекомендуется)
		Path:     "/",
		SameSite: http.SameSiteStrictMode, // дополнительная защита
	}

	http.SetCookie(w, cookie)

	resp := LoginResponse{
		AccessToken: accessToken,
	}
	httpData.ResponseJSON(w, resp, http.StatusOK)
}

func (h *authHandler) ReadRefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			http.Error(w, "Cookie не найден", http.StatusBadRequest)
		default:
			http.Error(w, "Ошибка чтения cookie", http.StatusInternalServerError)
		}
		return
	}

	// Используем значение cookie
	value := cookie.Value
	w.Write([]byte("Значение cookie: " + value))
}
