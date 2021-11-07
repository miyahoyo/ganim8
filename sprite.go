package ganim8

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteSize struct {
	W, H float64
}

// Sprite is a sprite that can be drawn to the screen.
// It can be animated by changing the current frame.
type Sprite struct {
	frames             []*image.Rectangle
	image              *ebiten.Image
	subImages          []*ebiten.Image
	size               SpriteSize
	length             int
	flippedH, flippedV bool
}

// NewSprite returns a new sprite.
func NewSprite(img *ebiten.Image, frames []*image.Rectangle) *Sprite {
	subImages := make([]*ebiten.Image, len(frames))
	for i, frame := range frames {
		subImages[i] = img.SubImage(*frame).(*ebiten.Image)
	}
	size := SpriteSize{0, 0}
	if len(frames) > 0 {
		size = SpriteSize{float64(frames[0].Dx()), float64(frames[0].Dy())}
	}
	return &Sprite{
		frames:    frames,
		image:     img,
		subImages: subImages,
		length:    len(frames),
		size:      size,
	}
}

// Size returns the size of the sprite.
func (spr *Sprite) Size() (float64, float64) {
	return spr.size.W, spr.size.H
}

// W is a shortcut for Size().X.
func (spr *Sprite) W() float64 {
	return spr.size.W
}

// H is a shortcut for Size().Y.
func (spr *Sprite) H() float64 {
	return spr.size.H
}

// IsEnd returns true if the current frame is the last frame.
func (spr *Sprite) IsEnd(index int) bool {
	return index >= spr.length-1
}

// FlipH flips the animation horizontally.
func (spr *Sprite) FlipH() {
	spr.flippedH = !spr.flippedH
}

// FlipV flips the animation horizontally.
func (spr *Sprite) FlipV() {
	spr.flippedV = !spr.flippedV
}

// Draw draws the current frame with the specified options.
func (spr *Sprite) Draw(screen *ebiten.Image, index int, opts *DrawOptions) {
	x, y := opts.X, opts.Y
	w, h := spr.Size()
	r := opts.Rotate
	ox, oy := opts.OriginX, opts.OriginY
	sx, sy := opts.ScaleX, opts.ScaleY

	op := &ebiten.DrawImageOptions{}
	op.ColorM = opts.ColorM
	op.CompositeMode = opts.CompositeMode

	if r != 0 {
		op.GeoM.Translate(-w*ox, -h*oy)
		op.GeoM.Rotate(r)
		op.GeoM.Translate(w*ox, h*oy)
	}

	if spr.flippedH {
		sx = sx * -1
	}
	if spr.flippedV {
		sy = sy * -1
	}

	if sx != 1 || sy != 1 {
		op.GeoM.Translate(-w*ox, -h*oy)
		op.GeoM.Scale(sx, sy)
		op.GeoM.Translate(w*ox, h*oy)
	}

	op.GeoM.Translate((x - w*ox), (y - h*oy))

	subImage := spr.subImages[index]
	screen.DrawImage(subImage, op)
}

// DrawWithShader draws the current frame with the specified options.
func (spr *Sprite) DrawWithShader(screen *ebiten.Image, index int, opts *DrawOptions, shaderOpts *ShaderOptions) {
	x, y := opts.X, opts.Y
	w, h := spr.Size()
	r := opts.Rotate
	ox, oy := opts.OriginX, opts.OriginY
	sx, sy := opts.ScaleX, opts.ScaleY

	op := &ebiten.DrawRectShaderOptions{}
	op.CompositeMode = opts.CompositeMode
	op.Uniforms = shaderOpts.Uniforms

	if r != 0 {
		op.GeoM.Translate(-w*ox, -h*oy)
		op.GeoM.Rotate(r)
		op.GeoM.Translate(w*ox, h*oy)
	}

	if spr.flippedH {
		sx = sx * -1
	}
	if spr.flippedV {
		sy = sy * -1
	}

	if sx != 1 || sy != 1 {
		op.GeoM.Translate(-w*ox, -h*oy)
		op.GeoM.Scale(sx, sy)
		op.GeoM.Translate(w*ox, h*oy)
	}

	op.GeoM.Translate((x - w*ox), (y - h*oy))

	subImage := spr.subImages[index]
	op.Images[0] = subImage
	for i := 0; i < 3; i++ {
		op.Images[i+1] = shaderOpts.Images[i]
	}
	screen.DrawRectShader(int(w), int(h), shaderOpts.Shader, op)
}