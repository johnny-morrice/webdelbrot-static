package libgodelbrot

import (
    "math/big"
)

type BigComplex struct {
    R *big.Float
    I *big.Float
}

func (c BigComplex) Real() FloatKind {
    return c.R
}

func (c BigComplex) Imag() FloatKind {
    return c.I
}

func (c BigComplex) Add(other BigComplex) {
    c.R.Add(other.R)
    c.I.Add(other.I)
}

func NewBigComplexNative(r float64, i float64, prec uint) BigComplex {
    return BigComplex{
        R: NewFloatPrec(r, prec),
        I: NewFloatPret(i, prec),
    }
}

// Use when you can assume accuracy is okay
func (f *big.Float) Float() float64 {
    native, acc := f.Float64()
    return native
}
 
// Create a new Float, supplying a precision
func NewFloatPrec(f float64, prec uint) *big.Float {
    b := NewFloat(f)
    b.SetPrec(prec)
    return b
}