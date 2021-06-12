package files

import (
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Local struct {
	maxFileSize int // max bytes of file
	basePath    string
}

func NewLocal(maxFileSize int, basePath string) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	return &Local{maxFileSize: maxFileSize, basePath: p}, nil
}

func (l *Local) Save(fn string, contents io.Reader) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	idString := id.String()
	path := filepath.Join(idString, fn)
	fp := l.fullPath(path)
	// make sure the target directory exists
	d := filepath.Dir(fp)
	err = os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return err
	}
	// create file
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()
	/*
		version 1: just copy the entire contents into the file
		```
		_, err = io.Copy(f, contents)
		if err != nil {
			return err
		}
		```
	*/
	/*
		version 2: use a buffered reader to read off bytes in a loop
	*/
	buf := make([]byte, l.maxFileSize)
	for {
		n, err := contents.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err = f.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

func (l *Local) fullPath(path string) string {
	return filepath.Join(l.basePath, path)
}
