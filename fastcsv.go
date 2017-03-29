package fastcsv

import (
	"bytes"
	"os"
	"syscall"
)

type FileReader struct {
	file      *os.File
	data      []byte
	current   []byte
	separator []byte
}

func NewFileReader(filename string, separator []byte) (*FileReader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, int(fi.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	return &FileReader{
		file:      f,
		data:      data,
		current:   data,
		separator: separator,
	}, nil
}

func (r *FileReader) Record() [][]byte {
	eol := bytes.IndexByte(r.current, '\n')
	if eol == -1 {
		return nil
	} else {
		row := r.current[:eol]
		r.current = r.current[eol+1:]
		return bytes.Split(row, r.separator)
	}
}

func (r *FileReader) Close() error {
	var err error

	err = syscall.Munmap(r.data)
	err = r.file.Close()

	return err
}
