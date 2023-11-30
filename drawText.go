package util

import (
	"image"
	"image/color"
	// "image/draw"
	// "time"

	// "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/font/basicfont"

	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	// "golang.org/x/mobile/exp/font"
)

// var f []byte = font.Default() 

func drawText(img *image.RGBA, x, y int, s string) {
    col := color.RGBA{0, 0, 0, 255}
    point := fixed.Point26_6{fixed.I(x), fixed.I(y)}

    d := &font.Drawer {
        Dst:  img,
        Src:  image.NewUniform(col),
        Face: basicfont.Face7x13,
        Dot:  point,
    }
    d.DrawString(s)
}

// Text draws text to screen
type Text struct {
	sz       size.Event
	images   *glutil.Images
	m        *glutil.Image
	// TODO: store *gl.Context
}

// NewText creates an Text tied to the current GL context.
func NewText(images *glutil.Images) *Text {
	return &Text{
		images:   images,
	}
}

// Draw draws text at the x, y coordinate and scaleX and scaleY specified by user
func (p *Text) Draw(sz size.Event, x, y int, scaleX, scaleY geom.Pt, s string) {
	if sz.WidthPx == 0 && sz.HeightPx == 0 {
		return
	}

	if p.sz != sz {
		p.sz = sz
		if p.m != nil {
			p.m.Release()
		}
		p.m = p.images.NewImage(sz.WidthPx, sz.HeightPx)
	}

	// clear image in p.m
	p.m.RGBA.Pix = make([]uint8, len(p.m.RGBA.Pix))

	drawText(p.m.RGBA, int(geom.Pt(x) / scaleX), int(geom.Pt(y) / scaleY), s)
	
	// copy img data to GL device
	p.m.Upload()

	p.m.Draw(
		sz,
		geom.Point{0, 0},	// topLeft	
		geom.Point{sz.WidthPt * scaleX, 0},	// topRight
		geom.Point{0, sz.HeightPt * scaleY},	// bottomLeft
		p.m.RGBA.Bounds(),
	)
}

func (t *Text) Release() {
	if t.m != nil {
		t.m.Release()
		t.m = nil
		t.images = nil
	}
}