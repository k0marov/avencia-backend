package static_store

import (
	"log"
	"os"
)

type (
	StaticFileCreator = func(data *[]byte, dir, filename string) (string, error)
	StaticDirDeleter  = func(dir string) error
)

var StaticDir = getStaticDir()
var StaticHost = getStaticHostStr()

func PathToURL(path string) string {
	if path == "" {
		return ""
	}
	return StaticHost + "/" + path
}

func getStaticHostStr() string {
	const staticHostEnv = "AVENCIA_STATIC_HOST"
	host, exists := os.LookupEnv(staticHostEnv)
	if !exists {
		log.Fatalf(`Environment variable %s is not set.
If this is a test, just set the environment variable to a dummy string.
If this is in production, set this environment variable to point to the URL from which the static directory can be accessed.`, staticHostEnv)
	}
	return host
}

func getStaticDir() string {
	const staticDirEnv = "AVENCIA_STATIC_DIR"
	dir, exists := os.LookupEnv(staticDirEnv)
	if !exists {
		log.Fatalf(`Environment variable %s is not set.
If this is a test, just set the environment variable to some path like ./static/.
If this is in production, set this environment variable to point to the file system path of a directory where static files will be stored`, staticDirEnv)
	}
	return dir
}
