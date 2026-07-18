package helper

import (
	"os"
	"os/user"
	"runtime"
)

func GetRuntimeUsername() string {
	u, err := user.Current()
	if err != nil {
		return "unknown"
	}
	return u.Username
}

func GetRuntimeHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func GetRuntimeOS() string {
	return runtime.GOOS
}

func GetRuntimeArch() string {
	return runtime.GOARCH
}
