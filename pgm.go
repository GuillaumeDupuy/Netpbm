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

// At returns the value of the pixel at (x, y).
func (pgm *PGM) Size() (int,int){
	return pgm.Width, pgm.Height
}

// At returns the value of the pixel at (x, y).
func (pgm *PGM) At(x, y int) uint8{
	return pgm.Data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8){
	pgm.Data[y][x] = value
}

// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error{
	file, err := os.Create(filename)
	if err != nil{
		return err
	}
	defer file.Close()

	// Write magic number
	fmt.Fprintf(file, "%s\n", pgm.MagicNumber)

	// Write width and height
	fmt.Fprintf(file, "%d %d\n", pgm.Width, pgm.Height)

	// Write max
	fmt.Fprintf(file, "%d\n", pgm.Max)

	// Write data
	switch pgm.MagicNumber {
	case "P2":
		for i := 0; i < pgm.Height; i++ {
			for j := 0; j < pgm.Width; j++ {
				fmt.Fprintf(file, "%d ", pgm.Data[i][j])
			}
			fmt.Fprintln(file)
		}
	case "P5":
		for i := 0; i < pgm.Height; i++ {
			for j := 0; j < pgm.Width; j++ {
				fmt.Fprintf(file, "%c", pgm.Data[i][j])
			}
		}
	default:
		return fmt.Errorf("invalid PGM file: invalid magic number")
	}

	return nil
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert(){
	for i := 0; i < pgm.Height; i++ {
		for j := 0; j < pgm.Width; j++ {
			pgm.Data[i][j] = uint8(pgm.Max) - pgm.Data[i][j]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip(){
	for i := 0; i < pgm.Height; i++ {
		for j := 0; j < pgm.Width/2; j++ {
			pgm.Data[i][j], pgm.Data[i][pgm.Width-j-1] = pgm.Data[i][pgm.Width-j-1], pgm.Data[i][j]
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop(){
	for i := 0; i < pgm.Height/2; i++ {
		for j := 0; j < pgm.Width; j++ {
			pgm.Data[i][j], pgm.Data[pgm.Height-i-1][j] = pgm.Data[pgm.Height-i-1][j], pgm.Data[i][j]
		}
	}
}

// SetMagicNumber sets the magic number of the PGM image and change the data format.
func (pgm *PGM) SetMagicNumber(magicNumber string){
	pgm.MagicNumber = magicNumber
	switch magicNumber {
	case "P2":
		for i := 0; i < pgm.Height; i++ {
			for j := 0; j < pgm.Width; j++ {
				pgm.Data[i][j] = uint8(pgm.Data[i][j])
			}
		}
	case "P5":
		for i := 0; i < pgm.Height; i++ {
			for j := 0; j < pgm.Width; j++ {
				pgm.Data[i][j] = uint8(pgm.Data[i][j])
			}
		}
	}
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8){
	pgm.Max = int(maxValue)
}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW(){
	pgm.Width, pgm.Height = pgm.Height, pgm.Width
	newData := make([][]uint8, pgm.Height)
	for i := range newData {
		newData[i] = make([]uint8, pgm.Width)
	}
	for i := 0; i < pgm.Height; i++ {
		for j := 0; j < pgm.Width; j++ {
			newData[i][j] = pgm.Data[pgm.Width-j-1][i]
		}
	}
	pgm.Data = newData
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM{
	var pbm PBM
	pbm.MagicNumber = "P1"
	pbm.Width = pgm.Width
	pbm.Height = pgm.Height
	pbm.Data = make([][]bool, pbm.Height)
	for i := range pbm.Data {
		pbm.Data[i] = make([]bool, pbm.Width)
	}
	for i := 0; i < pgm.Height; i++ {
		for j := 0; j < pgm.Width; j++ {
			if pgm.Data[i][j] > uint8(pgm.Max/2) {
				pbm.Data[i][j] = true
			} else {
				pbm.Data[i][j] = false
			}
		}
	}
	return &pbm
}