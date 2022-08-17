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
	sourceFileStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	source, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	if offset != 0 {
		if offset >= sourceFileStat.Size() {
			return ErrOffsetExceedsFileSize
		}
		source.Seek(offset, 0)
	}

	realSize := sourceFileStat.Size() - offset
	barLimit := realSize
	if limit != 0 {
		barLimit = limit
		if barLimit > realSize {
			barLimit = realSize
		}
	}
	bar := pb.Full.Start64(barLimit)
	limitReader := io.LimitReader(source, barLimit)
	barReader := bar.NewProxyReader(limitReader)

	var errCopy error
	_, errCopy = io.Copy(destination, barReader)
	bar.Finish()

	if errCopy != nil {
		return errCopy
	}
	return nil
}
