package swagger

import (
	"net/http"
	"path"
	"strings"

	log "github.com/samwang0723/jarvis/internal/logger"
)

func ServeSwaggerFile(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "swagger.json") {
		log.Errorf("Not Found: %s", r.URL.Path)
		http.NotFound(w, r)

		return
	}

	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join("api/", p)

	log.Infof("Serving swagger-file: %s", p)

	http.ServeFile(w, r, p)
}

func ServeAnalysisFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path)
}
