package swagger

import (
	"net/http"
	"path"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
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

func ServeSwaggerUI(mux *http.ServeMux) {
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}
