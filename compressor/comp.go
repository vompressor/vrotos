package compressor

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Item struct {
	RelPath string
	AbsPath string
}

type Pack struct {
	RootPath string
	Items    map[string]*Item
}

func NewPack(root string) *Pack {
	return &Pack{RootPath: root, Items: make(map[string]*Item)}
}

func (p *Pack) AddRelItem(relPath string) error {
	return p.AddItem(filepath.Join(p.RootPath, relPath))
}

func (p *Pack) AddItem(absPath string) error {
	f, err := os.Open(absPath)
	if err != nil {
		return err
	}

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if fi.IsDir() {
		f.Close()
		filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			p.add(path, info)
			return nil
		})
		return nil
	}

	return p.add(absPath, fi)
}

func (p *Pack) add(path string, fi os.FileInfo) error {
	if fi.IsDir() {
		return nil
	}
	relPath, err := filepath.Rel(p.RootPath, path)
	if strings.Contains(relPath, "..") {
		return errors.New("don't add path \"..\"")
	}
	if err != nil {
		return err
	}

	i := &Item{}
	i.RelPath = relPath
	i.AbsPath = path

	p.Items[i.RelPath] = i
	return nil
}

func (p *Pack) WriteTar(w io.Writer) error {
	tw := tar.NewWriter(w)

	for _, m := range p.Items {
		println(m.AbsPath)
		f, err := os.Open(m.AbsPath)
		if err != nil {
			println(err.Error())
			continue
		}

		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			println(err.Error())
			continue
		}
		header, err := tar.FileInfoHeader(fi, m.AbsPath)
		if err != nil {
			println(err.Error())
			continue
		}

		header.Name = filepath.ToSlash(m.RelPath)
		tw.WriteHeader(header)
		_, err = io.Copy(tw, f)
		if err != nil {
			println(err.Error())
			continue
		}
	}
	tw.Close()
	return nil
}
