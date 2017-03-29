package fastcsv

import (
	"encoding/csv"
	"io"
	"math/rand"
	"os"
	"testing"
)

var names = []string{
	"Reyes Palmer",
	"Lauran Sander",
	"Oswaldo Dyess",
	"Jamila Tiffany",
	"Shalonda Teti",
	"Monty Alcott",
	"Donald Brand",
	"Reginald Morningstar",
	"Elvie Aguiniga",
	"Cris Mulford",
	"Autumn Dahlquist",
	"Palmer Redman",
	"Merry Lesane",
	"Jannie Laura",
	"Reina Lofland",
	"Norma Valiente",
	"Millard Melville",
}

func createTestFile(filename string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0655)
	if os.IsExist(err) {
		return
	} else if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	for i := 0; i < 1000000; i++ {
		record := make([]string, 5)
		for i := range record {
			j := rand.Intn(len(names))
			record[i] = names[j]
		}
		w.Write(record)
	}

	w.Flush()
	if err := w.Error(); err != nil {
		panic(err)
	}
}

func BenchmarkRead(b *testing.B) {
	createTestFile("test.csv")
	r, err := NewFileReader("test.csv", []byte(","))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Record()
	}
}

func BenchmarkStandard(b *testing.B) {
	createTestFile("test.csv")
	f, err := os.Open("test.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := r.Read()
		if err == io.EOF {
			panic("Ran out of records")
		} else if err != nil {
			panic(err)
		}
	}
}

func TestComparison(t *testing.T) {
	createTestFile("test.csv")

	actual := make(chan [][]byte)
	expected := make(chan []string)
	actualNext := make(chan bool)
	expectedNext := make(chan bool)

	go func() {
		defer close(actual)
		r, err := NewFileReader("test.csv", []byte(","))
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := r.Close(); err != nil {
				panic(err)
			}
		}()

		for {
			actual <- r.Record()
			<-actualNext
		}
	}()
	go func() {
		defer close(expected)
		f, err := os.Open("test.csv")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		r := csv.NewReader(f)

		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}

			expected <- record
			<-expectedNext
		}
	}()

	nrecord := 0
	for {
		a := <-actual
		e := <-expected

		if a == nil && e == nil {
			break
		}

		if len(a) != len(e) {
			t.Fatalf("record %d: len(a) = %d, len(e) = %d", nrecord, len(a), len(e))
		}

		for i := 0; i < len(a); i++ {
			if string(a[i]) != e[i] {
				t.Fatalf("record %d: a[%d] = %q, e[%d] = %q", nrecord, a[i], e[i])
			}
		}

		actualNext <- true
		expectedNext <- true
		nrecord++
	}
}
