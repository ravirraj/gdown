package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ravirraj/gdown/internal/chunk"
	"github.com/ravirraj/gdown/internal/httpclient"
	"github.com/ravirraj/gdown/internal/merger"
	_ "github.com/ravirraj/gdown/internal/merger"
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


	//get all file info 
	FileInfo, err := httpclient.CheckUrl(arg)
	if err != nil {
		slog.Error("ERROR GETTING FILE DETAILS ", "error", err)
		return
	}

	fmt.Println(FileInfo)


	//make partes of that file 
	chunks := chunk.SplitIntoChuncks(FileInfo.Size, 4)
	fmt.Println(chunks)

	// err = httpclient.DownloadChunnk(arg, chunks[0], "ravi")
	// if err != nil {
	// 	panic(err)
	// }


	//downlaod every part 
	baseUrl := filepath.Base(arg)
	err = worker.StartWorkers(arg, chunks, baseUrl, 8)
	if err != nil {
		panic(err)
	}

	//merge the parts to one file 
	err = merger.MergerFiles(baseUrl,4)
	if err != nil {
		panic(err)
	}

	

}
