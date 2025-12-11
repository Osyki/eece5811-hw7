package main

import (
	"bytes"
	"fmt"
	"time"
)

// RunBenchmark executes the standard workload against a specific RAID implementation.
// It cleans any existing files, sets up the DiskArray, and measures read/write performance.
// Then, it returns the benchmark results.
func RunBenchmark(name string, factory func(*DiskArray) RAID) BenchResult {
	// Clean up any existing files
	CleanupFiles()
	phys, err := NewDiskArray()
	if err != nil {
		panic(err)
	}
	defer phys.Close()
	
	raid := factory(phys) // Create RAID instance
	numBlocks := BenchmarkData / BlockSize // Total number of blocks to process
	data := make([]byte, BlockSize)
	// Fill dummy data
	copy(data, []byte("derpderpderpderp"))

	fmt.Printf("Benchmarking %-8s | %d Blocks | ", name, numBlocks) // Show progress
	// Measure Write
	startWrite := time.Now()
	for i := range numBlocks {
		if err := raid.Write(i, data); err != nil {
			panic(fmt.Sprintf("Write error at block %d: %v", i, err))
		}
	}
	writeDuration := time.Since(startWrite)

	// Measure Read
	startRead := time.Now()
	for i := range numBlocks {
		if _, err := raid.Read(i); err != nil {
			panic(fmt.Sprintf("Read error at block %d: %v", i, err))
		}
	}
	readDuration := time.Since(startRead)

	fmt.Println("Done.")

	return BenchResult{
		Name:      name,
		WriteTime: writeDuration,
		ReadTime:  readDuration,
		Bytes:     BenchmarkData,
		Blocks:    numBlocks,
	}
}

func main() {
	fmt.Printf("RAID Benchmark Tool\n")
	fmt.Printf("Disk Size: %d MB\n", BenchmarkData/(1024*1024))
	fmt.Printf("Block Size: %d bytes\n", BlockSize)

	// Array to collect results 
	results := []BenchResult{}

	// RAID 0
	results = append(results, RunBenchmark("RAID 0", func(p *DiskArray) RAID {
		return &Raid0{phys: p}
	}))

	// RAID 1
	results = append(results, RunBenchmark("RAID 1", func(p *DiskArray) RAID {
		return &Raid1{phys: p}
	}))

	// RAID 4
	results = append(results, RunBenchmark("RAID 4", func(p *DiskArray) RAID {
		return &Raid4{phys: p}
	}))

	// RAID 5
	results = append(results, RunBenchmark("RAID 5", func(p *DiskArray) RAID {
		return &Raid5{phys: p}
	}))

	// PRINT SUMMARY
	fmt.Printf("\n%s\n", bytes.Repeat([]byte("-"), 90))
	fmt.Printf("| %-8s | %-22s | %-11s | %-22s | %-11s |\n",
		"RAID", "Write Time", "Write Speed", "Read Time", "Read Speed")
	fmt.Printf("%s\n", bytes.Repeat([]byte("-"), 90))

	for _, res := range results {
		// Calculate stats
		writePerBlock := res.WriteTime.Seconds() * 1000 / float64(res.Blocks) // ms
		readPerBlock := res.ReadTime.Seconds() * 1000 / float64(res.Blocks)   // ms

		// Throughput MB/s
		writeMBps := float64(res.Bytes) / (1024 * 1024) / res.WriteTime.Seconds()
		readMBps := float64(res.Bytes) / (1024 * 1024) / res.ReadTime.Seconds()

		fmt.Printf("| %-8s | %7.2fs (%4.1f ms/blk) | %6.1f MB/s | %7.2fs (%4.1f ms/blk) | %6.1f MB/s |\n",
			res.Name,
			res.WriteTime.Seconds(),
			writePerBlock,
			writeMBps,
			res.ReadTime.Seconds(),
			readPerBlock,
			readMBps,
		)
	}
	fmt.Printf("%s\n", bytes.Repeat([]byte("-"), 90))

	// Cleanup at the end
	CleanupFiles()
}
