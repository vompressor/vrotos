package async_rw

import (
	"io"
	"sync"
)

type AsyncReader struct {
	r  io.Reader
	m  *sync.Mutex
	wg *sync.WaitGroup
}

func NewAsyncReader(r io.Reader) *AsyncReader {
	ar := AsyncReader{}
	ar.m = new(sync.Mutex)
	ar.wg = new(sync.WaitGroup)
	ar.r = r
	return &ar
}

func (ar *AsyncReader) Read(b []byte) (n int, err error) {
	ar.wg.Add(1)
	go func() {
		ar.m.Lock()
		n, err = ar.r.Read(b)
		ar.m.Unlock()
		ar.wg.Done()
	}()
	ar.wg.Wait()
	return n, err
}

func (ar *AsyncReader) AsyncRead(b []byte, callback RWCallback) {
	go func() {
		ar.m.Lock()
		callback(ar.r.Read(b))
		ar.m.Unlock()
	}()
}

func (ar *AsyncReader) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, ar.r)
}

func (ar *AsyncReader) AsyncWriteTo(w io.Writer, callback CopyCallback) {
	go func() {
		ar.m.Lock()
		ar.m.Lock()
		n, err := ar.WriteTo(w)
		callback(n, err)
		ar.m.Unlock()
	}()
}
