package imgkit

import (
	"io"
	"net/http"
	"os"
)

func Get(url, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func Remove(filename string) error {
	return os.Remove(filename)
}
