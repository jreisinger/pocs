package main

func chunkSlice(slice []string, maxSize int) [][]string {
	if len(slice) <= 0 || maxSize < len(slice) {
		return nil
	}
	var chunks [][]string
	for maxSize < len(slice) {
		slice, chunks = slice[maxSize:], append(chunks, slice[0:maxSize:maxSize])
	}
	return append(chunks, slice)
}
