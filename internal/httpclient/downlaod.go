package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ravirraj/gdown/internal/types"
)

func (pw *ProgressWritter) Write(b []byte) (int, error) {
	n, err := pw.file.Write(b)
	if n > 0 {
		pw.progress <- int64(n)
	}

	return n, err
}

func DownloadChunnk(ctx context.Context, client *http.Client, url string, c types.Chunk, baseFileurl string, progressChan chan int64) error {

	// this will downlaod the perticular part of the file and if the downloading fails it will retry it , if it still fails it will stop

	//this to thorw last error after retrying
	var lastErr error

	//this to get file info and after that we will create a dir called downlaod and in that , all parts are saved
	fileThe, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return err

	}

	dowlaodDir := filepath.Join(fileThe, "download")
	err = os.MkdirAll(dowlaodDir, 0755)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("%v.part%v", baseFileurl, c.Index)

	filePath := filepath.Join(dowlaodDir, fileName)

	//retry logic
	for i := 0; i < 3; i++ {

		err := downlaod(ctx, client, url, c, filePath, progressChan)

		if err == nil {
			fmt.Println("NO ERRORS ")
			return nil
		}

		//removing the partially downlaoded file
		os.Remove(filePath)

		lastErr = err

		// 2 second sleep after the files fails
		if i < 2 {
			time.Sleep(2 * time.Second)
		}

	}

	//progress bar implemantation

	return fmt.Errorf("failed after retries %w", lastErr)
}

// this functions downlaod the actully file , main download logic is in here
func downlaod(ctx context.Context, client *http.Client, url string, c types.Chunk, filePath string, progressChan chan int64) error {
	resGet, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resGet.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", c.Start, c.End))

	respGet, err := client.Do(resGet)
	if err != nil {
		return err
	}

	defer respGet.Body.Close()

	if respGet.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("file does not support the partial downaload ")
	}

	expectedSize := (c.End - c.Start) + 1

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	//custom writter for progress

	progressWriter := &ProgressWritter{file: file, progress: progressChan}

	downlaodedFile, err := io.Copy(progressWriter, respGet.Body)
	if err != nil {
		return err
	}

	if downlaodedFile != expectedSize {
		return fmt.Errorf("Downlaod incomplete")
	}
	return nil
}
