# EECE 5811 - Homework 7 - RAID Simulation in Go

## Overview of Test Results

### 100MB Disk Size, 4KB Block Size, 5 Disks
| RAID   | Write Time             | Write Speed | Read Time              | Read Speed  |
|--------|------------------------|-------------|------------------------|-------------|
| RAID 0 | 97.64s ( 3.8 ms/blk)   | 1.0 MB/s    | 0.19s ( 0.0 ms/blk)    | 515.8 MB/s  |
| RAID 1 | 531.36s (20.8 ms/blk)  | 0.2 MB/s    | 0.67s ( 0.0 ms/blk)    | 149.3 MB/s  |
| RAID 4 | 224.35s ( 8.8 ms/blk)  | 0.4 MB/s    | 0.15s ( 0.0 ms/blk)    | 660.4 MB/s  |
| RAID 5 | 198.20s ( 7.7 ms/blk)  | 0.5 MB/s    | 0.17s ( 0.0 ms/blk)    | 600.2 MB/s  |

### 10MB Disk Size, 4KB Block Size, 5 Disks

#### Run 1
| RAID   | Write Time             | Write Speed | Read Time              | Read Speed  |
|--------|------------------------|-------------|------------------------|-------------|
| RAID 0 | 10.30s ( 4.0 ms/blk)   | 1.0 MB/s    | 0.02s ( 0.0 ms/blk)    | 653.3 MB/s  |
| RAID 1 | 51.17s (20.0 ms/blk)   | 0.2 MB/s    | 0.02s ( 0.0 ms/blk)    | 401.9 MB/s  |
| RAID 4 | 20.99s ( 8.2 ms/blk)   | 0.5 MB/s    | 0.01s ( 0.0 ms/blk)    | 1476.9 MB/s |
| RAID 5 | 20.33s ( 7.9 ms/blk)   | 0.5 MB/s    | 0.01s ( 0.0 ms/blk)    | 968.0 MB/s  |

#### Run 2
| RAID   | Write Time             | Write Speed | Read Time              | Read Speed  |
|--------|------------------------|-------------|------------------------|-------------|
| RAID 0 | 10.29s ( 4.0 ms/blk)   | 1.0 MB/s    | 0.01s ( 0.0 ms/blk)    | 743.1 MB/s  |
| RAID 1 | 51.17s (20.0 ms/blk)   | 0.2 MB/s    | 0.02s ( 0.0 ms/blk)    | 462.5 MB/s  |
| RAID 4 | 20.34s ( 7.9 ms/blk)   | 0.5 MB/s    | 0.01s ( 0.0 ms/blk)    | 973.1 MB/s  |
| RAID 5 | 20.38s ( 8.0 ms/blk)   | 0.5 MB/s    | 0.01s ( 0.0 ms/blk)    | 902.7 MB/s  |

### 1 MB Disk Size, 4KB Block Size, 5 Disks

#### Run 1
| RAID   | Write Time             | Write Speed | Read Time              | Read Speed  |
|--------|------------------------|-------------|------------------------|-------------|
| RAID 0 | 1.03s ( 4.0 ms/blk)    | 1.0 MB/s    | 0.00s ( 0.0 ms/blk)    | 455.6 MB/s  |
| RAID 1 | 5.11s (19.9 ms/blk)    | 0.2 MB/s    | 0.00s ( 0.0 ms/blk)    | 534.1 MB/s  |
| RAID 4 | 2.04s ( 8.0 ms/blk)    | 0.5 MB/s    | 0.00s ( 0.0 ms/blk)    | 714.2 MB/s  |
| RAID 5 | 2.03s ( 7.9 ms/blk)    | 0.5 MB/s    | 0.00s ( 0.0 ms/blk)    | 1746.7 MB/s |

