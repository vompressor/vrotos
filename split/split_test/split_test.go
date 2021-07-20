package split_test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/vompressor/vrotos/split"
)

const path = "./testBigData.txt"

func TestSplit(t *testing.T) {
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err.Error())
	}

	sr := split.NewSplitReader(f, 2048)

	for {
		r, err := sr.GetReader()
		if err != nil {
			t.Log(err.Error())
			return
		}

		b, _ := io.ReadAll(r)
		fmt.Printf("%s - %d\n", b, len(b))
	}

}

func TestSplit2(t *testing.T) {
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err.Error())
	}

	sr := split.NewSplitReader(f, 2048)

	for {
		r, err := sr.GetCopiedReader()
		if err != nil {
			t.Log(err.Error())
			return
		}

		b, _ := io.ReadAll(r)
		fmt.Printf("{%s} - %d\n", b, len(b))

	}

}

func TestMakeBigFile(t *testing.T) {
	f, _ := os.Create(path)

	i := 0

	w := bufio.NewWriter(f)

	for {
		if i > 5000 {
			break
		}
		w.WriteString(fmt.Sprintf("%d ", i))
		i++
	}
	w.Flush()
}

func TestSplitWrite(t *testing.T) {
	os.Mkdir("o", os.ModePerm)
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err.Error())
	}
	count := 0
	sw := split.NewSplitWriter(func() io.Writer {
		o, err := os.Create(fmt.Sprintf("o/testBigData.txt.%c", 'a'+count))
		if err != nil {
			t.Log(err.Error())
		}
		count++
		return o
	}, 1024*10)
	io.Copy(sw, f)
	sw.Flush()
}
