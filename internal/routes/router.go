package routes

import (
	"net/http"

	"github.com/Laelapa/PlateOps/internal/routes/handlers"

	"github.com/Laelapa/GoHome/logging"
)

// Setup initializes and returns a configured router with all application routes
// as well as static file serving.
//
// Parameters:
//   - staticDir: The directory containing static files to serve
func Setup(staticDir string, logger *logging.Logger) *http.ServeMux {

	mux := http.NewServeMux()
	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))

	h := &handlers.Handler{
		Logger: logger,
	}

	mux.Handle("GET /static/", fileServer)
	mux.HandleFunc("GET /health", h.HandleGetHealth)

	return mux
}
