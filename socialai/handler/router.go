package handler

import (
	"net/http"
	"os"
	"strings"

	"socialai/util"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var signingKey []byte

func InitRouter(tokenConfig *util.TokenInfo) http.Handler {
	signingKey = []byte(tokenConfig.Secret)

	router := mux.NewRouter()
	router.HandleFunc("/healthz", healthHandler).Methods(http.MethodGet)

	protected := router.NewRoute().Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return withJWTAuth(next, signingKey)
	})
	protected.HandleFunc("/upload", uploadHandler).Methods(http.MethodPost)
	protected.HandleFunc("/search", searchHandler).Methods(http.MethodGet)
	protected.HandleFunc("/post/{id}", deleteHandler).Methods(http.MethodDelete)
	protected.HandleFunc("/generate-image", generateImageHandler).Methods(http.MethodPost)

	router.HandleFunc("/signup", signupHandler).Methods(http.MethodPost)
	router.HandleFunc("/signin", signinHandler).Methods(http.MethodPost)
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	return handlers.CORS(
		handlers.AllowedOrigins(allowedOrigins()),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions}),
	)(router)
}

func allowedOrigins() []string {
	configuredOrigins := strings.TrimSpace(os.Getenv("SOCIALAI_ALLOWED_ORIGINS"))
	if configuredOrigins == "" {
		return []string{"http://localhost:3000"}
	}

	origins := strings.Split(configuredOrigins, ",")
	for index := range origins {
		origins[index] = strings.TrimSpace(origins[index])
	}
	return origins
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
