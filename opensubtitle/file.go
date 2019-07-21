package opensubtitle

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	ChunkSize = 65536 // 64k
)

type hashableFile interface {
	io.ReaderAt
	Size() int64
}

func hashFile(f hashableFile) (h string, err error) {
	if f.Size() < ChunkSize {
		return "", fmt.Errorf("file is too small")
	}

	// Read head and tail blocks.
	buf := make([]byte, ChunkSize*2)
	err = readChunk(f, 0, buf[:ChunkSize])
	if err != nil {
		return
	}
	err = readChunk(f, f.Size()-ChunkSize, buf[ChunkSize:])
	if err != nil {
		return
	}

	// Convert to uint64, and sum.
	var nums [(ChunkSize * 2) / 8]uint64
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &nums)
	if err != nil {
		return "", err
	}
	var hash uint64
	for _, num := range nums {
		hash += num
	}

	return fmt.Sprintf("%016x", hash+uint64(f.Size())), nil
}

func readChunk(file io.ReaderAt, offset int64, buf []byte) (err error) {
	n, err := file.ReadAt(buf, offset)
	if err != nil {
		return
	}
	if n != ChunkSize {
		return fmt.Errorf("invalid read %v", n)
	}
	return
}
