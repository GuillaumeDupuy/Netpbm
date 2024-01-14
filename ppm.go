package netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PPM struct{
    Data [][]Pixel
    Width, Height int
    MagicNumber string
    Max int
}

type Pixel struct{
    R, G, B uint8
}

type Point struct{
    X, Y int
}

// ReadPPM reads a PPM image from a file and returns a struct that represents the image.
func ReadPPM(filename string) (*PPM, error){
	file, err := os.Open(filename)
	if err != nil{
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var ppm PPM

	// Read magic number
	if scanner.Scan(){
		ppm.MagicNumber = scanner.Text()
	} else{
		return nil, fmt.Errorf("invalid PPM file: missing magic number")
	}

	// Read Width and Height
	for scanner.Scan(){
		line := scanner.Text()

		// Skip comments
		if strings.HasPrefix(line, "#"){
			continue
		}

		fmt.Sscanf(line, "%d %d", &ppm.Width, &ppm.Height)
		break
	}

	// Read Max
	for scanner.Scan(){
		line := scanner.Text()

		// Skip comments
		if strings.HasPrefix(line, "#"){
			continue
		}

		fmt.Sscanf(line, "%d", &ppm.Max)
		break
	}

	// Read Data
	Data := make([][]Pixel, ppm.Height)
	for i := range Data{
		Data[i] = make([]Pixel, ppm.Width)
	}

	ppm.Data = Data

	switch ppm.MagicNumber{
	case "P3":
		for i := 0; i < ppm.Height; i++{
			scanner.Scan()
			line := strings.Fields(scanner.Text())
			for j := 0; j < ppm.Width; j++{
				fmt.Sscanf(line[j], "%d %d %d", &ppm.Data[i][j].R, &ppm.Data[i][j].G, &ppm.Data[i][j].B)
			}
		}
	case "P6":
		for i := 0; i < ppm.Height; i++{
			scanner.Scan()
			line := scanner.Text()
			for j := 0; j < ppm.Width; j++{
				ppm.Data[i][j].R = line[j*3]
				ppm.Data[i][j].G = line[j*3+1]
				ppm.Data[i][j].B = line[j*3+2]
			}
		}
	default:
		return nil, fmt.Errorf("invalid PPM file: unknown magic number")
	}

	return &ppm, nil
}

// Size returns the Width and Height of the image.
func (ppm *PPM) Size() (int,int){
	return ppm.Width, ppm.Height
}

// At returns the value of the pixel at (x, y).
func (ppm *PPM) At(x, y int) Pixel{
	return ppm.Data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (ppm *PPM) Set(x, y int, value Pixel){
	ppm.Data[y][x] = value
}

// Save saves the PPM image to a file and returns an error if there was a problem.
func (ppm *PPM) Save(filename string) error{
	file, err := os.Create(filename)
	if err != nil{
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n", ppm.MagicNumber)
	fmt.Fprintf(file, "%d %d\n", ppm.Width, ppm.Height)
	fmt.Fprintf(file, "%d\n", ppm.Max)

	switch ppm.MagicNumber{
	case "P3":
		for i := 0; i < ppm.Height; i++{
			for j := 0; j < ppm.Width; j++{
				fmt.Fprintf(file, "%d %d %d ", ppm.Data[i][j].R, ppm.Data[i][j].G, ppm.Data[i][j].B)
			}
			fmt.Fprintln(file)
		}
	case "P6":
		for i := 0; i < ppm.Height; i++{
			for j := 0; j < ppm.Width; j++{
				fmt.Fprintf(file, "%c%c%c", ppm.Data[i][j].R, ppm.Data[i][j].G, ppm.Data[i][j].B)
			}
			fmt.Fprintln(file)
		}
	default:
		return fmt.Errorf("invalid PPM file: unknown magic number")
	}

	return nil
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert(){
	for i := 0; i < ppm.Height; i++{
		for j := 0; j < ppm.Width; j++{
			ppm.Data[i][j].R = uint8(ppm.Max) - ppm.Data[i][j].R
			ppm.Data[i][j].G = uint8(ppm.Max) - ppm.Data[i][j].G
			ppm.Data[i][j].B = uint8(ppm.Max) - ppm.Data[i][j].B
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip(){
	for i := 0; i < ppm.Height; i++{
		for j := 0; j < ppm.Width/2; j++{
			ppm.Data[i][j], ppm.Data[i][ppm.Width-j-1] = ppm.Data[i][ppm.Width-j-1], ppm.Data[i][j]
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop(){
	for i := 0; i < ppm.Height/2; i++{
		for j := 0; j < ppm.Width; j++{
			ppm.Data[i][j], ppm.Data[ppm.Height-i-1][j] = ppm.Data[ppm.Height-i-1][j], ppm.Data[i][j]
		}
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(MagicNumber string){
	ppm.MagicNumber = MagicNumber
	switch ppm.MagicNumber{
	case "P3":
		for i := 0; i < ppm.Height; i++{
			for j := 0; j < ppm.Width; j++{
				ppm.Data[i][j].R = uint8(ppm.Max) - ppm.Data[i][j].R
				ppm.Data[i][j].G = uint8(ppm.Max) - ppm.Data[i][j].G
				ppm.Data[i][j].B = uint8(ppm.Max) - ppm.Data[i][j].B
			}
		}
	case "P6":
		for i := 0; i < ppm.Height; i++{
			for j := 0; j < ppm.Width; j++{
				ppm.Data[i][j].R = uint8(ppm.Max) - ppm.Data[i][j].R
				ppm.Data[i][j].G = uint8(ppm.Max) - ppm.Data[i][j].G
				ppm.Data[i][j].B = uint8(ppm.Max) - ppm.Data[i][j].B
			}
		}
	}
}

// SetMaxValue sets the Max value of the PPM image.
func (ppm *PPM) SetMaxValue(MaxValue uint8){
	ppm.Max = int(MaxValue)
}

// Rotate90CW rotates the PPM image 90° clockwise.
func (ppm *PPM) Rotate90CW(){
	ppm.Width, ppm.Height = ppm.Height, ppm.Width
	newData := make([][]Pixel, ppm.Height)
	for i := range newData{
		newData[i] = make([]Pixel, ppm.Width)
	}
	for i := 0; i < ppm.Height; i++{
		for j := 0; j < ppm.Width; j++{
			newData[i][j] = ppm.Data[ppm.Width-j-1][i]
		}
	}
	ppm.Data = newData
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM{
	var pgm PGM
	pgm.MagicNumber = "P2"
	pgm.Width = ppm.Width
	pgm.Height = ppm.Height
	pgm.Max = ppm.Max
	Data := make([][]uint8, pgm.Height)
	for i := range Data{
		Data[i] = make([]uint8, pgm.Width)
	}
	pgm.Data = Data
	for i := 0; i < pgm.Height; i++{
		for j := 0; j < pgm.Width; j++{
			pgm.Data[i][j] = uint8((int(ppm.Data[i][j].R) + int(ppm.Data[i][j].G) + int(ppm.Data[i][j].B))/3)
		}
	}
	return &pgm
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM{
	var pbm PBM
	pbm.MagicNumber = "P1"
	pbm.Width = ppm.Width
	pbm.Height = ppm.Height
	Data := make([][]bool, pbm.Height)
	for i := range Data{
		Data[i] = make([]bool, pbm.Width)
	}
	pbm.Data = Data
	for i := 0; i < pbm.Height; i++{
		for j := 0; j < pbm.Width; j++{
			if uint8((int(ppm.Data[i][j].R) + int(ppm.Data[i][j].G) + int(ppm.Data[i][j].B))/3) > uint8(ppm.Max/2){
				pbm.Data[i][j] = true
			}
		}
	}
	return &pbm
}

// DrawLine draws a line between two points.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel){
	if p1.X > p2.X{
		p1, p2 = p2, p1
	}
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	if dx == 0{
		for y := p1.Y; y <= p2.Y; y++{
			ppm.Set(p1.X, y, color)
		}
		return
	}
	for x := p1.X; x <= p2.X; x++{
		y := p1.Y + dy*(x-p1.X)/dx
		ppm.Set(x, y, color)
	}
}

// DrawRectangle draws a rectangle.
func (ppm *PPM) DrawRectangle(p1 Point, width , height int, color Pixel){
	ppm.DrawLine(p1, Point{p1.X + width, p1.Y}, color)
	ppm.DrawLine(p1, Point{p1.X, p1.Y + height}, color)
	ppm.DrawLine(Point{p1.X + width, p1.Y}, Point{p1.X + width, p1.Y + height}, color)
	ppm.DrawLine(Point{p1.X, p1.Y + height}, Point{p1.X + width, p1.Y + height}, color)
}

// DrawFilledRectangle draws a filled rectangle.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width , height int, color Pixel){
	for i := 0; i < height; i++{
		ppm.DrawLine(Point{p1.X, p1.Y + i}, Point{p1.X + width, p1.Y + i}, color)
	}
}

// DrawCircle draws a circle.
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel){
	x := 0
	y := radius
	d := 1 - radius
	for x <= y{
		ppm.Set(center.X + x, center.Y + y, color)
		ppm.Set(center.X + x, center.Y - y, color)
		ppm.Set(center.X - x, center.Y + y, color)
		ppm.Set(center.X - x, center.Y - y, color)
		ppm.Set(center.X + y, center.Y + x, color)
		ppm.Set(center.X + y, center.Y - x, color)
		ppm.Set(center.X - y, center.Y + x, color)
		ppm.Set(center.X - y, center.Y - x, color)
		if d < 0{
			d += 2*x + 3
		} else{
			d += 2*(x-y) + 5
			y--
		}
		x++
	}
}

// DrawFilledCircle draws a filled circle.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel){
	x := 0
	y := radius
	d := 1 - radius
	for x <= y{
		ppm.DrawLine(Point{center.X - x, center.Y + y}, Point{center.X + x, center.Y + y}, color)
		ppm.DrawLine(Point{center.X - x, center.Y - y}, Point{center.X + x, center.Y - y}, color)
		ppm.DrawLine(Point{center.X - y, center.Y + x}, Point{center.X + y, center.Y + x}, color)
		ppm.DrawLine(Point{center.X - y, center.Y - x}, Point{center.X + y, center.Y - x}, color)
		if d < 0{
			d += 2*x + 3
		} else{
			d += 2*(x-y) + 5
			y--
		}
		x++
	}
	for i := 0; i < radius; i++{
		ppm.DrawLine(Point{center.X - i, center.Y}, Point{center.X + i, center.Y}, color)
	}
}

// DrawTriangle draws a triangle.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel){
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p1, p3, color)
	ppm.DrawLine(p2, p3, color)
}

// DrawFilledTriangle draws a filled triangle.
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel){
	if p1.Y > p2.Y{
		p1, p2 = p2, p1
	}
	if p1.Y > p3.Y{
		p1, p3 = p3, p1
	}
	if p2.Y > p3.Y{
		p2, p3 = p3, p2
	}
	if p1.Y == p2.Y{
		ppm.DrawFilledTriangleFlatTop(p1, p2, p3, color)
	} else if p2.Y == p3.Y{
		ppm.DrawFilledTriangleFlatBottom(p1, p2, p3, color)
	} else{
		p4 := Point{p1.X + int(float64(p2.Y-p1.Y)/float64(p3.Y-p1.Y)*float64(p3.X-p1.X)), p2.Y}
		ppm.DrawFilledTriangleFlatBottom(p1, p2, p4, color)
		ppm.DrawFilledTriangleFlatTop(p2, p4, p3, color)
	}
}

// DrawFilledTriangleFlatTop draws a filled triangle with a flat top.
func (ppm *PPM) DrawFilledTriangleFlatTop(p1, p2, p3 Point, color Pixel){
	s1 := float64(p3.X-p1.X)/float64(p3.Y-p1.Y)
	s2 := float64(p3.X-p2.X)/float64(p3.Y-p2.Y)
	x1 := float64(p3.X)
	x2 := float64(p3.X)
	for y := p3.Y; y > p1.Y; y--{
		ppm.DrawLine(Point{int(x1), y}, Point{int(x2), y}, color)
		x1 -= s1
		x2 -= s2
	}
}

// DrawFilledTriangleFlatBottom draws a filled triangle with a flat bottom.
func (ppm *PPM) DrawFilledTriangleFlatBottom(p1, p2, p3 Point, color Pixel){
	s1 := float64(p2.X-p1.X)/float64(p2.Y-p1.Y)
	s2 := float64(p3.X-p1.X)/float64(p3.Y-p1.Y)
	x1 := float64(p1.X)
	x2 := float64(p1.X)
	for y := p1.Y; y <= p2.Y; y++{
		ppm.DrawLine(Point{int(x1), y}, Point{int(x2), y}, color)
		x1 += s1
		x2 += s2
	}
}

// DrawPolygon draws a polygon.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel){
	for i := 0; i < len(points)-1; i++{
		ppm.DrawLine(points[i], points[i+1], color)
	}
	ppm.DrawLine(points[len(points)-1], points[0], color)
}

// DrawFilledPolygon draws a filled polygon.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel){
	for i := 0; i < len(points)-1; i++{
		ppm.DrawFilledTriangle(points[0], points[i], points[i+1], color)
	}
}

// DrawKochSnowflake draws a Koch snowflake.
// N is the number of iterations.
// Koch snowflake is a 3 times a Koch curve.
// Start is the top point of the snowflake.
// Width is the width all the lines.
// Color is the color of the lines.
func (ppm *PPM) DrawKochSnowflake(n int, start Point,width int,color Pixel){
	ppm.DrawKochCurve(n, start, width, color)
	ppm.DrawKochCurve(n, Point{start.X + width/2, start.Y + int(float64(width)*0.866)}, width, color)
	ppm.DrawKochCurve(n, Point{start.X + width, start.Y}, width, color)
}

// DrawKochCurve draws a Koch curve.
// N is the number of iterations.
// Start is the top point of the curve.
// Width is the width of the curve.
// Color is the color of the curve.
func (ppm *PPM) DrawKochCurve(n int, start Point,width int,color Pixel){
	if n == 0{
		ppm.DrawLine(start, Point{start.X + width, start.Y}, color)
	} else{
		ppm.DrawKochCurve(n-1, start, width/3, color)
		ppm.DrawKochCurve(n-1, Point{start.X + width/3, start.Y}, width/3, color)
		ppm.DrawKochCurve(n-1, Point{start.X + width/3*2, start.Y}, width/3, color)
		ppm.DrawKochCurve(n-1, Point{start.X + width, start.Y}, width/3, color)
	}
}

// N is the number of iterations.
// Start is the top point of the triangle.
// Width is the width all the lines.
// Color is the color of the lines.
// DrawSierpinskiTriangle draws a Sierpinski triangle.
func (ppm *PPM) DrawSierpinskiTriangle(n int, start Point,width int,color Pixel){
	ppm.DrawSierpinskiCurve(n, start, width, color)
	ppm.DrawSierpinskiCurve(n, Point{start.X + width/2, start.Y + int(float64(width)*0.866)}, width, color)
	ppm.DrawSierpinskiCurve(n, Point{start.X + width, start.Y}, width, color)
}

// DrawSierpinskiCurve draws a Sierpinski curve.
// N is the number of iterations.
// Start is the top point of the curve.
// Width is the width of the curve.
// Color is the color of the curve.
func (ppm *PPM) DrawSierpinskiCurve(n int, start Point,width int,color Pixel){
	if n == 0{
		ppm.DrawLine(start, Point{start.X + width, start.Y}, color)
	} else{
		ppm.DrawSierpinskiCurve(n-1, start, width/2, color)
		ppm.DrawSierpinskiCurve(n-1, Point{start.X + width/2, start.Y}, width/2, color)
		ppm.DrawSierpinskiCurve(n-1, Point{start.X + width/4, start.Y + int(float64(width)*0.866/2)}, width/2, color)
	}
}

// DrawPerlinNoise draws perlin noise.
// this function Draw a perlin noise of all the image.
// Color1 is the color of 0.
// Color2 is the color of 1.  
func (ppm *PPM) DrawPerlinNoise(color1 Pixel , color2 Pixel){
	
}

// KNearestNeighbors resizes the PPM image using the k-nearest neighbors algorithm.
func (ppm *PPM) KNearestNeighbors(newWidth, newHeight int){
	newData := make([][]Pixel, newHeight)
	for i := range newData{
		newData[i] = make([]Pixel, newWidth)
	}
	for i := 0; i < newHeight; i++{
		for j := 0; j < newWidth; j++{
			newData[i][j] = ppm.Data[i*ppm.Height/newHeight][j*ppm.Width/newWidth]
		}
	}
	ppm.Width = newWidth
	ppm.Height = newHeight
	ppm.Data = newData
}