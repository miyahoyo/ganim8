package ganim8

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/rand"
)

type SpriteSize struct {
	W, H int
}

type SpriteSizeF struct {
	W, H float64
}

// Sprite is a sprite that can be drawn to the screen.
// It can be animated by changing the current frame.
type Sprite struct {
	frames             []*image.Rectangle
	image              *ebiten.Image
	subImages          []*ebiten.Image
	size               SpriteSize
	sizeF              SpriteSizeF
	length             int
	flippedH, flippedV bool
	op                 *ebiten.DrawImageOptions
	shaderOp           *ebiten.DrawRectShaderOptions
}

// NewSprite returns a new sprite.
func NewSprite(img *ebiten.Image, frames []*image.Rectangle) *Sprite {
	subImages := make([]*ebiten.Image, len(frames))
	for i, frame := range frames {
		subImages[i] = img.SubImage(*frame).(*ebiten.Image)
	}
	size := SpriteSize{0, 0}
	sizeF := SpriteSizeF{0, 0}
	if len(frames) > 0 {
		size = SpriteSize{frames[0].Dx(), frames[0].Dy()}
		sizeF = SpriteSizeF{float64(frames[0].Dx()), float64(frames[0].Dy())}
	}
	return &Sprite{
		frames:    frames,
		image:     img,
		subImages: subImages,
		length:    len(frames),
		size:      size,
		sizeF:     sizeF,
		op:        &ebiten.DrawImageOptions{},
		shaderOp:  &ebiten.DrawRectShaderOptions{},
	}
}

// Size returns the size of the sprite.
func (spr *Sprite) Size() (int, int) {
	return spr.size.W, spr.size.H
}

// Width returns the width of the sprite.
func (spr *Sprite) Width() int {
	return spr.size.W
}

// W is a shortcut for Width().
func (spr *Sprite) W() int {
	return spr.Width()
}

// Height returns the height of the sprite.
func (spr *Sprite) Height() int {
	return spr.size.H
}

// H is a shortcut for Height().
func (spr *Sprite) H() int {
	return spr.Height()
}

// Length returns the number of frames.
func (spr *Sprite) Length() int {
	return spr.length
}

// RandomIndex returns random index of the sprite
func (spr *Sprite) RandomIndex() int {
	return rand.Intn(spr.length)
}

// LoopIndex returns loop index of the sprite.
func (spr *Sprite) LoopIndex(index int) int {
	if index >= spr.length {
		return index % spr.length
	}
	return index
}

// IsEnd returns true if the current frame is the last frame.
func (spr *Sprite) IsEnd(index int) bool {
	return index >= spr.length-1
}

// FlipH flips the sprite horizontally.
func (spr *Sprite) FlipH() {
	spr.flippedH = !spr.flippedH
}

// FlipV flips the sprite vertically.
func (spr *Sprite) FlipV() {
	spr.flippedV = !spr.flippedV
}

// SetFlipH flips the sprite horizontally.
func (spr *Sprite) SetFlipH(flipH bool) {
	spr.flippedH = flipH
}

// SetFlipV flips the sprite vertically.
func (spr *Sprite) SetFlipV(flipV bool) {
	spr.flippedV = flipV
}

// Draw draws the current frame with the specified options.
func (spr *Sprite) Draw(screen *ebiten.Image, index int, opts *DrawOptions) {
	x, y := opts.X, opts.Y
	w, h := spr.sizeF.W, spr.sizeF.H
	r := opts.Rotate
	ox, oy := opts.OriginX, opts.OriginY
	sx, sy := opts.ScaleX, opts.ScaleY

	op := spr.op
	op.GeoM.Reset()
	op.ColorM = opts.ColorM
	op.CompositeMode = opts.CompositeMode

	if spr.flippedH {
		sx = sx * -1
		ox = 1 - ox
	}
	if spr.flippedV {
		sy = sy * -1
		oy = 1 - oy
	}

	if sx != 1 || sy != 1 {
		op.GeoM.Translate(-w*ox, -h*oy)
		op.GeoM.Scale(sx, sy)
		op.GeoM.Translate(w*ox, h*oy)
	}

	if r != 0 {
		op.GeoM.Translate(-w*ox, -h*oy)
		op.GeoM.Rotate(r)
		op.GeoM.Translate(w*ox, h*oy)
	}

	op.GeoM.Translate((x - w*ox), (y - h*oy))

	subImage := spr.subImages[index]
	screen.DrawImage(subImage, op)
}

// DrawWithShader draws the current frame with the specified options.
func (spr *Sprite) DrawWithShader(screen *ebiten.Image, index int, opts *DrawOptions, shaderOpts *ShaderOptions) {
	x, y := opts.X, opts.Y
	w, h := spr.sizeF.W, spr.sizeF.H
	r := opts.Rotate
	ox, oy := opts.OriginX, opts.OriginY
	sx, sy := opts.ScaleX, opts.ScaleY

	op := spr.shaderOp
	op.GeoM.Reset()
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
	op.Images[1] = shaderOpts.Images[0]
	op.Images[2] = shaderOpts.Images[1]
	op.Images[3] = shaderOpts.Images[2]
	screen.DrawRectShader(int(w), int(h), shaderOpts.Shader, op)
}

func (spr *Sprite) Clone() *Sprite {
	s := *spr
	s.op = &ebiten.DrawImageOptions{}
	s.shaderOp = &ebiten.DrawRectShaderOptions{}
	return &s
}
