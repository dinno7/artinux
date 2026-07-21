package helper

import (
	"slices"
	"strings"
)

var validOSArch = map[string][]string{
	"aix":       {"ppc64"},
	"android":   {"386", "amd64", "arm", "arm64"},
	"darwin":    {"amd64", "arm64"},
	"dragonfly": {"amd64"},
	"freebsd":   {"386", "amd64", "arm", "arm64", "riscv64"},
	"illumos":   {"amd64"},
	"ios":       {"arm64"},
	"js":        {"wasm"},
	"linux": {
		"386",
		"amd64",
		"amd64p32",
		"arm",
		"arm64",
		"arm64be",
		"armbe",
		"loong64",
		"mips",
		"mips64",
		"mips64le",
		"mipsle",
		"ppc",
		"ppc64",
		"ppc64le",
		"riscv",
		"riscv64",
		"s390",
		"s390x",
		"sparc",
		"sparc64",
	},
	"netbsd":  {"386", "amd64", "arm", "arm64"},
	"openbsd": {"386", "amd64", "arm", "arm64", "mips64", "ppc64", "riscv64"},
	"plan9":   {"386", "amd64", "arm"},
	"solaris": {"amd64"},
	"wasip1":  {"wasm"},
	"windows": {"386", "amd64", "arm", "arm64"},
}

func IsValidOSAndArch(os, arch string) bool {
	arches, ok := validOSArch[strings.ToLower(os)]
	if !ok {
		return false
	}
	return slices.Contains(arches, strings.ToLower(arch))
}
