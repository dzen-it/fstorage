package api

import (
	"net/http"

	"github.com/dzen-it/fstorage/storage"
	"github.com/dzen-it/fstorage/utils/log"

	"github.com/go-chi/chi"
)

func (api *API) DeleteFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")

		err := api.filestorage.DeleteFile(hash)
		if err == storage.ErrFileNotFound {
			writeErrorHTTP(w, ErrFileNotFound)
			return
		}

		if err != nil {
			log.Errorw("error delete file from storage.", "hash", hash, "error", err)
			writeErrorHTTP(w, ErrInternalError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
