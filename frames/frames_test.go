package frames_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/knei-knurow/lidar-tools/frames"
)

var testCases = []struct {
	header   []byte
	data     []byte
	checksum byte
}{
	{
		header:   []byte{'L', 'D'},
		data:     []byte{},
		checksum: 0x00,
	},
	{
		header:   []byte{'L', 'D'},
		data:     []byte{'A'},
		checksum: 0x41,
	},
	{
		header:   []byte{'L', 'D'},
		data:     []byte{'t', 'e', 's', 't'},
		checksum: 0x16,
	},
	{
		header:   []byte{'L', 'D'},
		data:     []byte{'d', 'u', 'p', 'c', 'i', 'a'},
		checksum: 0x0a,
	},
	{
		header:   []byte{'L', 'D'},
		data:     []byte{'l', 'o', 'l', 'x', 'd'},
		checksum: 0x73,
	},
	{
		header:   []byte{'B', 'I', 'G'},
		data:     []byte{'d', 'o', 'n', 'd', 'u'},
		checksum: 0x30,
	},
	{
		header:   []byte{},
		data:     []byte{},
		checksum: 0x08,
	},
}

func TestCreate(t *testing.T) {
	for i, tc := range testCases {
		testName := fmt.Sprintf("test %d", i)
		t.Run(testName, func(t *testing.T) {
			gotFrame := frames.Create(tc.header, tc.data)

			if !bytes.Equal(gotFrame.Header(), tc.header) {
				t.Errorf("got header % x, want header % x", gotFrame.Header(), tc.header)
			}

			if !bytes.Equal(gotFrame.Data(), tc.data) {
				t.Errorf("got data % x, want data % x", gotFrame.Data(), tc.data)
			}

			if gotFrame.Checksum() != tc.checksum {
				t.Errorf("got checksum % x, want checksum % x", gotFrame.Checksum(), tc.checksum)
			}
		})
	}
}

func TestAssemble(t *testing.T) {
	for i, tc := range testCases {
		testName := fmt.Sprintf("test %d", i)
		t.Run(testName, func(t *testing.T) {
			gotFrame := frames.Assemble(tc.header, tc.data, tc.checksum)

			if !bytes.Equal(gotFrame.Header(), tc.header) {
				t.Errorf("got header % x, want header % x", gotFrame.Header(), tc.header)
			}

			if !bytes.Equal(gotFrame.Data(), tc.data) {
				t.Errorf("got data % x, want data % x", gotFrame.Data(), tc.data)
			}

			if gotFrame.Checksum() != tc.checksum {
				t.Errorf("got checksum % x, want checksum % x", gotFrame.Checksum(), tc.checksum)
			}
		})
	}
}
