package fastcsv

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"syscall"
)

type FileReader struct {
	file      *os.File
	data      []byte
	current   []byte
	separator byte
	dest      []reflect.Value
}

func NewFileReader(filename string, separator byte, dest interface{}) (*FileReader, error) {
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

	r := &FileReader{
		file:      f,
		data:      data,
		current:   data,
		separator: separator,
	}

	unsetFields := make(map[string]bool)
	destFields := make(map[string]reflect.Value)
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	for i := 0; i < destValue.NumField(); i++ {
		field := destValue.Type().Field(i)

		if name, ok := field.Tag.Lookup("csv"); ok {
			destFields[name] = destValue.Field(i)
		} else {
			destFields[field.Name] = destValue.Field(i)
			unsetFields[field.Name] = true
		}
	}

	header := r.byteRecord()

	r.dest = make([]reflect.Value, len(header))
	for index, name := range header {
		field, ok := destFields[string(name)]
		if !ok {
			continue
		}

		r.dest[index] = field
		delete(unsetFields, string(name))
	}

	if len(unsetFields) > 0 {
		return nil, fmt.Errorf("Not all fields present in csv")
	}

	return r, nil
}

func (r *FileReader) byteRecord() [][]byte {
	eol := bytes.IndexByte(r.current, '\n')
	if eol == -1 {
		return nil
	} else {
		row := r.current[:eol]
		r.current = r.current[eol+1:]
		return bytes.Split(row, []byte{r.separator})
	}
}

func (r *FileReader) Scan() bool {
	if len(r.current) == 0 {
		return false
	}

	nfield := 0
	start := 0
	i := 0
	for {
		if i >= len(r.current) || r.current[i] == '\n' || r.current[i] == r.separator {
			if r.dest[nfield].CanSet() {
				r.dest[nfield].SetBytes(r.current[start:i])
			}
			start = i + 1
			nfield++
		}

		if i >= len(r.current) {
			r.current = nil
			break
		} else if r.current[i] == '\n' {
			r.current = r.current[i+1:]
			break
		}

		i++
	}

	return nfield >= len(r.dest)
}

func (r *FileReader) Close() error {
	var err error

	err = syscall.Munmap(r.data)
	err = r.file.Close()

	return err
}
