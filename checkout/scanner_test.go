package checkout

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/Joshswooft/thinkmoney-test/sku"
)

const nullCharacter = '\x00'

// Removes the null character
func removeNull(r rune) rune {
	if r == nullCharacter {
		return -1
	}
	return r
}

func TestStringScanner_Read(t *testing.T) {
	type fields struct {
		reader io.Reader
	}

	file, err := os.Open("./testdata/skus.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	tests := []struct {
		name       string
		fields     fields
		wantPos    int
		want       string
		err        error
		bufferSize int
	}{
		{
			name: "reads in empty string",
			fields: fields{
				reader: strings.NewReader(""),
			},
			wantPos: 0,
			err:     io.EOF,
		},
		{
			name:       "reads in non-alphabetical string",
			fields:     fields{reader: strings.NewReader("6&%")},
			bufferSize: 3,
			wantPos:    3,
			want:       "",
			err:        nil,
		},
		{
			name:       "reads in string removing special characters and transforming to upper case",
			fields:     fields{reader: strings.NewReader("a69B420Cd")},
			bufferSize: 9,
			want:       "ABCD",
			wantPos:    9,
			err:        nil,
		},
		{
			name:       "reads data in from a file",
			fields:     fields{reader: file},
			bufferSize: 20,
			want:       "GUMMYBAR",
			wantPos:    10,
			err:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &skuScanner{
				reader: tt.fields.reader,
			}
			buf := make([]byte, tt.bufferSize)

			n, err := s.Read(buf)
			if err != tt.err {
				t.Errorf("StringScanner.Read() error = %v, wantErr %v", err, tt.err)
				return
			}
			if n != tt.wantPos {
				t.Errorf("StringScanner.Read() pos = %v, wantPos %v", n, tt.wantPos)
			}

			got := strings.Map(removeNull, string(buf))

			if got != tt.want {
				t.Errorf("scanner results are incorrect! got: '%q', want: '%q'", got, tt.want)
			}
		})
	}
}

func TestSkuScanner_Scan(t *testing.T) {

	skuGen := func(r rune) sku.SKU {
		return skuGenerator(t, r)
	}

	file, err := os.Open("./testdata/skus.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	type fields struct {
		reader io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		want   sku.SKU
		err    error
	}{
		{
			name:   "should return valid sku",
			fields: fields{reader: strings.NewReader("z")},
			want:   skuGen('Z'),
			err:    nil,
		},
		{
			name:   "should return a sku error",
			fields: fields{reader: strings.NewReader("69")},
			want:   sku.SKU{},
			err:    sku.ErrNoSpecialCharacters,
		},
		{
			name:   "should return an empty sku when given an empty input",
			fields: fields{reader: strings.NewReader("")},
			want:   sku.SKU{},
			err:    io.EOF,
		},
		{
			name:   "reads in first sku from a file",
			fields: fields{reader: file},
			want:   skuGen('G'),
			err:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewSkuScanner(tt.fields.reader)

			if err != nil {
				t.Errorf("failed to initialize scanner, error = %v", err)
			}

			got, err := s.Scan()
			if err != tt.err {
				t.Errorf("SkuScanner.Scan() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SkuScanner.Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}
