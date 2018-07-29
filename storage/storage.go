package storage

import (
	"io"
	"os"
	"path/filepath"
)

const (
	ctxStorage = "storage"
)

type Error uint16

func (e Error) Error() string {
	return ""
}

const (
	ErrFileNotFound Error = iota
)

type Storage interface {
	SaveFile(hash string, reader io.Reader) error
	GetFile(hash string) (io.Reader, error)
	DeleteFile(hash string) error
}

type fileStorage struct {
	limit int64
	path  string
}

func NewFileStorage(dir string, limit int64) (f *fileStorage, err error) {
	f = new(fileStorage)
	f.path = dir

	if _, err = os.Stat(f.path); os.IsNotExist(err) {
		os.Mkdir(f.path, os.ModePerm)
		err = nil
	}

	return f, err
}

func (f *fileStorage) SaveFile(hash string, reader io.Reader) error {
	subdir := hash[:2]
	endpath := filepath.Join(f.path, subdir)
	err := os.MkdirAll(endpath, os.ModePerm)
	if err != nil {
		return err
	}

	filename := filepath.Join(endpath, hash)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

func (f *fileStorage) GetFile(hash string) (io.Reader, error) {
	subdir := hash[:2]
	filename := filepath.Join(f.path, subdir, hash)
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil, ErrFileNotFound
	}

	return file, err
}

func (f *fileStorage) DeleteFile(hash string) error {
	subdir := hash[:2]
	filename := filepath.Join(f.path, subdir, hash)
	err := os.Remove(filename)
	if os.IsNotExist(err) {
		return ErrFileNotFound
	}
	return err
}

func createDirIfnotExists(path string) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
		return nil
	}
	return err
}
