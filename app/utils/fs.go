package utils

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

const caCertFile = "/etc/ssl/certs/ca-certificates.pem"
const caCertFolder = "/etc/ssl/certs"

func GetAbsolutePath(relativePath string) string {
	return filepath.Join("/mnt/onboard/.adds/kobotracker", relativePath)
}

func SetupCertificates() error {
	if _, err := os.Stat(caCertFile); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(caCertFolder, os.ModePerm)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		sourceFile, err := os.Open(GetAbsolutePath("assets/ca-certificates.pem"))
		if err != nil {
			log.Println(err.Error())
			return err
		}
		defer sourceFile.Close()

		newFile, err := os.Create(caCertFile)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, sourceFile)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	return nil
}
