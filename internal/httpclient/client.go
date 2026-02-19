package httpclient

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ravirraj/gdown/internal/types"
)

func CheckUrl(url string) (types.FileInfo, error) {

	file := filepath.Base(url)

	var fileInfo types.FileInfo

	res, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return fileInfo, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(res)
	if err != nil {
		return fileInfo, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fileInfo, fmt.Errorf("unexpected status: %s", resp.Status)
	}


	fileSize := resp.ContentLength

	if fileSize < 0 {

		return fileInfo, fmt.Errorf("NO ContentLength")

	}

	if resp.Header.Get("Accept-Ranges") == "bytes" {
		fileInfo.SupportRange = true
	}

	fileInfo.Size = fileSize

	fileInfo.FileName = file

	return fileInfo, nil
}
