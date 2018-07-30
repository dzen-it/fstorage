package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	maxLengthFileName = 255
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func MiddlewareValidateFilename(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := chi.URLParam(r, "filename")
		l := len(filename)

		if l == 0 || l > maxLengthFileName {
			writeErrorHTTP(w, ErrInvalidFilename)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MiddlewareValidateHash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")

		if !regexpMatchKeccak256HexString.MatchString(hash) {
			writeErrorHTTP(w, ErrFileNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MiddlewareValidateHeadersControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers, err := getHeadersHashControl(r)
		if err != nil {
			writeErrorHTTP(w, err)
			return
		}

		ctx := r.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		ctx = context.WithValue(ctx, "headerscontrol", headers)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
