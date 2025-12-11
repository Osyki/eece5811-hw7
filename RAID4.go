package main

type Raid4 struct {
	phys *DiskArray
}

// Writes data to the RAID 4 array at the given logical block address (LBA).
func (r *Raid4) Write(lba int, newData []byte) error {
	// RAID 4 uses dedicated parity disk (last disk)
	dataDisks := NumDisks - 1 // 4 data disks since 1 is for parity
	parityDiskIdx := 4        // Last disk is always parity in RAID 4

	// Map LBA to physical data disk
	stripeIdx := lba / dataDisks // Which stripe
	diskIdx := lba % dataDisks // Which data disk within the stripe
	offset := int64(stripeIdx * BlockSize) // Calculate offset within the disk

	// READ-MODIFY-WRITE Sequence for Parity Update
	// Read Old Data
	oldData, _ := r.phys.ReadPhysical(diskIdx, offset)
	// Read Old Parity
	oldParity, _ := r.phys.ReadPhysical(parityDiskIdx, offset)

	// Modify Parity: NewParity = OldParity ^ OldData ^ NewData
	tempXor := xorBlocks(oldParity, oldData)
	newParity := xorBlocks(tempXor, newData)

	// Write New Data
	if err := r.phys.WritePhysical(diskIdx, offset, newData); err != nil {
		return err
	}
	// Write New Parity
	return r.phys.WritePhysical(parityDiskIdx, offset, newParity)
}

// Reads data from the RAID 4 array at the given logical block address (LBA).
func (r *Raid4) Read(lba int) ([]byte, error) {
	dataDisks := NumDisks - 1 // 4 data disks since 1 is for parity
	diskIdx := lba % dataDisks // Map LBA to physical data disk
	offset := int64((lba / dataDisks) * BlockSize) // Calculate offset within the disk
	return r.phys.ReadPhysical(diskIdx, offset) // Perform the read
}
