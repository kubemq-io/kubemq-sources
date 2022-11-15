//go:build !linux
// +build !linux

package logger

import (
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}
