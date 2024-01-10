package netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PGM struct{
    Data [][]uint8
    Width, Height int
    MagicNumber string
    Max int
}

func ReadPGM(filename string) (*PGM, error){
	file, err := os.Open(filename)
	if err != nil{
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var pgm PGM

	// Read magic number
	if scanner.Scan(){
		pgm.MagicNumber = scanner.Text()
	} else{
		return nil, fmt.Errorf("invalid PGM file: missing magic number")
	}

	// Read width and height
	for scanner.Scan(){
		line := scanner.Text()

		// Skip comments
		if strings.HasPrefix(line, "#"){
			continue
		}

		fmt.Sscanf(line, "%d %d", &pgm.Width, &pgm.Height)
		break
	}

	// Read max
	for scanner.Scan(){
		line := scanner.Text()

		// Skip comments
		if strings.HasPrefix(line, "#"){
			continue
		}

		fmt.Sscanf(line, "%d", &pgm.Max)
		break
	}

	// Read data
	data := make([][]uint8, pgm.Height)
	for i := range data {
		data[i] = make([]uint8, pgm.Width)
	}

	pgm.Data = data

	switch pgm.MagicNumber {
	case "P2":
		for i := 0; i < pgm.Height; i++ {
			scanner.Scan()
			line := strings.Fields(scanner.Text())
			for j, val := range line {
				fmt.Sscanf(val, "%d", &pgm.Data[i][j])
			}
		}
	case "P5":
		for i := 0; i < pgm.Height; i++ {
			scanner.Scan()
			line := scanner.Bytes() // Read the line as bytes
			for j := 0; j < pgm.Width; j++ {
				pgm.Data[i][j] = uint8(line[j]) // Convert byte to uint8
			}
		}
	default:
		return nil, fmt.Errorf("invalid PGM file: invalid magic number")
	}

	return &pgm, nil
}