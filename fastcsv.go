package fastcsv

import (
	"bufio"
	"bytes"
	"os"
)

type FileReader struct {
	file      *os.File
	scanner   *bufio.Scanner
	separator []byte
}

func NewFileReader(filename string, separator []byte) (*FileReader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &FileReader{
		file:      f,
		scanner:   bufio.NewScanner(f),
		separator: separator,
	}, nil
}

func (r *FileReader) Close() error {
	return r.file.Close()
}

func (r *FileReader) Scan() bool {
	return r.scanner.Scan()
}

func (r *FileReader) Record() [][]byte {
	return bytes.Split(r.scanner.Bytes(), r.separator)
}

func (r *FileReader) Err() error {
	return r.scanner.Err()
}
