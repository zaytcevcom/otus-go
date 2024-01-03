package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	tmpFile, _ := os.CreateTemp("", "test")
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name      string
		fromPath  string
		checkPath string
		offset    int64
		limit     int64
	}{
		{
			name:      "Success offset 0 limit 0",
			fromPath:  "testdata/input.txt",
			checkPath: "testdata/out_offset0_limit0.txt",
			offset:    int64(0),
			limit:     int64(0),
		},
		{
			name:      "Success offset 0 limit 10",
			fromPath:  "testdata/input.txt",
			checkPath: "testdata/out_offset0_limit10.txt",
			offset:    int64(0),
			limit:     int64(10),
		},
		{
			name:      "Success offset 0 limit 1000",
			fromPath:  "testdata/input.txt",
			checkPath: "testdata/out_offset0_limit1000.txt",
			offset:    int64(0),
			limit:     int64(1000),
		},
		{
			name:      "Success offset 0 limit 10000",
			fromPath:  "testdata/input.txt",
			checkPath: "testdata/out_offset0_limit10000.txt",
			offset:    int64(0),
			limit:     int64(10000),
		},
		{
			name:      "Success offset 100 limit 1000",
			fromPath:  "testdata/input.txt",
			checkPath: "testdata/out_offset100_limit1000.txt",
			offset:    int64(100),
			limit:     int64(1000),
		},
		{
			name:      "Success offset 6000 limit 1000",
			fromPath:  "testdata/input.txt",
			checkPath: "testdata/out_offset6000_limit1000.txt",
			offset:    int64(6000),
			limit:     int64(1000),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(tc.fromPath, tmpFile.Name(), tc.offset, tc.limit)
			assert.NoError(t, err)

			toInfo, _ := os.Stat(tmpFile.Name())
			checkInfo, _ := os.Stat(tc.checkPath)

			assert.Equal(t, toInfo.Size(), checkInfo.Size())
		})
	}
}

func TestCopyFailed(t *testing.T) {
	tmpFile, _ := os.CreateTemp("", "test")
	defer os.Remove(tmpFile.Name())

	t.Run("Offset exceeds file size (empty file)", func(t *testing.T) {
		fromPath := "testdata/empty.txt"
		toPath := tmpFile.Name()
		offset := int64(40)
		limit := int64(0)

		err := Copy(fromPath, toPath, offset, limit)
		assert.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})

	t.Run("Offset exceeds file size", func(t *testing.T) {
		fromPath := "testdata/input.txt"
		toPath := tmpFile.Name()
		offset := int64(1000000)
		limit := int64(0)

		err := Copy(fromPath, toPath, offset, limit)
		assert.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})
}
