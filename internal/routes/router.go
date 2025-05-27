package routes

import (
	"net/http"

	"github.com/Laelapa/PlateOps/auth/tokenauthority"
	"github.com/Laelapa/PlateOps/internal/middleware"
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

	// to simplify the mux.Handle parameters for authenticated routes
	withAuth := func (handler func(http.ResponseWriter, *http.Request)) http.Handler {
		return middleware.AuthenticateWithJWT(tokenAuthority, logger)(http.HandlerFunc(handler))
	}

	// -- Commented out routes are not implemented yet.
	// They will be serving HTML pages in the future.
	mux.Handle("GET /static/", fileServer)
	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
    	w.Header().Set("Content-Type", "application/json")
    	http.ServeFile(w, r, "docs/openapi.json")
	})

	// Add Swagger UI route
mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
    html := `<!DOCTYPE html>
	<html>
		<head>
    		<title>PlateOps API</title>
    		<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.22.0/swagger-ui.css" />
		</head>
		<body>
    		<div id="swagger-ui"></div>
    		<script src="https://unpkg.com/swagger-ui-dist@5.22.0/swagger-ui-bundle.js"></script>
    		<script>
        		SwaggerUIBundle({
            		url: '/openapi.json',
            		dom_id: '#swagger-ui',
            		presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.presets.standalone]
        		});
    		</script>
		</body>
	</html>`
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(html))
})


	mux.HandleFunc("GET /health", h.HandleGetHealth)
	// -- mux.HandleFunc("GET /signup", h.HandleGetSignup)
	mux.HandleFunc("POST /signup", h.HandlePostSignup)
	mux.HandleFunc("POST /refresh", h.HandlePostRefresh)
	// -- mux.HandleFunc("GET /login", h.HandleGetLogin)
	mux.HandleFunc("POST /login", h.HandlePostLogin)
	mux.HandleFunc("POST /logout", h.HandlePostLogout)
	// -- mux.HandleFunc("GET /reset-password", h.HandleGetResetPassword)
	// mux.HandleFunc("POST /reset-password", h.HandlePostResetPassword)

	// mux.HandleFunc("GET /food/id/{id}", h.HandleGetFoodById)
	// mux.HandleFunc("GET /food/gtin/{gtin}", h.HandleGetFoodByGtin)

	mux.Handle("POST /food", withAuth(h.HandlePostFood))
	// mux.HandleFunc("PUT /food/id/{id}", h.HandlePutFood)
	// mux.HandleFunc("DELETE /food/id/{id}", h.HandleDeleteFood)

	// mux.HandleFunc("GET /foods", h.HandleGetFoods)
	// mux.HandleFunc("GET /foods/category/{category}", h.HandleGetFoodsByCategory)
	// mux.HandleFunc("GET /foods/name/{name}", h.HandleGetFoodsByName)

	return mux
}
