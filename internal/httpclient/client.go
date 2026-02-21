package httpclient

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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

	res.Header.Set("Range", "bytes=0-1")
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(res)
	if err != nil {
		return fileInfo, err
	}
	fmt.Println(resp)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return fileInfo, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	fileSize := resp.Header.Get("Content-range")
	part := strings.Split(fileSize, "/")

	fileSizeAfterSlash := part[1]

	actualFileSize, err := strconv.ParseInt(fileSizeAfterSlash, 10, 64)
	if err != nil {
		return fileInfo, err
	}

	if actualFileSize < 0 {

		return fileInfo, fmt.Errorf("NO ContentLength")

	}

	// fmt.Println(resp.Header)
	if resp.Header.Get("Accept-Ranges") == "bytes" {
		fileInfo.SupportRange = true
	}

	fileInfo.Size = actualFileSize

	fileInfo.FileName = file

	// send a get request to confirm the if the file supprots the parital downlaod

	if resp.StatusCode != http.StatusPartialContent {
		return fileInfo, fmt.Errorf("server does not support partial content")

	}

	fileInfo.SupportRange = true
	return fileInfo, nil

}
