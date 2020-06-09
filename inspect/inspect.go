package inspect

import (
	_ "crypto/sha256"
	_ "fmt"
	"io/ioutil"
	"log"
	"os"
	_ "strings"
	"time"
)

const format = "2006-01-02 15:04:05.000000 MST"

var timezone = time.UTC

func ParseFile(filename string, blockDetails, printLines, storeBlocks bool) (*ChunkHeader, *LokiChunk, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("%s: %v", filename, err)
		return nil, nil, err
	}
	defer f.Close()

	_, err = f.Stat()
	if err != nil {
		log.Println("failed to stat file", err)
		return nil, nil, err
	}

	h, err := DecodeHeader(f)
	if err != nil {
		log.Printf("%s: %v", filename, err)
		return nil, nil, err
	}

	/*fmt.Println()
	fmt.Println("Chunks file:", filename)
	fmt.Println("Metadata length:", h.MetadataLength)
	fmt.Println("Data length:", h.DataLength)
	fmt.Println("UserID:", h.UserID)
	from, through := h.From.Time().In(timezone), h.Through.Time().In(timezone)
	fmt.Println("From:", from.Format(format))
	fmt.Println("Through:", through.Format(format), "("+through.Sub(from).String()+")")
	fmt.Println("Labels:")

	for _, l := range h.Metric {
		fmt.Println("\t", l.Name, "=", l.Value)
	}*/

	lokiChunk, err := parseLokiChunk(h, f)
	if err != nil {
		log.Printf("%s: %v", filename, err)
		return nil, nil, err
	}
	/*
		fmt.Println("Encoding:", lokiChunk.encoding)
		fmt.Print("Blocks Metadata Checksum: ", fmt.Sprintf("%08x", lokiChunk.metadataChecksum))
		if lokiChunk.metadataChecksum == lokiChunk.computedMetadataChecksum {
			fmt.Println(" OK")
		} else {
			fmt.Println(" BAD, computed checksum:", fmt.Sprintf("%08x", lokiChunk.computedMetadataChecksum))
		}
		if blockDetails {
			fmt.Println("Found", len(lokiChunk.Blocks), "block(s)")
		} else {
			fmt.Println("Found", len(lokiChunk.Blocks), "block(s), use -b to show block details")
		}
		if len(lokiChunk.Blocks) > 0 {
			fmt.Println("Minimum time (from first block):", time.Unix(0, lokiChunk.Blocks[0].minT).In(timezone).Format(format))
			fmt.Println("Maximum time (from last block):", time.Unix(0, lokiChunk.Blocks[len(lokiChunk.Blocks)-1].maxT).In(timezone).Format(format))
		}

		if blockDetails {
			fmt.Println()
		}

		totalSize := 0

		for ix, b := range lokiChunk.Blocks {
			if blockDetails {
				cksum := ""
				if b.storedChecksum == b.computedChecksum {
					cksum = fmt.Sprintf("%08x OK", b.storedChecksum)
				} else {
					cksum = fmt.Sprintf("%08x BAD (computed: %08x)", b.storedChecksum, b.computedChecksum)
				}
				fmt.Printf("Block %4d: position: %8d, original length: %6d (stored: %6d, ratio: %.2f), minT: %v maxT: %v, checksum: %s\n",
					ix, b.dataOffset, len(b.originalData), len(b.rawData), float64(len(b.originalData))/float64(len(b.rawData)),
					time.Unix(0, b.minT).In(timezone).Format(format), time.Unix(0, b.maxT).In(timezone).Format(format),
					cksum)
				fmt.Printf("Block %4d: digest compressed: %02x, original: %02x\n", ix, sha256.Sum256(b.rawData), sha256.Sum256(b.originalData))
			}

			totalSize += len(b.originalData)

			if printLines {
				for _, l := range b.Entries {
					fmt.Printf("%v\t%s\n", time.Unix(0, l.timestamp).In(timezone).Format(format), strings.TrimSpace(l.Line))
				}
			}

			if storeBlocks {
				writeBlockToFile(b.rawData, ix, fmt.Sprintf("%s.block.%d", filename, ix))
				writeBlockToFile(b.originalData, ix, fmt.Sprintf("%s.original.%d", filename, ix))
			}
		}

		fmt.Println("Total size of original data:", totalSize, "file size:", si.Size(), "ratio:", fmt.Sprintf("%0.3g", float64(totalSize)/float64(si.Size())))*/

	return h, lokiChunk, nil
}

func writeBlockToFile(data []byte, blockIndex int, filename string) {
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Println("Failed to store block", blockIndex, "to file", filename, "due to error:", err)
	} else {
		log.Println("Stored block", blockIndex, "to file", filename)
	}
}
