package exp

import (
	"image/draw"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"golang.org/x/image/math/fixed"
	"log"
	"github.com/x6uderek/expgraphy/paint"
)

type Graphy interface {
	Measure(fontSize float64) (fixed.Point26_6)
	Draw(dst draw.Image, dot fixed.Point26_6)
}

var (
	DefaultFont *truetype.Font
	DefaultColor *image.Uniform
)

const DPI = 72.0

func FloatToFixed(f float64) fixed.Int26_6 {
	return fixed.Int26_6(f * DPI * 64.0 / 72.0)
}

func FixedToFloat(f fixed.Int26_6) float64 {
	return float64(f>>6) + float64(f&(1<<6 - 1))/64
}

func init() {
	defaultFont,err := truetype.Parse(gomono.TTF)
	if err != nil {
		log.Fatalf("parse font err: %v", err)
	}
	DefaultFont = defaultFont
	DefaultColor = image.NewUniform(color.Black)
}

//----------------------------------------------------------------------------
type VarGraphy struct {  //simple var, numbers
	Name string

	face font.Face
}

func (v *VarGraphy)Measure(fontSize float64) (fixed.Point26_6) {
	v.face = truetype.NewFace(DefaultFont, &truetype.Options{Size:fontSize})
	w := font.MeasureString(v.face, v.Name)
	return fixed.Point26_6{X:w, Y:FloatToFixed(fontSize)}
}

func (v *VarGraphy)Draw(dst draw.Image, dot fixed.Point26_6) {
	drawer := &font.Drawer{
		Dst:dst,
		Src:DefaultColor,
		Face:v.face,
		Dot: dot,
		}
	drawer.DrawString(v.Name)
}

//----------------------------------------------------------------------------
/*type NumberGraphy struct {
	Number string
}*/

//----------------------------------------------------------------------------
type UnaryGraphy struct {
	Op string
	X Graphy

	measureSize fixed.Point26_6
	xDot fixed.Point26_6 //x偏移
	postDot fixed.Point26_6 //后半个)的偏移
	face font.Face
	wrap bool
}

func (u *UnaryGraphy)Measure(fontSize float64) (fixed.Point26_6) {
	switch u.Op {
	case "+":
		u.measureSize = u.X.Measure(fontSize)
		u.xDot = fixed.Point26_6{}
		u.face = truetype.NewFace(DefaultFont, &truetype.Options{Size:FixedToFloat(u.measureSize.Y)})
		u.wrap = false
	case "-":
		childSize := u.X.Measure(fontSize)
		u.measureSize.Y = childSize.Y
		u.face = truetype.NewFace(DefaultFont, &truetype.Options{Size:FixedToFloat(u.measureSize.Y)})
		switch x := u.X.(type) {
		case *BinaryGraphy:
			if x.Op == "+" || x.Op == "-" {
				preW := font.MeasureString(u.face, "-(")
				postW := font.MeasureString(u.face, ")")
				u.measureSize.X = preW + childSize.X + postW
				u.xDot = fixed.Point26_6{X:preW, Y:0}
				u.postDot = fixed.Point26_6{X:preW + childSize.X, Y:0}
				u.wrap = true
			} else {
				preW := font.MeasureString(u.face, "-")
				u.measureSize.X = preW + childSize.X
				u.xDot = fixed.Point26_6{X:preW, Y:0}
				u.wrap = false
			}
		default:
			preW := font.MeasureString(u.face, "-")
			u.measureSize.X = preW + childSize.X
			u.xDot = fixed.Point26_6{X:preW, Y:0}
			u.wrap = false
		}
	}
	return u.measureSize
}

