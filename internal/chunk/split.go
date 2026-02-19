package chunk

import (
	"github.com/ravirraj/gdown/internal/types"
)

func SplitIntoChuncks(fileSize int64, workers int) []types.Chunk {

	var chuncks []types.Chunk
	chunckSize := fileSize / int64(workers)
	var start int64 = 0

	// end := chunckSize - 1

	// for i := 0; i < workers; i++ {

	// 	chuncks = append(chuncks, types.Chunk{Index: i, Start: int64(i*start), End: end})

	// }

	// for i := range workers - 1 {
	// 	chuncks = append(chuncks, types.Chunk{Index: i, Start: int64(start), End: end})
	// 	start += int(chunckSize)
	// 	end = end + chunckSize
	// }

	for i := 0; i < workers; i++ {
		end := start + chunckSize - 1

		if i == workers-1 {
			end = fileSize - 1
		}

		chuncks = append(chuncks, types.Chunk{
			Index: i,
			Start: start,
			End:   end,
		})

		start += chunckSize
	}

	// chuncks = append(chuncks, types.Chunk{Index: 4, Start: int64(start), End: fileSize - 1})
	return chuncks
}
