package botto

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	materialBlue      = color.RGBA{0x46, 0x88, 0xF1, 0xFF}
	materialRed       = color.RGBA{0xe8, 0x45, 0x3C, 0xFF}
	materialGray      = color.RGBA{0xb3, 0xb3, 0xB3, 0xFF}
	materialLightGray = color.RGBA{0xe3, 0xe3, 0xE3, 0xFF}
)

type textAlignment = int

const (
	alignLeft textAlignment = iota
	alignCenter
	alignRight
)

type customFontCache map[string]*truetype.Font

func (fc customFontCache) Store(fd draw2d.FontData, font *truetype.Font) {
	fc[fd.Name] = font
}

func (fc customFontCache) Load(fd draw2d.FontData) (*truetype.Font, error) {
	font, stored := fc[fd.Name]
	if !stored {
		return nil, fmt.Errorf("font %s is not stored in font cache", fd.Name)
	}
	return font, nil
}

func init() {
	fontCache := customFontCache{}

	TTFs := map[string]([]byte){
		"goregular": goregular.TTF,
	}

	for fontName, TTF := range TTFs {
		font, err := truetype.Parse(TTF)
		if err != nil {
			panic(err)
		}
		fontCache.Store(draw2d.FontData{Name: fontName}, font)
	}

	draw2d.SetFontCache(fontCache)
}

type Punchcard struct {
	canvas *image.RGBA
	inner  image.Rectangle
	gc     *draw2dimg.GraphicContext
}

func NewPunchcard() *Punchcard {
	canvas := image.NewRGBA(image.Rect(0, 0, 600, 1080))
	gc := draw2dimg.NewGraphicContext(canvas)

	gc.SetFontData(draw2d.FontData{
		Name:   "goregular",
		Family: draw2d.FontFamilyMono,
		Style:  draw2d.FontStyleBold | draw2d.FontStyleItalic,
	})

	pc := &Punchcard{
		canvas: canvas,
		inner:  image.Rect(32, 120, 540, 1020),
		gc:     gc,
	}

	pc.drawIntervals()

	return pc
}

func (pc *Punchcard) SetDay(date time.Time) {
	pc.gc.SetFontSize(22)

	text := date.Format("Monday, January 2")
	pc.withColor(materialBlue).rect(image.Rect(0, 0, 600, 60))
	pc.withColor(color.White).text(text, alignCenter, 300, 30)
}

func (pc *Punchcard) AddTask(text string, start, end time.Time) error {
	pc.gc.SetFontSize(16)

	if start.After(end) {
		return fmt.Errorf("start time is after end time")
	}

	x0 := pc.inner.Min.X + 96
	x1 := pc.inner.Max.X - 12

	y0 := pc.timeToYCoord(start)
	y1 := pc.timeToYCoord(end)

	if math.IsNaN(y0) || math.IsNaN(y1) {
		return fmt.Errorf("time is outside printable range")
	}

	pc.withColor(materialRed).roundedRect(image.Rect(x0, int(y0), x1, int(y1)), 10)
	pc.withColor(color.White).text(text, alignLeft, float64(x0+12), y0+24)

	return nil
}

func (pc *Punchcard) Rasterize() image.Image {
	return pc.canvas
}

func (pc *Punchcard) withColor(color color.Color) *Punchcard {
	pc.gc.SetFillColor(color)
	pc.gc.SetStrokeColor(color)
	return pc
}

func (pc *Punchcard) line(x0, y0, x1, y1 float64) {
	pc.gc.BeginPath()
	pc.gc.MoveTo(x0, y0)
	pc.gc.LineTo(x1, y1)
	pc.gc.Close()
	pc.gc.FillStroke()
}

func (pc *Punchcard) rect(rect image.Rectangle) {
	x0, y0 := float64(rect.Min.X), float64(rect.Min.Y)
	x1, y1 := float64(rect.Max.X), float64(rect.Max.Y)

	pc.gc.BeginPath()
	pc.gc.MoveTo(x0, y0)
	pc.gc.LineTo(x1, y0)
	pc.gc.LineTo(x1, y1)
	pc.gc.LineTo(x0, y1)
	pc.gc.LineTo(x0, y0)
	pc.gc.Close()
	pc.gc.FillStroke()
}

func (pc *Punchcard) roundedRect(rect image.Rectangle, radius float64) {
	x0, y0 := float64(rect.Min.X), float64(rect.Min.Y)
	x1, y1 := float64(rect.Max.X), float64(rect.Max.Y)

	pc.gc.BeginPath()
	pc.gc.MoveTo(x0+radius, y0)

	pc.gc.LineTo(x1-radius, y0)
	pc.gc.ArcTo(x1-radius, y0+radius, radius, radius, math.Pi*3/2, math.Pi/2)

	pc.gc.LineTo(x1, y1-radius)
	pc.gc.ArcTo(x1-radius, y1-radius, radius, radius, 0, math.Pi/2)

	pc.gc.LineTo(x0+radius, y1)
	pc.gc.ArcTo(x0+radius, y1-radius, radius, radius, math.Pi/2, math.Pi/2)

	pc.gc.LineTo(x0, y0+radius)
	pc.gc.ArcTo(x0+radius, y0+radius, radius, radius, math.Pi, math.Pi/2)

	pc.gc.Close()
	pc.gc.FillStroke()
}

func (pc *Punchcard) text(text string, align textAlignment, x float64, y float64) {
	l, t, r, b := pc.gc.GetStringBounds(text)
	w, h := math.Abs(r-l), math.Abs(t-b)

	y = y - t - (h / 2)

	switch align {
	case alignCenter:
		x = x - l - (w / 2)
	case alignRight:
		x = x - l + (w / 2)
	}

	pc.gc.FillStringAt(text, x, y)
}

func (pc *Punchcard) drawIntervals() {
	pc.gc.SetFontSize(16)

	pc.withColor(color.White).rect(pc.canvas.Rect)

	for h := 7; h < 19; h++ {
		x := float64(pc.inner.Min.X + 12)
		y := pc.timeToYCoord(time.Date(0, 0, 0, h, 0, 0, 0, time.Local))

		pc.withColor(materialGray).text(fmt.Sprintf("%02d:00", h), alignLeft, x, y)
		pc.withColor(materialLightGray).line(x+72, y, float64(pc.inner.Max.X), y)
	}
}

func (pc *Punchcard) timeToYCoord(t time.Time) float64 {
	start := 7 // 07:00 AM
	end := 19  // 07:00 PM

	delta := float64(end - start)
	step := float64(pc.inner.Dy()) / delta

	if t.Hour() < 7 || t.Hour() >= 19 {
		return math.NaN()
	}

	y := float64(pc.inner.Min.Y)
	y += step * float64(t.Hour()-start)
	y += step * (float64(t.Minute()) / 60.0)

	return 16.0 + y
}
