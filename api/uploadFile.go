package api

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dzen-it/fstorage/storage"
	"github.com/dzen-it/fstorage/utils/log"

	"golang.org/x/crypto/sha3"
)

func (api *API) UploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		// filename := chi.URLParam(r, "filename")

		headersHashControl := r.Context().Value("headerscontrol").(map[string]string)

		reader := io.TeeReader(r.Body, &sizeLimiter{
			limit: api.filesizeLimit,
		})

		hash, err := processUploadFile(api.filestorage, reader, headersHashControl, sha3.New256(), api.filesizeLimit)
		if err != nil {
			writeErrorHTTP(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, hash)
	}
}

func processUploadFile(s storage.Storage, body io.Reader, expectedHashes map[string]string, basicHash hash.Hash, limitFileSize int64) (string, error) {
	tmpFile, err := ioutil.TempFile("", "fstorage-")
	if err != nil {
		log.Errorw("Error create temporary file", "err", err.Error())
		return "", ErrInternalError
	}

	defer func() {
		tmpFile.Close()
	}()

	var factedHashes map[string]hash.Hash

	// it's necessary for stream calculating of a hash
	// while writing a temporary file
	body, factedHashes = executeHashesForFileReader(expectedHashes, io.TeeReader(body, basicHash))

	if err = writeToTmpFile(tmpFile, body); err != nil {
		return "", err
	}

	if err = compareHashes(expectedHashes, factedHashes); err != nil {
		return "", err
	}

	tmp, err := os.Open(tmpFile.Name())
	if err != nil {
		log.Errorw("Open temperary file", "Error", err.Error())
		return "", ErrInternalError
	}
	defer tmp.Close()

	hashsum := hex.EncodeToString(basicHash.Sum(nil))
	if err = saveFile(s, hashsum, tmp); err != nil {
		return "", err
	}

	return hashsum, nil
}

func writeToTmpFile(tmpFile io.Writer, reader io.Reader) error {
	n, err := io.Copy(tmpFile, reader)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			return ErrFileTransferFailure
		}
		log.Errorw("Error uploading file", "Error", err.Error())
		return err
	}

	if n == 0 {
		return ErrEmptyFile
	}

	return nil
}

func saveFile(s storage.Storage, hash string, reader io.Reader) error {
	if err := s.SaveFile(hash, reader); err != nil {
		log.Errorw("error save file to storage", "error", err, "hash", hash)
		return ErrInternalError
	}

	return nil
}

// Calculate hashes from one thread of a Reader.
func executeHashesForFileReader(expectedHashes map[string]string, reader io.Reader) (io.Reader, map[string]hash.Hash) {
	var (
		h hash.Hash
	)

	factedHashes := make(map[string]hash.Hash)

	for header := range expectedHashes {
		switch header {
		case XFSHashControlKeccak256:
			h = sha3.New256()
		case XFSHashControlKeccak512:
			h = sha3.New512()
		case XFSHashControlMD5:
			h = md5.New()
		case XFSHashControlSHA1:
			h = sha1.New()
		case XFSHashControlSHA256:
			h = sha256.New()
		case XFSHashControlSHA512:
			h = sha512.New()
		}

		factedHashes[header] = h
		reader = io.TeeReader(reader, h)
	}

	return reader, factedHashes
}

// returns an error if one of the hashes does not match
func compareHashes(expectedHashes map[string]string, factedHashes map[string]hash.Hash) error {
	var (
		header  string
		hash    hash.Hash
		hashsum string
	)

	for header, hash = range factedHashes {
		hashsum = hex.EncodeToString(hash.Sum(nil))
		if hashsum != expectedHashes[header] {
			return ErrTheseHashesDoNotMatch
		}
	}

	return nil
}

type sizeLimiter struct {
	n     int64
	limit int64
}

func (s *sizeLimiter) Write(p []byte) (n int, err error) {
	n = len(p)
	s.n += int64(n)
	if s.n > s.limit {
		return 0, ErrMaxSizeFileLimit
	}
	return n, nil
}
