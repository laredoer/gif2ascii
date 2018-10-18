package gif

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/golang/freetype"
)

var (
	path, output string
	delay        int
	dpi          = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile     = flag.String("fontfile", "./testdata/luxisr.ttf", "filename of the ttf font")
	hinting      = flag.String("hinting", "none", "none | full")
	size         = flag.Float64("size", 1, "font size in points")
	spacing      = flag.Float64("spacing", 1, "line spacing (e.g. 2 means double spaced)")
	wonb         = flag.Bool("whiteonblack", false, "white text on a black background")
)

//GetEach 获取每一帧
func GetEach(filePath string) []image.Image {

	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	//gif解码
	g, err := gif.DecodeAll(f)
	if err != nil {
		panic(err)
	}
	//获取每一帧图片数据
	return Create(g.Image, g.Config.Width, g.Config.Height)
}

//Create .
func Create(p []*image.Paletted, w, h int) (i []image.Image) {
	for _, v := range p {
		i = append(i, v.SubImage(image.Rect(0, 0, w, h)))
	}
	return i
}

//Ascllimage 图片转为字符画（简易版）
func Ascllimage(m image.Image) (str string) {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	//arr := []string{"M", "N", "H", "Q", "$", "O", "C", "?", "7", ">", "!", ":", "~", ";", "."}   //粗糙
	//arr := []string{"M", "N", "H", "Q", "$", "O", "C", "?", "*", ">", "!", ":", "-", ";", "."}	//粗糙
	arr := []string{"M", "N", "D", "8", "O", "Z", "$", "7", "I", "?", "+", "=", "~", ":", "."} //细腻
	for i := 0; i < dy; i++ {
		for j := 0; j < dx; j++ {
			colorRgb := m.At(j, i)
			_, g, _, _ := colorRgb.RGBA()
			avg := uint8(g >> 8)
			num := avg / 18
			str = strings.Join([]string{str, arr[num]}, "")
			if j == dx-1 {
				str = strings.Join([]string{str, "\n"}, "")
			}
		}
	}
	return
}

//DrawImg is
func DrawImg(str string, k int) {
	txt := []byte(str)
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fg, bg := image.Black, image.White
	ruler := color.RGBA{0x22, 0x22, 0x22, 0xff}

	my := strings.Split(str, "\n")
	pngY := len(my) * int(*size)
	pngX := len([]byte(my[1])) * int(*size)
	rgba := image.NewRGBA(image.Rect(0, 0, pngX, pngY))

	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	c := freetype.NewContext()

	c.SetDPI(float64(*dpi))
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	rgba.Set(0, 0, ruler)

	pt := freetype.Pt(1, int(c.PointToFixed(*size)>>6))

	for _, s := range txt {
		_, err = c.DrawString(string(s), pt)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		pt.X += c.PointToFixed(*size)
		if string(s) == "\n" {
			pt.X = 1
			pt.Y += c.PointToFixed(*size)
		}

	}
	PathExists("tmp")
	outFile, err := os.Create(fmt.Sprintf("./tmp/%d.png", k))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	//
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func PathExists(path string) {
	_, err := os.Stat(path)
	if err == nil {

	}
	if os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func CreateGif() {

	path = "./tmp"
	output = "output.gif"
	delay = 1
	if path == "" {
		fmt.Println("请输入图片路径")
		flag.PrintDefaults()
		return
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	anim := gif.GIF{}
	for i := 0; i < len(files); i++ {
		f, err := os.Open(fmt.Sprintf(path+"/%d.png", i))
		if err != nil {
			fmt.Printf("Could not open file %d.png. Error: %s\n", i, err)
			return
		}
		defer f.Close()
		img, _, _ := image.Decode(f)

		paletted := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.FloydSteinberg.Draw(paletted, img.Bounds(), img, image.ZP)

		anim.Image = append(anim.Image, paletted)
		anim.Delay = append(anim.Delay, delay*15)
	}

	f, _ := os.Create(output)
	defer f.Close()
	gif.EncodeAll(f, &anim)
}
