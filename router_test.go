package webpack

import (
	"net/http"
	"testing"
)

const (
	entrypointfilepath = "../frontend/dist/index.html"
	resourcedirpath    = "../frontend/dist/"
)

func TestExampleUsage(t *testing.T) {
	router := NewRouter()
	entrypoints := NewEntrypointRouter()
	resources := NewResourceRouter()

	entrypoints.RegisterFile("/", entrypointfilepath)
	resources.RegisterDirectory(resourcedirpath)

	router.GET("/$", entrypoints.ServeHTTP)
	router.GET("(\\.js|\\.json|\\.css|\\.png|\\.gif|\\.jpe?g|\\.ico)$", resources.ServeHTTP)

	server := &http.Server{
		Addr:    "127.0.0.1:1339",
		Handler: router,
	}

	t.Log(server.ListenAndServe())
}
