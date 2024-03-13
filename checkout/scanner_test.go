package checkout

import (
	"io"
	"strings"
	"testing"
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
