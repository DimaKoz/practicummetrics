package common

import (
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

var workDir = ""

func GetWD() string {
	if workDir != "" {
		return workDir
	}
	wDirTemp, _ := os.Getwd()

	zap.S().Infof("wd started: %s", wDirTemp)

	for !strings.HasSuffix(wDirTemp, "practicummetrics") {
		wDirTemp = filepath.Dir(wDirTemp)
	}
	workDir = wDirTemp

	return workDir
}
