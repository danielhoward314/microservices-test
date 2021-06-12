package files

import "io"

// interface abstracts behavior so implementation aside from `local` can be created
type Storage interface {
	Save(path string, file io.Reader) error
}
