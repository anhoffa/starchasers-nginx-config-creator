package main

import (
	"io"
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

func saveInPersistentVolume(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			log.Errorw("Error closing file", "error", err)
		}
	}(sourceFile)

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {
			log.Errorw("Error closing file", "error", err)
		}
	}(destFile)

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
