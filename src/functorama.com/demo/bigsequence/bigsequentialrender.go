package bigsequence

import (
	"functorama.com/demo/base"
	"functorama.com/demo/bigbase"
	"functorama.com/demo/sequence"
)

type BigSequenceNumerics struct {
	bigbase.BigBaseNumerics
	area int
}

// Check that BigSequenceNumerics implements SequenceNumerics interface
var _ sequence.SequenceNumerics = (*BigSequenceNumerics)(nil)

func Make(app bigbase.RenderApplication) BigSequenceNumerics {
	w, h := app.PictureDimensions()
	return BigSequenceNumerics{
		BigBaseNumerics: bigbase.Make(app),
		area: int(w * h),
	}
}

func (bsn *BigSequenceNumerics) Area() int {
	return bsn.area
}


func (bsn *BigSequenceNumerics) Sequence() <-chan base.PixelMember {
	imageLeft, imageTop := bsn.PictureMin()
	imageRight, imageBottom := bsn.PictureMax()
	iterlim := bsn.IterateLimit

	out := make(chan base.PixelMember)

	go func() {
		pos := bigbase.BigComplex{
			R: bsn.RealMin,
		}
		for i := imageLeft; i < imageRight; i++ {
			pos.I = bsn.ImagMax
			for j := imageTop; j < imageBottom; j++ {
				member := bigbase.BigMandelbrotMember{
					C: &pos,
					SqrtDivergeLimit: &bsn.SqrtDivergeLimit,
				}
				member.Mandelbrot(iterlim)
				pos.I.Sub(&pos.I, &bsn.Iunit)
			}
			pos.R.Add(&pos.R, &bsn.Runit)
		}
		close(out)
	}()

	return out
}