func (u *UnaryGraphy)Draw(dst draw.Image, dot fixed.Point26_6) {
	switch u.Op {
	case "+":
		u.X.Draw(dst, dot)
	case "-":
		drawer := &font.Drawer{
			Dst:dst,
			Src:DefaultColor,
			Face:u.face,
			Dot: dot,
		}
		if u.wrap {
			drawer.DrawString("-(")
			u.X.Draw(dst, dot.Add(u.xDot))
			drawer.Dot = dot.Add(u.postDot)
			drawer.DrawString(")")
		} else {
			drawer.DrawString("-")
			u.X.Draw(dst, dot.Add(u.xDot))
		}
	}
}

//----------------------------------------------------------------------------
type BinaryGraphy struct {
	Op string  //"+"  "-" "*"
	X, Y Graphy

	measureSize fixed.Point26_6
	xDot,yDot,opDot,bracketsDot fixed.Point26_6 //偏移
	xWrap, yWrap bool
	face font.Face
}

func (b *BinaryGraphy)Measure(fontSize float64) (fixed.Point26_6) {
	xSize := b.X.Measure(fontSize)
	ySize := b.Y.Measure(fontSize)
	if xSize.Y > ySize.Y {
		b.measureSize.Y = xSize.Y
	} else {
		b.measureSize.Y = ySize.Y
	}
	b.face = truetype.NewFace(DefaultFont, &truetype.Options{Size:FixedToFloat(b.measureSize.Y)})
	leftBracketsW := font.MeasureString(b.face, "(")
	rightBracketsW := font.MeasureString(b.face, ")")
	opW := font.MeasureString(b.face, b.Op)
	switch b.Op {
	case "+","-":
		b.xWrap = false
		b.yWrap = false
	case "*":
		switch x := b.X.(type) {
		case *BinaryGraphy:
			if x.Op == "+" || x.Op == "-" {
				b.xWrap = true
			} else {
				b.xWrap = false
			}
		default:
			b.xWrap = false
		}

		switch y := b.Y.(type) {
		case *BinaryGraphy:
			if y.Op == "+" || y.Op == "-" {
				b.yWrap = true
			} else {
				b.yWrap = false
			}
		default:
			b.yWrap = false
		}
	}
	if b.xWrap {
		b.xDot = fixed.Point26_6{X:leftBracketsW}
		b.opDot = fixed.Point26_6{X:leftBracketsW + xSize.X}
		if b.yWrap {
			b.yDot = fixed.Point26_6{X:leftBracketsW + xSize.X + rightBracketsW + opW + leftBracketsW}
			b.bracketsDot = fixed.Point26_6{X:leftBracketsW + xSize.X + rightBracketsW + opW + leftBracketsW + ySize.X}
			b.measureSize.X = leftBracketsW + xSize.X + rightBracketsW + opW + leftBracketsW + ySize.X + rightBracketsW
		} else {
			b.yDot = fixed.Point26_6{X:leftBracketsW + xSize.X + rightBracketsW + opW}
			b.measureSize.X = leftBracketsW + xSize.X + rightBracketsW + opW + ySize.X
		}
	} else {
		b.opDot = fixed.Point26_6{X:xSize.X}
		if b.yWrap {
			b.yDot = fixed.Point26_6{X:xSize.X + opW + leftBracketsW}
			b.bracketsDot = fixed.Point26_6{X:xSize.X + opW + leftBracketsW + ySize.X}
			b.measureSize.X = xSize.X + opW + leftBracketsW + ySize.X + rightBracketsW
		} else {
			b.yDot = fixed.Point26_6{X:xSize.X + opW}
			b.measureSize.X = xSize.X + opW + ySize.X
		}
	}
	return b.measureSize
}
func (b *BinaryGraphy) Draw(dst draw.Image, dot fixed.Point26_6) {
	drawer := &font.Drawer{
		Dst:dst,
		Src:DefaultColor,
		Face:b.face,
		Dot:dot,
	}
	opStr := b.Op
	if b.xWrap {
		drawer.DrawString("(")
		opStr = ")" + opStr
	}
	if b.yWrap {
		drawer.Dot = dot.Add(b.bracketsDot)
		drawer.DrawString(")")
		opStr = opStr + "("
	}
	drawer.Dot = dot.Add(b.opDot)
	drawer.DrawString(opStr)
	b.X.Draw(dst, dot.Add(b.xDot))
	b.Y.Draw(dst, dot.Add(b.yDot))
}

