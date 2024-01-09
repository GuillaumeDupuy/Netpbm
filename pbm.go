package netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPBM(filename string) (*PBM, error){
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var pbm PBM

	// Read magic number
	if scanner.Scan() {
		pbm.magicNumber = scanner.Text()
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

		fmt.Sscanf(line, "%d %d", &pbm.width, &pbm.height)
		break
	}

	// Read data
	data := make([][]bool, pbm.height)
	for i := range data {
		data[i] = make([]bool, pbm.width)
	}

	switch pbm.magicNumber {
	case "P1":
		for i := 0; i < pbm.height; i++ {
			scanner.Scan()
			line := scanner.Text()
			for j, c := range line {
				if c == '1' {
					data[i][j] = true
				} else if c != '0' {
					return nil, fmt.Errorf("invalid PBM file: invalid character '%c'", c)
				}
			}
		}
	case "P4":
		for i := 0; i < pbm.height; i++ {
			scanner.Scan()
			line := scanner.Bytes()
			for j := 0; j < pbm.width; j++ {
				if line[j/8]&(1<<(7-uint(j%8))) != 0 {
					data[i][j] = true
				}
			}
		}
	default:
		return nil, fmt.Errorf("invalid PBM file: invalid magic number '%s'", pbm.magicNumber)
	}
	
	return &pbm, nil
}