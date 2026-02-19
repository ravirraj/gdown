# gdown

A concurrent file downloader written in Go. Downloads files in parallel chunks using HTTP Range requests and a worker pool.

## Features

- Parallel chunk-based downloads using goroutines
- HTTP Range request support for resumable downloads
- Configurable worker pool and chunk size
- Automatic file merging after download

## Build

```bash
go build -o gdown ./cmd
```

## Run

```bash
./gdown <url>
```

The downloader splits the file into 4 chunks by default and uses 4 workers to download them concurrently.

## Example

```bash
./gdown https://example.com/large-file.zip
```

The downloaded file will be saved in the current directory with its original name.

## Structure

```
gdown/
├── cmd/
│   └── main.go           # entry point, orchestrates download flow
├── internal/
│   ├── chunk/
│   │   └── split.go      # splits file into byte ranges
│   ├── httpclient/
│   │   ├── client.go     # checks URL and file metadata
│   │   └── download.go   # downloads individual chunks
│   ├── merger/
│   │   └── merger.go     # merges downloaded parts into final file
│   ├── types/
│   │   └── types.go      # shared data structures
│   └── worker/
│       └── worker.go     # manages worker pool and channels
└── download/             # temporary directory for chunk storage
```
