// Day 9: Disk Fragmenter
// https://adventofcode.com/2024/day/9

package main

import (
	"fmt"
	"os"
	"strconv"
)

type file struct {
	id, size int
}

type filesystem struct {
	files  []file
	blocks []*file
}

// Move a block from one location to another, leaving free space behind
func (fs *filesystem) moveBlock(srcIdx, destIdx int) {
	block := fs.blocks[srcIdx]

	fs.blocks[destIdx] = block
	fs.blocks[srcIdx] = nil
}

// Move a file from one location to another, leaving free space behind
func (fs *filesystem) moveFile(f file, destIdx int) {
	blocks := fs.fileLocations(f)

	// Iterate over the blocks the file currently resides at
	// Move it to the destination space
	// Source blocks become free space
	for i, srcIdx := range blocks {
		fs.blocks[destIdx+i] = &f
		fs.blocks[srcIdx] = nil
	}
}

// Returns all block indexes that make up a file
// Block indexes will be returned in ascending order
func (fs *filesystem) fileLocations(file file) []int {
	indexes := make([]int, 0)

	for i, b := range fs.blocks {
		if b != nil && *b == file {
			indexes = append(indexes, i)
		}
	}

	return indexes
}

func (fs *filesystem) String() string {
	str := ""

	for _, block := range fs.blocks {
		if block != nil {
			str += fmt.Sprintf("%d", block.id)
		} else {
			str += "."
		}
	}

	return str
}

func newFilesystem(diskmap []byte) *filesystem {
	// Diskmap is a series of block sizes
	// Even indexes are files, odd indexes are empty space
	// File IDs are in order, but are only incremented per file rather than per block
	files := make([]file, 0)
	blocks := make([]*file, 0)

	for i, b := range diskmap {
		blockSize, _ := strconv.Atoi(string(b))
		var f *file

		if i%2 == 0 {
			// Add file details to file listing
			fileId := i / 2
			f = &file{
				id:   fileId,
				size: blockSize,
			}

			files = append(files, *f)
		}

		// Write blocks
		// Each block points to the file details, or nil if free space
		for range blockSize {
			blocks = append(blocks, f)
		}
	}

	return &filesystem{files, blocks}
}

func main() {
	if cap(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [inputfile]\n", os.Args[0])
		os.Exit(-1)
	}

	content, err := os.ReadFile(os.Args[1])

	if err != nil {
		panic(err)
	}

	part1(content)
	part2(content)
}

// Part one is to compress a filesystem by moving blocks from the end of the filesystem to empty space at the beginning
// Part one does not care about fragmenting files
func part1(input []byte) {
	fs := newFilesystem(input)

	compress(fs)

	fmt.Println("Part 1 Checksum:", checksum(fs))
}

// Part two compresses the filesystem but also respects file fragmentation
// Files are kept together and must be moved to a span of empty space big enough to house the whole file
func part2(input []byte) {
	fs := newFilesystem(input)

	fileCompress(fs)

	fmt.Println("Part 2 Checksum:", checksum(fs))
}

// Move blocks from the end of the filesystem to fill in empty space at the start of the filesystem
// Used for Part 1
func compress(fs *filesystem) {
	tailIdx := len(fs.blocks) - 1

	for curIdx := 0; curIdx < tailIdx; curIdx++ {
		curBlock := fs.blocks[curIdx]

		// Is this block free space?
		if curBlock == nil {
			// Find a block at the end of the filesystem to move
			for tailIdx > curIdx {
				tailBlock := fs.blocks[tailIdx]

				// If a file is found, move it
				if tailBlock != nil {
					fs.moveBlock(tailIdx, curIdx)
					break
				}

				tailIdx--
			}
		}
	}
}

// Move full files from the end of the filesystem to available space at the start of the filesystem, ensuring that files are not fragmented
// Used for Part 2
func fileCompress(fs *filesystem) {
	// Starting at the last file index, descend the file list trying to move the file to free space at the start of the filesystem
	for fileIdx := len(fs.files) - 1; fileIdx >= 0; fileIdx-- {
		file := fs.files[fileIdx]

		fileStartIdx := fs.fileLocations(file)[0]

		// Look for free space starting from the beginning of the filesystem that can house the file
		if freeBlockIdx, ok := seekfree(fs, file.size, fileStartIdx); ok {
			fs.moveFile(file, freeBlockIdx)
		}
	}
}

// Looks for empty space starting at the beginning of the fs to contain the number of blocks given
func seekfree(fs *filesystem, size int, beforeIdx int) (int, bool) {
	for curIdx := 0; curIdx < beforeIdx; curIdx++ {
		if fs.blocks[curIdx] != nil {
			continue
		}

		freeStart := curIdx
		freeSize := 0

		// Found a block of free space
		// Loop over the next blocks to find where free space ends
		for fs.blocks[curIdx] == nil {
			freeSize++
			curIdx++
		}

		if freeSize >= size {
			return freeStart, true
		}
	}

	return 0, false
}

// Checksum multiplies each block's index with its file ID and adds them all up
func checksum(fs *filesystem) int {
	checksum := 0

	for i, block := range fs.blocks {
		if block != nil {
			checksum += i * block.id
		}
	}

	return checksum
}
