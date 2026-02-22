package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/ravirraj/gdown/internal/chunk"
	"github.com/ravirraj/gdown/internal/httpclient"
	"github.com/ravirraj/gdown/internal/merger"
	"github.com/ravirraj/gdown/internal/worker"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	progress := make(chan int64, 100)

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

	totalFileSize := FileInfo.Size
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	go func() {
		var downloaded int64
		for {
			select {
			case n, ok := <-progress:
				if !ok {
					return
				}
				downloaded += n

			case <-ticker.C:
				percent := float64(downloaded) / float64(totalFileSize) * 100
				fmt.Printf("\rDownloading: %.2f%%", percent)

			}
			// fmt.Printf("Downlaoded %v",downlaoded)
			// percentage := downloaded / totalFileSize * 100
		}
	}()

	baseUrl := filepath.Base(arg)
	err = worker.StartWorkers(ctx ,arg, chunks, baseUrl, 4, progress)
	if err != nil {
		panic(err)
	}

	//merge the parts to one file
	err = merger.MergerFiles(baseUrl, 4)
	if err != nil {
		panic(err)
	}

}
