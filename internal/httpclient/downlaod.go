package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ravirraj/gdown/internal/types"
)

func DownloadChunnk(client *http.Client, url string, c types.Chunk, baseFileurl string) error {

	resGet, err := http.NewRequest("GET", url, nil)
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

	//we have multiple files , so we need to loop over every file

	// for i := 0; i < c.Index; i++ {
	// 	fileName := fmt.Sprintf("baseFileurl.part{%v}", c.Index)
	// 	file, err := os.Create(fileName)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	defer file.Close()

	// 	_, err = io.Copy(file, respGet.Body)

	// }

	fileThe, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return err

	}

	dowlaodDir := filepath.Join(fileThe, "downlaod")

	err = os.MkdirAll(dowlaodDir, 0755)
	if err != nil {
		return err
	}

	expectedSize := (c.End - c.Start) + 1
	fileName := fmt.Sprintf("%v.part%v", baseFileurl, c.Index)

	filePath := filepath.Join(dowlaodDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	downlaodedFile, err := io.Copy(file, respGet.Body)
	if err != nil {
		return err
	}

	if downlaodedFile != expectedSize {
		return fmt.Errorf("Downlaod incomplete")
	}

	return err

}
