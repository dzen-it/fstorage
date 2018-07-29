package api

import (
	"github.com/dzen-it/fstorage/storage"
)

type API struct {
	filestorage   storage.Storage
	filesizeLimit int64
}

func NewAPI(filestorage storage.Storage, filesizeLimit int64) *API {
	return &API{
		filestorage:   filestorage,
		filesizeLimit: filesizeLimit,
	}
}
