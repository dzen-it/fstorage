package storage

import (
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/dzen-it/fstorage/utils/log"
)

const (
	cleaningInterval = 60 * 5 // 5min in sec
)

type File struct {
	Filename   string
	AccessTime int64
	Size       int64
}

type Files []*File

func (f Files) Len() int      { return len(f) }
func (f Files) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

type ByAccessTime struct{ Files }

func (f ByAccessTime) Less(i, j int) bool { return f.Files[i].AccessTime < f.Files[j].AccessTime }

func cleanerWorker(dir string, limit int64) {
	for {
		if space := detectLimitMemmory(dir, limit); space > 0 {
			files := getListFiles(dir, space)
			cleanFiles(dir, files)
		}

		time.Sleep(time.Second * time.Duration(cleaningInterval))
	}
}

func getListFiles(dir string, freeSpace int64) (files []*File) {
	var sum int64

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		fstat, ok := info.Sys().(*syscall.Stat_t)
		if ok {
			sec, _ := fstat.Atim.Unix()
			files = append(files, &File{
				Filename:   info.Name(),
				AccessTime: sec,
				Size:       info.Size(),
			})

			sum += info.Size()

			for {
				if sum <= freeSpace {
					break
				}

				sort.Sort(ByAccessTime{files})

				lastIndex := len(files) - 1
				sum -= files[lastIndex].Size
				files = files[:lastIndex]
			}
		}
		return nil
	})

	return
}

func cleanFiles(dir string, files []*File) {
	var (
		err  error
		path string
	)

	for _, f := range files {
		path = filepath.Join(dir, f.Filename[:2], f.Filename)
		if err = os.Remove(path); err != nil {
			log.Errorw("Remove file", "path", path, "error", err)
		}
	}
}

func detectLimitMemmory(dir string, limit int64) (delta int64) {
	size := dirSize(dir)
	if size > limit {
		return size - limit
	}
	return 0
}

func dirSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size
}
