package handlers

import "net/http"

func (h *Handler)HandleSwaggerUI(w http.ResponseWriter, r *http.Request) {
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
	if _,err := w.Write([]byte(html)); err != nil {
		h.logger.LogAppError("Couldn't write response", err)
	}
}
