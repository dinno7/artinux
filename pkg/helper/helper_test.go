package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMBToBytes(t *testing.T) {
	tests := []struct {
		mb   int64
		want int64
	}{
		{0, 0},
		{1, 1024 * 1024},
		{10, 10 * 1024 * 1024},
		{100, 100 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.want, MBToBytes(tt.mb))
		})
	}
}

func TestIsValidOSAndArch(t *testing.T) {
	tests := []struct {
		os   string
		arch string
		want bool
	}{
		{"linux", "amd64", true},
		{"linux", "arm64", true},
		{"linux", "386", true},
		{"darwin", "amd64", true},
		{"darwin", "arm64", true},
		{"windows", "amd64", true},
		{"windows", "arm64", true},
		{"Linux", "AMD64", true},
		{"LINUX", "ARM64", true},
		{"unknown", "amd64", false},
		{"", "amd64", false},
		{"linux", "m68k", false},
		{"linux", "", false},
		{"darwin", "386", false},
	}

	for _, tt := range tests {
		t.Run(tt.os+"/"+tt.arch, func(t *testing.T) {
			got := IsValidOSAndArch(tt.os, tt.arch)
			assert.Equal(t, tt.want, got)
		})
	}
}
