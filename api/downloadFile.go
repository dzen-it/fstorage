package api

import (
	"io"
	"net/http"

	"github.com/dzen-it/fstorage/storage"
	"github.com/dzen-it/fstorage/utils/log"

	"github.com/go-chi/chi"
)

func (api *API) DownloadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")

		file, err := api.filestorage.GetFile(hash)
		if err == storage.ErrFileNotFound {
			writeErrorHTTP(w, ErrFileNotFound)
			return
		}

		if err != nil {
			log.Errorw("error get file from storage.", "hash", hash, "error", err)
			writeErrorHTTP(w, ErrInternalError)
			return
		}

		w.WriteHeader(http.StatusOK)

		if _, err = io.Copy(w, file); err != nil {
			log.Errorw("error write file", "hash", hash, "error", err)
			writeErrorHTTP(w, ErrInternalError)
			return
		}
	}
}
