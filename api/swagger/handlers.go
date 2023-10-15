package swagger

import (
	"net/http"
	"path"
	"strings"

	"github.com/rs/zerolog"
)

func ServeSwaggerFile(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "swagger.json") {
		zerolog.Ctx(r.Context()).Error().Msgf("Not Found: %s", r.URL.Path)
		http.NotFound(w, r)

		return
	}

	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join("api/", p)

	zerolog.Ctx(r.Context()).Info().Msgf("Serving swagger-file: %s", p)

	http.ServeFile(w, r, p)
}

func ServeAnalysisFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path)
}
