package compressor_test

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/vompressor/vrotos/async_rw"
	"github.com/vompressor/vrotos/compressor"
	"github.com/vompressor/vrotos/split"
)

func TestGetFiles(t *testing.T) {
	os.Mkdir("o", os.ModePerm)
	pack := compressor.NewPack(".")
	pack.AddRelItem("verybigdata.txt")

	for i := range pack.Items {
		println(i)
	}
}
func TestCompFiles(t *testing.T) {
	os.Mkdir("o", os.ModePerm)
	pack := compressor.NewPack(".")
	pack.AddRelItem("verybigdata.txt")

	for i := range pack.Items {
		println(i)
	}

	f, _ := os.Create("o/verybigdata.txt.tar")

	pack.WriteTar(f)
}

func TestCompGzFiles(t *testing.T) {
	os.Mkdir("o", os.ModePerm)
	pack := compressor.NewPack(".")
	pack.AddRelItem("verybigdata.txt")

	for i := range pack.Items {
		println(i)
	}

	f, _ := os.Create("o/verybigdata.txt.tgz")
	gzw := gzip.NewWriter(f)
	pack.WriteTar(gzw)
}

func TestSplitComp(t *testing.T) {
	os.Mkdir("o", os.ModePerm)
	pack := compressor.NewPack(".")
	pack.AddItem("verybigdata.txt")

	j := 0
	files := make([]*os.File, 0)
	sw := split.NewSplitWriter(func() io.Writer {
		f, err := os.Create(fmt.Sprintf("o/verybigdata.txt.tar.%c", 'a'+j))
		files = append(files, f)
		if err != nil {
			t.Fatal(err.Error())
		}
		j++
		return f
	}, 1024*1024)

	err := pack.WriteTar(sw)
	if err != nil {
		t.Fatal(err.Error())
	}
	sw.Flush()
	sw.Close()
	for _, n := range files {
		n.Close()
	}
}

func TestSplitCompGz(t *testing.T) {
	os.Mkdir("o", os.ModePerm)
	pack := compressor.NewPack(".")
	pack.AddItem("verybigdata.txt")

	j := 0
	files := make([]*os.File, 0)
	sw := split.NewSplitWriter(func() io.Writer {
		f, err := os.Create(fmt.Sprintf("o/verybigdata.txt.tgz.%c", 'a'+j))
		files = append(files, f)
		if err != nil {
			t.Fatal(err.Error())
		}
		j++
		return f
	}, 1024*1024)
	gzw := gzip.NewWriter(sw)
	err := pack.WriteTar(gzw)
	if err != nil {
		t.Fatal(err.Error())
	}
	gzw.Close()
	sw.Flush()
	sw.Close()
	for _, n := range files {
		n.Close()
	}
}

func TestAsyncSplitComp(t *testing.T) {
	os.Mkdir("o", os.ModePerm)
	pack := compressor.NewPack(".")
	pack.AddItem("verybigdata.txt")

	j := 0
	files := make([]*os.File, 0)
	sw := split.NewSplitWriter(func() io.Writer {
		f, err := os.Create(fmt.Sprintf("o/verybigdata.txt.tar.%c", 'a'+j))
		files = append(files, f)
		if err != nil {
			t.Fatal(err.Error())
		}
		j++

		return async_rw.NewAsyncWriter(f)
	}, 1024*1024)

	err := pack.WriteTar(sw)
	if err != nil {
		t.Fatal(err.Error())
	}
	sw.Flush()
	sw.Close()
	for _, n := range files {
		n.Close()
	}
}
