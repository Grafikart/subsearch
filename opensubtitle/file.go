package opensubtitle

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	ChunkSize = 65536 // 64k
)

func HashFile(file *os.File) (h string, err error) {
	fi, err := file.Stat()
	if err != nil {
		return
	}
	if fi.Size() < ChunkSize {
		return "", fmt.Errorf("file is too small")
	}

	// Read head and tail blocks.
	buf := make([]byte, ChunkSize*2)
	err = readChunk(file, 0, buf[:ChunkSize])
	if err != nil {
		return
	}
	err = readChunk(file, fi.Size()-ChunkSize, buf[ChunkSize:])
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

	return fmt.Sprintf("%016x", hash+uint64(fi.Size())), nil
}

func readChunk(file *os.File, offset int64, buf []byte) (err error) {
	n, err := file.ReadAt(buf, offset)
	if err != nil {
		return
	}
	if n != ChunkSize {
		return fmt.Errorf("invalid read %v", n)
	}
	return
}
