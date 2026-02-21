package httpclient
import "os"
type ProgressWritter struct {

	file *os.File
	progress chan int64
} 