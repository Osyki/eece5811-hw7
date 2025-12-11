package main

type Raid0 struct {
	phys *DiskArray
}

// Writes data to the RAID 0 array at the given logical block address (LBA).
func (r *Raid0) Write(lba int, data []byte) error {
	// RAID 0 stripes data across all disks
	diskIdx := lba % NumDisks // Determine which disk to write to
	offset := int64((lba / NumDisks) * BlockSize) // Calculate offset within the disk
	return r.phys.WritePhysical(diskIdx, offset, data) // Perform the write
}

// Reads data from the RAID 0 array at the given logical block address (LBA).
func (r *Raid0) Read(lba int) ([]byte, error) {
	diskIdx := lba % NumDisks // Determine which disk to read from
	offset := int64((lba / NumDisks) * BlockSize) // Calculate offset within the disk
	return r.phys.ReadPhysical(diskIdx, offset) // Perform the read
}
