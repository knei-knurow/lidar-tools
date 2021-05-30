package frames

import (
	"testing"
)

// FIXME: fix tests
func TestCalculateCRC(t *testing.T) {
	cases := []struct {
		data []byte
		crc  byte
		name string
	}{
		{
			data: []byte("0"),
			crc:  48,
			name: "Test Case 1",
		},
		{
			data: []byte("01"),
			crc:  1,
			name: "Test Case 2",
		},
		{
			data: []byte("ABC"),
			crc:  64,
			name: "Test Case 3",
		},
		{
			data: []byte{1, 2, 3, 4, 5},
			crc:  1,
			name: "Test Case 4",
		},
		{
			data: []byte{123, 153, 223},
			crc:  61,
			name: "Test Case 5",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateChecksum(tc.data)
			want := tc.crc

			if got != want {
				t.Errorf("got %s, want %s", DescribeByte(got), DescribeByte(want))
			}
		})
	}
}
