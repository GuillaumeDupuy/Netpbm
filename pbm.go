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

	pbm.Data = data

	switch pbm.MagicNumber {
	case "P1":
		for i := 0; i < pbm.Height; i++ {
			scanner.Scan()
			line := strings.Fields(scanner.Text())
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

// Size returns the width and height of the image.
func (pbm *PBM) Size() (int,int){
	return pbm.Width, pbm.Height
}

// At returns the value of the pixel at (x, y).
func (pbm *PBM) At(x, y int) bool{
	return pbm.Data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool){
	pbm.Data[y][x] = value
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error{
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n", pbm.MagicNumber)
	fmt.Fprintf(file, "%d %d\n", pbm.Width, pbm.Height)

	switch pbm.MagicNumber {
	case "P1":
		for i := 0; i < pbm.Height; i++ {
			for j := 0; j < pbm.Width; j++ {
				if pbm.Data[i][j] {
					fmt.Fprintf(file, "1 ")
				} else {
					fmt.Fprintf(file, "0 ")
				}
			}
			fmt.Fprintln(file)
		}
	case "P4":
		for i := 0; i < pbm.Height; i++ {
			for j := 0; j < pbm.Width; j++ {
				if pbm.Data[i][j] {
					fmt.Fprintf(file, "1")
				} else {
					fmt.Fprintf(file, "0")
				}
			}
		}
	}

	return nil
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert(){
	for i := 0; i < pbm.Height; i++ {
		for j := 0; j < pbm.Width; j++ {
			pbm.Data[i][j] = !pbm.Data[i][j]
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip(){
	for i := 0; i < pbm.Height; i++ {
		for j := 0; j < pbm.Width/2; j++ {
			pbm.Data[i][j], pbm.Data[i][pbm.Width-1-j] = pbm.Data[i][pbm.Width-1-j], pbm.Data[i][j]
		}
	}
}

// Flop flops the PBM image vertically.
func (pbm *PBM) Flop(){
	for i := 0; i < pbm.Height/2; i++ {
		for j := 0; j < pbm.Width; j++ {
			pbm.Data[i][j], pbm.Data[pbm.Height-1-i][j] = pbm.Data[pbm.Height-1-i][j], pbm.Data[i][j]
		}
	}
}

// SetMagicNumber sets the magic number of the PBM image and change the data format.
func (pbm *PBM) SetMagicNumber(magicNumber string){
	pbm.MagicNumber = magicNumber
	switch magicNumber {
	case "P1":
		for i := 0; i < pbm.Height; i++ {
			for j := 0; j < pbm.Width; j++ {
				if pbm.Data[i][j] {
					pbm.Data[i][j] = true
				} else {
					pbm.Data[i][j] = false
				}
			}
		}
	case "P4":
		for i := 0; i < pbm.Height; i++ {
			for j := 0; j < pbm.Width; j++ {
				if pbm.Data[i][j] {
					pbm.Data[i][j] = true
				} else {
					pbm.Data[i][j] = false
				}
			}
		}
	}
}