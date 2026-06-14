package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: ./patch_source <target_nrx_file> <old_block_file> <new_block_file>")
		os.Exit(1)
	}

	nrxFile := os.Args[1]
	oldFile := os.Args[2]
	newFile := os.Args[3]

	nrxContent, err := os.ReadFile(nrxFile)
	if err != nil {
		fmt.Printf("Error reading target file %s: %v\n", nrxFile, err)
		os.Exit(1)
	}

	oldContent, err := os.ReadFile(oldFile)
	if err != nil {
		fmt.Printf("Error reading old block file %s: %v\n", oldFile, err)
		os.Exit(1)
	}

	newContent, err := os.ReadFile(newFile)
	if err != nil {
		fmt.Printf("Error reading new block file %s: %v\n", newFile, err)
		os.Exit(1)
	}

	// Normalize line endings and whitespace to prevent match failure
	normalize := func(b []byte) []byte {
		b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
		return bytes.TrimSpace(b)
	}

	normOld := normalize(oldContent)
	normNew := normalize(newContent)
	normNrx := bytes.ReplaceAll(nrxContent, []byte("\r\n"), []byte("\n"))

	// We search for normOld in normNrx
	if !bytes.Contains(normNrx, normOld) {
		// Try a fallback search with trimmed lines or direct containment checking if spacing differs
		fmt.Println("Error: Old block not found in target NetRexx file.")
		os.Exit(1)
	}

	patchedNrx := bytes.Replace(normNrx, normOld, normNew, 1)

	err = os.WriteFile(nrxFile, patchedNrx, 0644)
	if err != nil {
		fmt.Printf("Error writing patched file %s: %v\n", nrxFile, err)
		os.Exit(1)
	}

	fmt.Println("Success: Source file patched successfully.")
}
