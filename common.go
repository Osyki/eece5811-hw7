package main

import "time"

const (
	BlockSize     = 4096             // 4 KB
	NumDisks      = 5                // disk0 to disk4
	BenchmarkData = 10 * 1024 * 1024 // 100 MB
	DiskDir       = "./disks"        // Directory to store disk files
)

// RAID interface defines the public contract for all implementations
type RAID interface {
	Write(blockNum int, data []byte) error
	Read(blockNum int) ([]byte, error)
}

// BenchResult holds the results of a benchmark run.
type BenchResult struct {
	Name      string
	WriteTime time.Duration
	ReadTime  time.Duration
	Bytes     int
	Blocks    int
}
