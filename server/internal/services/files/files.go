package files

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func DownloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http %s", resp.Status)
	}

	// Пытаемся взять имя из Content-Disposition или из URL
	filename := "downloaded_file"
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		// можно распарсить filename="..." (но для простоты опустим)
	}
	if filename == "downloaded_file" {
		filename = path.Base(url) // берём имя из URL
	}

	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return filename, err
}
