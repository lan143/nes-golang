package rom

import "os"

type Rom interface {
	Load(file *os.File) error
	GetData() []byte
}
