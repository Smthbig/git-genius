package system

import (
	"fmt"
	"os"
	"time"
)

const (
	geniusDir = ".git/.genius"
	errorLog  = geniusDir + "/error.log"
)

func LogError(context string, err error) {
	if err == nil {
		return
	}

	os.MkdirAll(geniusDir, 0700)

	f, ferr := os.OpenFile(errorLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if ferr != nil {
		return // last-resort: silently fail
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("[%s] %s: %v\n", timestamp, context, err)

	f.WriteString(line)
}
