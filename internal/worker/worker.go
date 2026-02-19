package worker

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/ravirraj/gdown/internal/httpclient"
	"github.com/ravirraj/gdown/internal/types"
)

func StartWorkers(url string, c []types.Chunk, baseUrl string, workers int) error {

	client := &http.Client{
		// Timeout: 30 * time.Second,
	}

	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()
	//create a wg(wait group)
	var wg sync.WaitGroup

	// create a channel
	jobs := make(chan types.Chunk)
	errChan := make(chan error, 1)

	//loop over the workres
	for i := 0; i < workers; i++ {
		wg.Add(1)

		//start go routines
		go func() {
			defer wg.Done()

			// fmt.Println("worker stared")
			for job := range jobs {
				select {
				case <-ctx.Done():
					fmt.Println("worker stoped ")
					return

				default:
					err := httpclient.DownloadChunnk(client, url, job,baseUrl)
					if err != nil {
						select {
						case errChan <- err:
						default:
						}
						cancle()
						return
					}
					fmt.Println("working on", job.Index)
				}

			}
		}()
	}

	//give job to workers
sendingJob:
	for _, job := range c {

		select {
		case jobs <- job:
		case <-ctx.Done():
			break sendingJob
		}
	}

	//close the job
	close(jobs)

	//wait for workres to complete the job
	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
	}
	fmt.Println("all work done")
	return nil
}
