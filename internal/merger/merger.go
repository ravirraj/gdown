package merger

import (
	"fmt"
	"io"
	"os"
)

func MergerFiles(baseUrl string, totalChunks int) error {
	file, err := os.Create(baseUrl)
	_ = file
	if err != nil {
		return err
	}

	defer file.Close()

	curretDir, err := os.Getwd()

	if err != nil {
		return err
	}

	for i := 0; i < totalChunks; i++ {
		filePath := fmt.Sprintf("%v/download/%v.part%v", curretDir, baseUrl, i)

		src, err := os.Open(filePath)
		if err != nil {
			return err
		}

		_, err = io.Copy(file, src)
		src.Close()

		if err != nil {
			return err
		}

	}

	//all this to just delete some file lol maybe this is not an optimal approch but it is what it is

	fileDele := fmt.Sprintf("%v/download", curretDir)

	op, err := os.Open(fileDele)
	if err != nil {
		return err
	}

	defer op.Close()

	entries, err := op.Readdir(-1)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		err := os.Remove(fmt.Sprintf("%v/%v", fileDele, entry.Name()))
		if err != nil {
			return err
		}
	}

	return nil

}
