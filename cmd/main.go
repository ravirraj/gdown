package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ravirraj/gdown/internal/chunk"
	"github.com/ravirraj/gdown/internal/httpclient"
	_ "github.com/ravirraj/gdown/internal/types"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("reqired a field ")
		return
	}

	arg := os.Args[1]

	fmt.Println(arg)

	FileInfo, err := httpclient.CheckUrl(arg)
	if err != nil {
		slog.Error("ERROR GETTING FILE DETAILS ", "error", err)
		return
	}

	fmt.Println(FileInfo)

	// resp, _ := http.NewRequest("GET", arg, nil)

	// resp.Header.Set("Range", "bytes=0-1023")
	// client := &http.Client{}

	// respg, _ := client.Do(resp)

	// defer respg.Body.Close()

	// fmt.Println(respg)

	chunks := chunk.SplitIntoChuncks(5242880, 4)
	fmt.Println(chunks)

	fmt.Println(chunks)
}
