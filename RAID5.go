package main

type Raid5 struct {
	phys *DiskArray
}

// GetMapping maps a logical block address (LBA) to physical disk index and offset
func (r *Raid5) GetMapping(lba int) (dataDiskIdx int, parityDiskIdx int, offset int64) {
	// RAID 5 uses distributed parity in a rotational manner
	dataDisks := NumDisks - 1 // We subtract 1 here since one disk is used for parity
	stripeRow := lba / dataDisks // Which stripe row
	logicalCol := lba % dataDisks // Which data chunk in the stripe (0..3)

	// Calculate Parity Disk
	// We use Left-Symmetric Rotation for Parity Distribution
	// Row 0: Parity on Disk 4. Row 1: Parity on Disk 3, etc.
	parityDiskIdx = (NumDisks - 1) - (stripeRow % NumDisks)

	// Calculate Physical Data Disk
	// We skip the parity disk when mapping the logical column
	if logicalCol < parityDiskIdx {
		dataDiskIdx = logicalCol
	} else {
		dataDiskIdx = logicalCol + 1
	}
	// Calculate Offset within the disk
	offset = int64(stripeRow * BlockSize)
	return
}

// Writes data to the RAID 5 array at the given logical block address (LBA).
func (r *Raid5) Write(lba int, newData []byte) error {
	diskIdx, parityIdx, offset := r.GetMapping(lba) // Get physical disk and offset

	// READ-MODIFY-WRITE Sequence
	// Read Old Data
	oldData, _ := r.phys.ReadPhysical(diskIdx, offset)
	// Read Old Parity
	oldParity, _ := r.phys.ReadPhysical(parityIdx, offset)

	// Modify Parity: NewParity = OldParity ^ OldData ^ NewData
	tempXor := xorBlocks(oldParity, oldData)
	newParity := xorBlocks(tempXor, newData)

	// Write New Data
	if err := r.phys.WritePhysical(diskIdx, offset, newData); err != nil {
		return err
	}
	// Write New Parity
	return r.phys.WritePhysical(parityIdx, offset, newParity)
}

// Reads data from the RAID 5 array at the given logical block address (LBA).
func (r *Raid5) Read(lba int) ([]byte, error) {
	diskIdx, _, offset := r.GetMapping(lba) // Get physical disk and offset
	return r.phys.ReadPhysical(diskIdx, offset) // Perform the read
}