package routes

import (
	"net/http"

	"github.com/Laelapa/PlateOps/auth/tokenauthority"
	"github.com/Laelapa/PlateOps/internal/repository"
	"github.com/Laelapa/PlateOps/internal/routes/handlers"

	"github.com/Laelapa/GoHome/logging"
)

// Setup initializes and returns a configured router with all application routes
// as well as static file serving.
//
// Parameters:
//   - staticDir: The directory containing static files to serve
func Setup(staticDir string, logger *logging.Logger, queries *repository.Queries, tokenAuthority *tokenauthority.TokenAuthority) *http.ServeMux {

	mux := http.NewServeMux()
	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))

	h := handlers.New(logger, queries, tokenAuthority)

	// -- Commented out routes are not implemented yet.
	// They will be serving HTML pages in the future.
	mux.Handle("GET /static/", fileServer)
	mux.HandleFunc("GET /health", h.HandleGetHealth)
	// -- mux.HandleFunc("GET /signup", h.HandleGetSignup)
	mux.HandleFunc("POST /signup", h.HandlePostSignup)
	// -- mux.HandleFunc("GET /login", h.HandleGetLogin)
	// mux.HandleFunc("POST /login", h.HandlePostLogin)
	// mux.HandleFunc("POST /logout", h.HandleGetLogout)
	// -- mux.HandleFunc("GET /reset-password", h.HandleGetResetPassword)
	// mux.HandleFunc("POST /reset-password", h.HandlePostResetPassword)

	return mux
}
