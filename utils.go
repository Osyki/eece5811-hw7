package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// DiskArray represents the collection of disk files used in the RAID setup.
type DiskArray struct {
	disks []*os.File
}

// NewDiskArray initializes the file handles and pre-allocates disk space.
// We allocate the full size upfront to stabilize performance during benchmarks.
func NewDiskArray() (*DiskArray, error) {
	// Initialize DiskArray
	da := &DiskArray{
		disks: make([]*os.File, NumDisks),
	}

	// Ensure the storage directory exists
	if err := os.MkdirAll(DiskDir, 0755); err != nil {
		return nil, err
	}

	for i := range da.disks {
		// Create file path: ./disks/disk{i}.dat
		filename := filepath.Join(DiskDir, fmt.Sprintf("disk%d.dat", i))
		
		// Open or create the disk file
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, err // Return error if file cannot be opened/created
		}

		// We set the file size to the full BenchmarkData size immediately.
		// This stabilizes performance by preventing metadata updates during the benchmark.
		if err := f.Truncate(BenchmarkData); err != nil {
			f.Close()
			return nil, err
		}

		da.disks[i] = f
	}
	return da, nil
}

// Close ensures all file handles are closed properly.
func (d *DiskArray) Close() {
	for _, f := range d.disks {
		f.Close() // Close each disk file handle
	}
}

// CleanupFiles removes existing disk directory to ensure a clean benchmark state.
func CleanupFiles() {
	os.RemoveAll(DiskDir)
}

// ReadPhysical reads a physical block from a specific disk index.
func (d *DiskArray) ReadPhysical(diskIndex int, offset int64) ([]byte, error) {
	buffer := make([]byte, BlockSize) // Pre-allocate buffer
	_, err := d.disks[diskIndex].ReadAt(buffer, offset) // Read into buffer
	// Check for errors except EOF
	if err != nil && err != io.EOF { // Check for errors except EOF
		return nil, err // Return error if read fails
	}
	return buffer, nil // If EOF, return zero-filled buffer (valid for XOR logic)
}

// WritePhysical writes to a specific disk and performs fsync.
func (d *DiskArray) WritePhysical(diskIndex int, offset int64, data []byte) error {
	// Write data to the specified disk at the given offset
	if _, err := d.disks[diskIndex].WriteAt(data, offset); err != nil {
		return err // Return error if write fails
	}
	// Simulate real hardware write latency/durability
	return d.disks[diskIndex].Sync()
}

// XOR two blocks of data. Used in RAID 4/5 for parity calculations.
func xorBlocks(a, b []byte) []byte {
	result := make([]byte, len(a)) // Assume a and b are of equal length
	for i := range a { // XOR byte by byte
		result[i] = a[i] ^ b[i]
	}
	return result
}