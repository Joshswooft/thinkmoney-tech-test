package checkout

import (
	"errors"
	"io"
	"log"

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

// scans a single sku, for use wrap in an infinite loop and scan until io.EOF or internal err
func (s *skuScanner) Scan() (sku.SKU, error) {
	var b [1]byte
	_, err := s.Read(b[:])
	if err != nil {
		return sku.SKU{}, err
	}

	// Validate and create SKU directly from byte
	skuInstance, err := sku.New(rune(b[0]))
	if err != nil {
		log.Printf("error reading in sku: err: %v, b: '%q'", err, string(b[0]))
		return sku.SKU{}, err
	}

	return skuInstance, nil
}
