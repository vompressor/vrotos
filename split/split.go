package split

import (
	"bytes"
	"io"
)

type SplitReader struct {
	r    io.Reader
	buf  []byte
	size int
}

func NewSplitReader(source io.Reader, splitSize int) *SplitReader {
	s := &SplitReader{}
	s.buf = make([]byte, 0)
	s.r = source
	s.size = splitSize
	return s
}

func (sr *SplitReader) GetReader() (io.Reader, error) {
	n, err := io.ReadFull(sr.r, sr.buf)
	if err != nil {
		if n > 0 {
			return bytes.NewReader(sr.buf[:n]), nil
		}
		return nil, err
	}
	return bytes.NewReader(sr.buf[:n]), nil
}

func (sr *SplitReader) GetCopiedReader() (io.Reader, error) {
	n, err := io.ReadFull(sr.r, sr.buf)
	if err != nil {
		if n > 0 {
			buf := make([]byte, n)
			copy(buf, sr.buf[:n])
			return bytes.NewReader(buf), nil
		}
		return nil, err
	}
	buf := make([]byte, n)
	copy(buf, sr.buf[:n])
	return bytes.NewReader(buf), nil
}

type SplitWriter struct {
	buf       *bytes.Buffer
	size      int
	setWriter func() io.Writer
}

func NewSplitWriter(f func() io.Writer, size int) *SplitWriter {
	sw := SplitWriter{}
	sw.buf = bytes.NewBuffer(make([]byte, 0))

	sw.setWriter = f
	sw.size = size
	return &sw
}

func (sw *SplitWriter) Write(b []byte) (int, error) {
	written := 0
	br := bytes.NewReader(b)
	if sw.buf.Len() != 0 {
		m, _ := io.CopyN(sw.buf, br, int64(sw.size-sw.buf.Len()))
		written += int(m)
		if sw.buf.Len() >= sw.size {
			io.CopyN(sw.setWriter(), sw.buf, int64(sw.size))
			sw.buf.Reset()
		}
	}
	for {

		if br.Len() >= sw.size {
			m, _ := io.CopyN(sw.setWriter(), br, int64(sw.size))
			written += int(m)
		} else {
			if br.Len() == 0 {
				return written, nil
			}
			m, _ := io.Copy(sw.buf, br)
			written += int(m)

		}
	}

}

func (sw *SplitWriter) Flush() (err error) {
	err = nil
	if sw.buf.Len() > 0 {
		_, err = io.Copy(sw.setWriter(), sw.buf)
	}
	return
}

func (sw *SplitWriter) Close() error {
	return sw.Flush()
}
