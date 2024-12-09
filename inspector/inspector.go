package inspector

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/llmcontext/gomcp/config"
)

//go:embed html
var html embed.FS

var t = template.Must(template.ParseFS(html, "html/*"))

func StartInspector(config *config.InspectorInfo) {
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	})

	server := &http.Server{
		Addr:    config.ListenAddress,
		Handler: router,
	}

	server.ListenAndServe()
}
