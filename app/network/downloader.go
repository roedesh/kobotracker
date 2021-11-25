package network

import (
	"io"
	"os"
)

func DownloadFile(filepath string, url string) error {
	client := GetHttpClient()

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
