package main

type Raid1 struct {
	phys *DiskArray
}

// Writes data to the RAID 1 array at the given logical block address (LBA).
func (r *Raid1) Write(lba int, data []byte) error {
	// RAID 1 mirrors: We must write the same data to ALL disks at the same offset
	offset := int64(lba * BlockSize) // Calculate offset within the disks
	for i := range NumDisks { // Write to all disks
		if err := r.phys.WritePhysical(i, offset, data); err != nil { // Write to all disks
			return err
		}
	}
	return nil // All writes successful
}

// Reads data from the RAID 1 array at the given logical block address (LBA).
func (r *Raid1) Read(lba int) ([]byte, error) {
	// Read from primary (Disk 0) for simplicity. No need to worry about failures here.
	offset := int64(lba * BlockSize) // Calculate offset within the disk
	return r.phys.ReadPhysical(0, offset) // Perform the read
}