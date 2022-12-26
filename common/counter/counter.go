package counter

import (
	"fmt"
	"sync/atomic"
	"github.com/inhies/go-bytesize"
)
var _totalReadBytes uint64 = 0
var _totalWrittenBytes uint64 = 0
func IncrReadBytes(n int) {
	atomic.AddUint64(&_totalReadBytes, uint64(n))
}
func IncrWrittenBytes(n int) {
	atomic.AddUint64(&_totalWrittenBytes, uint64(n))
}
func GetReadBytes() uint64 {
	return _totalReadBytes
}
func GetWrittenBytes() uint64 {
	return _totalWrittenBytes
}
func PrintBytes(serverMode bool) string {
	if serverMode {
		return fmt.Sprintf("download %v upload %v", bytesize.New(float64(GetWrittenBytes())).String(), bytesize.New(float64(GetReadBytes())).String())
	}
	return fmt.Sprintf("download %v upload %v", bytesize.New(float64(GetReadBytes())).String(), bytesize.New(float64(GetWrittenBytes())).String())
}
