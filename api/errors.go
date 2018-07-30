package api

import (
	"net/http"

	"github.com/dzen-it/fstorage/utils/log"
)

type Error uint16

func (e Error) Error() string {
	return ""
}

const (
	ErrInternalError Error = iota
	ErrEmptyFile
	ErrFileTransferFailure
	ErrInvalidFilename
	ErrFileNotFound
	ErrInvalidHeader
	ErrTheseHashesDoNotMatch
	ErrMaxSizeFileLimit
)

func writeErrorHTTP(w http.ResponseWriter, err error) {
	switch err {
	case ErrInternalError:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	case ErrEmptyFile:
		http.Error(w, "Could not uplpoad empty file", http.StatusBadRequest)
	case ErrFileTransferFailure:
		return
	case ErrInvalidFilename:
		http.Error(w, "Invalid filename", http.StatusBadRequest)
	case ErrFileNotFound:
		http.Error(w, "File not found", http.StatusNotFound)
	case ErrInvalidHeader:
		http.Error(w, "Invalid header value format ", http.StatusBadRequest)
	case ErrTheseHashesDoNotMatch:
		http.Error(w, "These hashes do not match", http.StatusBadRequest)
	case ErrMaxSizeFileLimit:
		http.Error(w, "File size too large", http.StatusRequestEntityTooLarge)
	default:
		log.Errorw("Undefined error", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
