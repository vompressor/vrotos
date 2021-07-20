package async_rw

import (
	"io"
	"sync"
)

type AsyncWriter struct {
	w  io.Writer
	m  *sync.Mutex
	wg *sync.WaitGroup
}

func NewAsyncWriter(w io.Writer) *AsyncWriter {
	aw := AsyncWriter{}
	aw.m = new(sync.Mutex)
	aw.w = w
	aw.wg = new(sync.WaitGroup)
	return &aw
}

func (aw *AsyncWriter) Write(b []byte) (n int, err error) {
	aw.wg.Add(1)
	go func() {
		aw.m.Lock()
		n, err = aw.w.Write(b)
		aw.m.Unlock()
		aw.wg.Done()
	}()
	aw.wg.Wait()
	return
}

func (aw *AsyncWriter) Asyncwrite(b []byte, callback RWCallback) {
	go func() {
		aw.m.Lock()
		callback(aw.w.Write(b))
		aw.m.Unlock()
	}()
}

// func (aw *AsyncWriter) ReadFrom(r io.Reader) (n int64, err error) {
// 	return io.Copy(aw.w, r)
// }

// func (aw *AsyncWriter) AsyncReadFrom(r io.Reader, callback CopyCallback) {
// 	go func() {
// 		aw.m.Lock()
// 		callback(aw.ReadFrom(r))
// 		aw.m.Unlock()
// 	}()
// }
