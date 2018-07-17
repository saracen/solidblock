package solidblock_test

import (
	"bytes"
	"compress/gzip"
	"hash/crc32"
	"io"
	"os"

	"github.com/saracen/solidblock"
)

func ExampleSolidblock() {
	// file contents
	files := [][]byte{
		[]byte("file 1\n"),
		[]byte("file 2\n"),
	}

	// file metadata
	var metadata struct {
		sizes []uint64
		crcs  []uint32
	}
	metadata.sizes = []uint64{
		uint64(len(files[0])),
		uint64(len(files[1])),
	}
	metadata.crcs = []uint32{
		crc32.ChecksumIEEE(files[0]),
		crc32.ChecksumIEEE(files[1]),
	}

	// Concatenate files to compressed block
	block := new(bytes.Buffer)
	w := gzip.NewWriter(block)
	w.Write(files[0])
	w.Write(files[1])
	w.Close()

	// Open gzip reader to compressed block
	r, err := gzip.NewReader(block)
	if err != nil {
		panic(err)
	}

	// Create a new solidblock reader
	s := solidblock.New(r, metadata.sizes, metadata.crcs)

	for {
		err := s.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		io.Copy(os.Stdout, s)
	}

	// Output:
	// file 1
	// file 2
}
