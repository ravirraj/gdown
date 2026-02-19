package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ravirraj/gdown/internal/chunk"
	"github.com/ravirraj/gdown/internal/httpclient"
	_ "github.com/ravirraj/gdown/internal/types"
	"github.com/ravirraj/gdown/internal/worker"
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

	chunks := chunk.SplitIntoChuncks(5242880, 4)
	fmt.Println(chunks)

	// err = httpclient.DownloadChunnk(arg, chunks[0], "ravi")
	// if err != nil {
	// 	panic(err)
	// }

	err = worker.StartWorkers(arg, chunks, "ravi", 8)
	if err != nil {
		panic(err)
	}
}
