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
		Timeout: 30 * time.Second,
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

	// send a get request to confirm the if the file supprots the parital downlaod

	resGet, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fileInfo, err
	}
	resGet.Header.Set("Range", "bytes=0-1023")

	respGet, err := client.Do(resGet)
	if err != nil {
		return fileInfo, err
	}
	defer respGet.Body.Close()

	if respGet.StatusCode != http.StatusPartialContent {
		return fileInfo, fmt.Errorf("server does not support partial content")

	}

	fileInfo.RangeVerified = true
	return fileInfo, nil
}
