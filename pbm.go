package netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PBM struct {
	Data          [][]bool
	Width, Height int
	MagicNumber   string
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var pbm PBM

	// Read magic number
	if scanner.Scan() {
		pbm.MagicNumber = scanner.Text()
	} else {
		return nil, fmt.Errorf("invalid PBM file: missing magic number")
	}

	// Read width and height
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		fmt.Sscanf(line, "%d %d", &pbm.Width, &pbm.Height)
		break
	}

	// Read data
	data := make([][]bool, pbm.Height)
	for i := range data {
		data[i] = make([]bool, pbm.Width)
	}

	switch pbm.MagicNumber {
	case "P1":
		for i := 0; i < pbm.Height; i++ {
			scanner.Scan()
			line := strings.Fields(scanner.Text()) // Split the line into fields
			for j, val := range line {
				if val == "1" {
					data[i][j] = true
				} else if val != "0" {
					return nil, fmt.Errorf("invalid PBM file: invalid character '%s'", val)
				}
			}
		}
	case "P4":
		for i := 0; i < pbm.Height; i++ {
			scanner.Scan()
			line := scanner.Bytes()
			for j := 0; j < pbm.Width; j++ {
				if line[j/8]&(1<<(7-uint(j%8))) != 0 {
					data[i][j] = true
				}
			}
		}
	default:
		return nil, fmt.Errorf("invalid PBM file: invalid magic number '%s'", pbm.MagicNumber)
	}

	return &pbm, nil
}
