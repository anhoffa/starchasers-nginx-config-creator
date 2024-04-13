package main

import (
	"os"
)

func checkForSavedConfig(src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {
			log.Errorw("Error closing file", "error", err)
		}
	}(srcFile)
	return nil
}
