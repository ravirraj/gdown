package types

type FileInfo struct {
	Size          int64
	SupportRange  bool
	FileName      string
	RangeVerified bool
}

type Chunk struct {
	Index int
	Start int64
	End   int64
}
