package frame

import (
	"bytes"
	"testing"
)

func TestCalculateCRC(t *testing.T) {
	cases := []struct {
		rawFrame []byte
		crc      byte
		name     string
	}{
		{
			rawFrame: []byte("LD+\x00\x00#"),
			crc:      0,
			name:     "Test Case 1",
		},
		{
			rawFrame: []byte("LD+\x00\x05#"),
			crc:      5,
			name:     "Test Case 2",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateCRC(tc.rawFrame)
			want := tc.crc

			if got != want {
				t.Errorf("got %b, want %b", got, want)
				t.Errorf("aka got %d, want %d", got, want)
			}
		})
	}

}

func TestCreateFrame(t *testing.T) {
	t.Run("valid frame 1", func(t *testing.T) {
		want := []byte("LD+\x00\x00#")

		got := CreateRawFrame(0)

		if !bytes.Equal(got, want) {
			t.Errorf("got %b, want %b", got, want)
			t.Errorf("aka got %d, want %d", got, want)
		}
	})

	t.Run("valid frame 2", func(t *testing.T) {
		want := []byte("LD+\x00\x05#")

		got := CreateRawFrame(5)

		if !bytes.Equal(got, want) {
			t.Errorf("got %b, want %b", got, want)
			t.Errorf("got %d, want %d", got, want)
		}
	})
}
