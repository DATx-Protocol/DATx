package util

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentPath() string {
	file, _ := os.Getwd()

	path, _ := filepath.Abs(file)

	log.Printf("get current path:%v", path)

	if len(path) == 0 {
		return path
	}

	ins := strings.Split(path, string(os.PathSeparator))
	return strings.Join(ins[:], string(os.PathSeparator))
}