#### Run 2
| RAID   | Write Time             | Write Speed | Read Time              | Read Speed  |
|--------|------------------------|-------------|------------------------|-------------|
| RAID 0 | 1.03s ( 4.0 ms/blk)    | 1.0 MB/s    | 0.00s ( 0.0 ms/blk)    | 478.0 MB/s  |
| RAID 1 | 5.10s (19.9 ms/blk)    | 0.2 MB/s    | 0.00s ( 0.0 ms/blk)    | 489.3 MB/s  |
| RAID 4 | 2.03s ( 7.9 ms/blk)    | 0.5 MB/s    | 0.00s ( 0.0 ms/blk)    | 745.2 MB/s  |
| RAID 5 | 2.02s ( 7.9 ms/blk)    | 0.5 MB/s    | 0.00s ( 0.0 ms/blk)    | 1051.7 MB/s |


## Analysis

Comparing the performance numbers from my results to the theoretical expectations found in the textbook, my results are slightly different than what is expected. In my results, RAID 0 is the fastest for write operations, which aligns with the theoretical expectation found in the book. This is followed by RAID 5, RAID 4, and finally RAID 1 write speeds which also align with the textbook. However, looking closely at the actual write speeds, my results show lower speeds across all RAID levels compared to the theoretical values. The textbook states that RAID 0 should have T latency, RAID 1 should have T latency, RAID 4 should have 2T latency, and RAID 5 should have 2T latency for write operations. My results show that I have RAID 0 with T latency, RAID 1 with 5T latency, RAID 4 with 2T latency, and RAID 5 with 2T latency. The textbook is assuming parallel writes to disks whereas my implementation does sequential writes to disks. Nevertheless, the trend is still consistent with RAID 0 > RAID 5 > RAID 4 > RAID 1 for write speeds.

However, for read speeds my results are not consistent with the textbook's theoretical expectations. The textbook suggests that RAID 0 should have T latency, RAID 1 should have T latency, RAID 4 should have T latency, and RAID 5 should have T latency for read operations. Looking at the 100MB test run, my results show that RAID 0 has T latency, RAID 1 has approximately 4T latency, RAID 4 has T latency, and RAID 5 has approximately T latency. Likewise, the textbook suggests I should have a throughput of 1x for RAID 0, 0.5x for RAID 1, 0.8x for RAID 4, and 0.8x for RAID 5. My results show that I have a throughput of 1x for RAID 0, 0.3x for RAID 1, 1.3x for RAID 4, and 1.3x for RAID 5. The discrepancy in read speeds is most likely caused by how the OS handles page caching. 

#### Sources:
- Class lectures and materials.
- https://pages.cs.wisc.edu/~remzi/OSTEP/file-raid.pdf
- https://pkg.go.dev/builtin#make
    - Used to understand how to initialize slices in Go.
- https://pkg.go.dev/builtin#panic
    - Used to understand how to handle errors and panic in Go.
- https://pkg.go.dev/os
    - Used to understand file operations in Go.
- https://pkg.go.dev/io/ioutil
    - Used to understand reading and writing files in Go.
- https://pkg.go.dev/path/filepath
    - Used to understand file path manipulations in Go.
- https://www.ufsexplorer.com/articles/storage-technologies/raid-data-organization/
    - Reference for RAID levels and data organization.
- https://spdk.io/news/2024/02/12/raid5f/
    - Reference for RAID 5 implementation details.


## Overview of Code

### How to Run

1. Have Go installed on your machine. You can download it from https://golang.org/dl/.
2. Clone the repository to your local machine.
3. Navigate to the project directory in your terminal.
4. Run the following command to execute the program: `go run .`

### Design of Program
- main.go: Contains the main function to run the RAID simulation.
- common.go: Contains common constants and utility functions used across the program. You can change the disk size and block size here.
- utils.go: Contains utility functions for file operations and data generation. We create the disk files here in the directory './disks' which can be changed in common.go. The directory and files created and then deleted after each run.
- RAID\<level>.go: Contains the implementation of each RAID level (RAID 0, RAID 1, RAID 4, RAID 5).

### Libraries Used
- bytes: For byte slice manipulations and reading/writing data.
- fmt: For formatted I/O operations.
- time: For measuring performance.
- io: For I/O operations.
- path/filepath: For file path manipulations.