package infra

import (
	"encoding/json"
	"os"
)

type LogWriter struct{}

func (ref *LogWriter) detectLogLevel(p []byte) string {
	var logs map[string]string
	_ = json.Unmarshal(p, &logs)
	if _, ok := logs["level"]; ok {
		return logs["level"]
	}
	return ""
}

func (ref *LogWriter) Write(p []byte) (n int, err error) {
	level := ref.detectLogLevel(p)
	if level == "info" {
		return os.Stdout.Write(p)
	}
	return os.Stderr.Write(p)
}
