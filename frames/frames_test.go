package frames_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/knei-knurow/lidar-tools/frames"
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		inputHeader      []byte
		inputData        []byte
		expectedChecksum byte
	}{
		{
			inputHeader:      []byte{'L', 'D'},
			inputData:        []byte{},
			expectedChecksum: 0x00,
		},
		{
			inputHeader:      []byte{'L', 'D'},
			inputData:        []byte{'A'},
			expectedChecksum: 0x41,
		},
		{
			inputHeader:      []byte{'L', 'D'},
			inputData:        []byte{'t', 'e', 's', 't'},
			expectedChecksum: 0x16,
		},
		{
			inputHeader:      []byte{'L', 'D'},
			inputData:        []byte{'d', 'u', 'p', 'c', 'i', 'a'},
			expectedChecksum: 0x0a,
		},
		{
			inputHeader:      []byte{'L', 'D'},
			inputData:        []byte{'l', 'o', 'l', 'x', 'd'},
			expectedChecksum: 0x73,
		},
		{
			inputHeader:      []byte{'B', 'I', 'G'},
			inputData:        []byte{'d', 'o', 'n', 'd', 'u'},
			expectedChecksum: 0x30,
		},
		{
			inputHeader:      []byte{},
			inputData:        []byte{},
			expectedChecksum: 0x08,
		},
	}

	for i, tc := range testCases {
		testName := fmt.Sprintf("test %d", i)
		t.Run(testName, func(t *testing.T) {
			gotFrame := frames.Create(tc.inputHeader, tc.inputData)

			if !bytes.Equal(gotFrame.Header(), tc.inputHeader) {
				t.Errorf("got header % x, want header % x", gotFrame.Header(), tc.inputHeader)
			}

			if !bytes.Equal(gotFrame.Data(), tc.inputData) {
				t.Errorf("got data % x, want data % x", gotFrame.Data(), tc.inputData)
			}

			if gotFrame.Checksum() != tc.expectedChecksum {
				t.Errorf("got checksum % x, want checksum % x", gotFrame.Checksum(), tc.expectedChecksum)
			}
		})
	}
}

func TestAssemble(t *testing.T) {
	testCases := []struct {
		inputHeader   []byte
		inputData     []byte
		inputChecksum byte
	}{
		{
			inputHeader:   []byte{'L', 'D'},
			inputData:     []byte{},
			inputChecksum: 0x00,
		},
		{
			inputHeader:   []byte{'L', 'D'},
			inputData:     []byte{'A'},
			inputChecksum: 0x41,
		},
		{
			inputHeader:   []byte{'L', 'D'},
			inputData:     []byte{'t', 'e', 's', 't'},
			inputChecksum: 0x16,
		},
		{
			inputHeader:   []byte{'L', 'D'},
			inputData:     []byte{'d', 'u', 'p', 'c', 'i', 'a'},
			inputChecksum: 0x0a,
		},
		{
			inputHeader:   []byte{'L', 'D'},
			inputData:     []byte{'l', 'o', 'l', 'x', 'd'},
			inputChecksum: 0x73,
		},
		{
			inputHeader:   []byte{'B', 'I', 'G'},
			inputData:     []byte{'d', 'o', 'n', 'd', 'u'},
			inputChecksum: 0x30,
		},
		{
			inputHeader:   []byte{},
			inputData:     []byte{},
			inputChecksum: 0x08,
		},
	}

	for i, tc := range testCases {
		testName := fmt.Sprintf("test %d", i)
		t.Run(testName, func(t *testing.T) {
			gotFrame := frames.Create(tc.inputHeader, tc.inputData)

			if !bytes.Equal(gotFrame.Header(), tc.inputHeader) {
				t.Errorf("got header % x, want header % x", gotFrame.Header(), tc.inputHeader)
			}

			if !bytes.Equal(gotFrame.Data(), tc.inputData) {
				t.Errorf("got data % x, want data % x", gotFrame.Data(), tc.inputData)
			}

			if gotFrame.Checksum() != tc.inputChecksum {
				t.Errorf("got checksum % x, want checksum % x", gotFrame.Checksum(), tc.inputChecksum)
			}
		})
	}
}
