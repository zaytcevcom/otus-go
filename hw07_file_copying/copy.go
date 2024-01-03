package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	if offset > info.Size() {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		limit = info.Size()
	}

	if offset > 0 {
		_, _ = src.Seek(offset, io.SeekStart)
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	bar := pb.Full.Start64(limit)
	_, err = io.CopyN(dst, bar.NewProxyReader(src), limit)
	bar.Finish()

	if err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}

	return nil
}