//----------------------------------------------------------------------------
type DivideGraphy struct { // "÷"
	X Graphy
	Y Graphy

	measureSize fixed.Point26_6
	xDot,yDot fixed.Point26_6
	lineW fixed.Int26_6
}

func (d *DivideGraphy)Measure(fontSize float64) (fixed.Point26_6) {
	xSize := d.X.Measure(fontSize*4/9)
	ySize := d.Y.Measure(fontSize*4/9)
	if xSize.X > ySize.X {
		d.measureSize.X = xSize.X*9/8
	} else {
		d.measureSize.X = ySize.X*9/8
	}
	d.measureSize.Y = (xSize.Y + ySize.Y)*9/8
	d.yDot = fixed.Point26_6{X:(d.measureSize.X - ySize.X)/2, Y:0}
	d.xDot = fixed.Point26_6{X:(d.measureSize.X - xSize.X)/2, Y:-d.measureSize.Y*5/9}
	d.lineW = d.measureSize.Y/40
	return d.measureSize
}

func (d *DivideGraphy)Draw(dst draw.Image, dot fixed.Point26_6) {
	d.X.Draw(dst, dot.Add(d.xDot))
	d.Y.Draw(dst, dot.Add(d.yDot))
	lineStart := dot.Add(fixed.Point26_6{Y:-d.measureSize.Y*4/9})
	lienEnd := lineStart.Add(fixed.Point26_6{X:d.measureSize.X})
	DrawLineFixPoint(dst, lineStart, lienEnd, FixedToFloat(d.lineW))
}

//----------------------------------------------------------------------------
type PowGraphy struct {
	X Graphy
	Y Graphy

	measureSize fixed.Point26_6
	xWrap bool
	xDot fixed.Point26_6
	yDot fixed.Point26_6
	faceBrackets font.Face
	dotBrackets fixed.Point26_6
}

func (p *PowGraphy) Measure(fontSize float64) (fixed.Point26_6) {
	xSize := p.X.Measure(fontSize)//x使用同等字体，因为普通变量x可以看做x的一次方，是平级的
	ySize := p.Y.Measure(fontSize*3/5)
	p.faceBrackets = truetype.NewFace(DefaultFont, &truetype.Options{Size:fontSize})
	switch x := p.X.(type) {
	case *UnaryGraphy:
		if x.Op == "-" {
			p.xWrap = true
		}
	case *BinaryGraphy:
		p.xWrap = true
	case *DivideGraphy:
		p.xWrap = true
	case *PowGraphy,*TriangleGraphy:// 暂时不支持(sin(x))^2把指数放到sin后面的写法，只能先用括号表示
		p.xWrap = true
	default:
		p.xWrap = false
	}
	preW := font.MeasureString(p.faceBrackets, "(")
	postW := font.MeasureString(p.faceBrackets, ")")
	if p.xWrap {
		p.xDot = fixed.Point26_6{X:preW,Y:0}
		p.yDot = fixed.Point26_6{X:preW + xSize.X + postW, Y:-xSize.Y*3/5}
		p.dotBrackets = fixed.Point26_6{X:preW+xSize.X, Y:0}
		p.measureSize = fixed.Point26_6{X:preW+xSize.X+postW+ySize.X, Y:xSize.Y+ySize.Y}
	} else {
		p.xDot = fixed.Point26_6{}
		p.yDot = fixed.Point26_6{X:xSize.X, Y:-xSize.Y*3/5}
		p.measureSize = fixed.Point26_6{X:xSize.X+ySize.X, Y:xSize.Y+ySize.Y}
	}
	return p.measureSize
}

