package static_store

import "github.com/AvenciaLab/avencia-backend/lib/core"

type (
	StaticFileCreator = func(file core.File, dir, filename string) (string, error)
	StaticDirDeleter  = func(dir string) error
)
