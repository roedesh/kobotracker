package config

import (
	"cryptokobo/app/utils"
	"errors"
	"io"
	"log"
	"os"
)

const caCertFile = "/etc/ssl/certs/ca-certificates.pem"
const caCertFolder = "/etc/ssl/certs"

func SetupLogger() func() {
	logFile, err := os.OpenFile(utils.GetAbsolutePath("debug.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFile)
	}

	return func() {
		if logFile != nil {
			logFile.Close()
		}
	}
}

func SetupSSLCertificates() error {
	if _, err := os.Stat(caCertFile); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(caCertFolder, os.ModePerm)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		sourceFile, err := os.Open(utils.GetAbsolutePath("assets/ca-certificates.pem"))
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
