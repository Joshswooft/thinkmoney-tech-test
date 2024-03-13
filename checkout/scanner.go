package checkout

import (
	"errors"
	"io"

	"github.com/Joshswooft/thinkmoney-test/sku"
)

var errNilReaderProvided = errors.New("nil reader provided")

// Scanner reads product information and extracts the product meta data
// reads in SKU's from a given reader this could be a buffer, file, string etc.
type skuScanner struct {
	reader io.Reader
}

func NewSkuScanner(reader io.Reader) (*skuScanner, error) {
	if reader == nil {
		return nil, errNilReaderProvided
	}

	return &skuScanner{
		reader: reader,
	}, nil
}

// Implements the io.Reader so it can be combined with other readers
// reads in skus from a given byte array
func (s *skuScanner) Read(p []byte) (int, error) {

	bytesRead, err := s.reader.Read(p)
	if err != nil {
		return bytesRead, err
	}

	buf := make([]byte, bytesRead)

	for i := 0; i < bytesRead; i++ {
		sku, err := sku.New(rune(p[i]))

		if err == nil {
			buf[i] = byte(sku.Value())
		}
	}

	copy(p, buf)

	return bytesRead, nil
}
