//go:build !container
// +build !container

package global

const (
	DefaultApiPort = 8083
	EnableLogFile  = true
	LoggerType     = "console"
)
