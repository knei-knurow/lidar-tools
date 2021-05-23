package frame

import (
	"bytes"
	"testing"
)

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
			got := CalculateCRC(tc.data)
			want := tc.crc

			if got != want {
				t.Errorf("got %s, want %s", DescribeByte(got), DescribeByte(want))
			}
		})
	}

}

func TestCreateFrame(t *testing.T) {
	t.Run("valid frame 1", func(t *testing.T) {
		want := []byte("LD+\x00\x00#")

		got := EncodeRawFrame(0)

		if !bytes.Equal(got, want) {
			t.Errorf("got %b, want %b", got, want)
			t.Errorf("aka got %d, want %d", got, want)
		}
	})
}
