package static_store

type (
	StaticFileCreator = func(data *[]byte, dir, filename string) (string, error)
	StaticDirDeleter  = func(dir string) error
)