func (p *PowGraphy)Draw(dst draw.Image, dot fixed.Point26_6) {
	drawer := &font.Drawer{
		Dst:dst,
		Src:DefaultColor,
		Face:p.faceBrackets,
		Dot:dot,
	}
	if p.xWrap {
		drawer.DrawString("(")
		drawer.Dot = dot.Add(p.dotBrackets)
		drawer.DrawString(")")
	}
	p.X.Draw(dst, dot.Add(p.xDot))
	p.Y.Draw(dst, dot.Add(p.yDot))
}

//----------------------------------------------------------------------------
type TriangleGraphy struct {
	Op string  //sin cos tan arcsin ...
	X Graphy

	measureSize fixed.Point26_6
	xDot fixed.Point26_6
	bracketsDot fixed.Point26_6
	face font.Face
}

func (t *TriangleGraphy)Measure(fontSize float64) (fixed.Point26_6) {
	xSize := t.X.Measure(fontSize)
	t.face = truetype.NewFace(DefaultFont, &truetype.Options{Size:fontSize})
	preW := font.MeasureString(t.face, t.Op+"(")
	postW := font.MeasureString(t.face, "(")
	t.xDot = fixed.Point26_6{X:preW, Y:0}
	t.bracketsDot = fixed.Point26_6{X:preW + xSize.X}
	t.measureSize = fixed.Point26_6{X:preW + xSize.X + postW, Y:xSize.Y}
	return t.measureSize
}

func (t *TriangleGraphy)Draw(dst draw.Image, dot fixed.Point26_6) {
	drawer := &font.Drawer{
		Dst:dst,
		Src:DefaultColor,
		Face:t.face,
		Dot:dot,
	}
	drawer.DrawString(t.Op+"(")
	drawer.Dot = dot.Add(t.bracketsDot)
	drawer.DrawString(")")
	t.X.Draw(dst, dot.Add(t.xDot))
}

//----------------------------------------------------------------------------
type SqrtGraphy struct {
	X Graphy

	measureSize fixed.Point26_6
	xDot fixed.Point26_6
	lineWidth float64
	preW fixed.Int26_6
}

func (s *SqrtGraphy)Measure(fontSize float64) (fixed.Point26_6) {
	xSize := s.X.Measure(fontSize)
	s.lineWidth = fontSize/30
	s.preW = FloatToFixed(fontSize/3)
	s.xDot = fixed.Point26_6{X: s.preW, Y:0}
	s.measureSize = fixed.Point26_6{X: s.preW+xSize.X, Y: xSize.Y*5/4}
	return s.measureSize
}

func (s *SqrtGraphy)Draw(dst draw.Image, dot fixed.Point26_6) {
	s.X.Draw(dst, dot.Add(s.xDot))
	lineStart := dot.Add(fixed.Point26_6{X:s.xDot.X, Y:-s.measureSize.Y*7/10})
	lineEnd := lineStart.Add(fixed.Point26_6{X:s.measureSize.X-s.preW})
	DrawLineFixPoint(dst, lineStart, lineEnd, s.lineWidth)
	lineEnd = dot.Add(fixed.Point26_6{X:s.preW*3/5, Y:s.measureSize.Y/5})
	DrawLineFixPoint(dst, lineStart, lineEnd, s.lineWidth)
	lineStart = dot.Add(fixed.Point26_6{X:s.preW/10,Y:-s.measureSize.Y*3/10})
	DrawLineFixPoint(dst, lineStart, lineEnd, s.lineWidth)
	lineEnd = dot.Add(fixed.Point26_6{Y:-s.measureSize.Y/10})
	DrawLineFixPoint(dst, lineStart, lineEnd, s.lineWidth)
}

func DrawLineFixPoint(dst draw.Image, from, to fixed.Point26_6, width float64) {
	paint.DrawLine(dst, FixedToFloat(from.X), FixedToFloat(from.Y), FixedToFloat(to.X), FixedToFloat(to.Y), width, DefaultColor)
}