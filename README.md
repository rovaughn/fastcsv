fastcsv
=======

CSV parser for golang that aims to be as fast as possible through the use of
mmap and happens to provide a convenient interface.

Example usage:

    var record struct {
        Name []byte `csv:"name"`
        Age  []byte `csv:"age"`
        Sex  []byte `csv:"sex"`
    }

    reader, err := NewFileReader("file.csv", ',', &record)
    if err != nil {
        panic(err)
    }
    defer reader.Close()

    for reader.Scan() {
        fmt.Printf("%s|%s|%s\n", record.Name, record.Age, record.Sex)
    }

This setup currently gets 3.01 ns/row on my laptop, whereas the standard library
gets 1872 ns/row.  This can be tested with `go test -bench=.`.

**BUG**: Currently the reader does not implemented the csv spec correctly yet;
namely it does not handle quoted fields or escaped characters.  It also always
assumed rows are terminated by LF instead of CRLF.
