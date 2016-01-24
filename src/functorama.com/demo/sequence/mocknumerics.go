package sequence

import (
    "image"
    "functorama.com/demo/base"
)

type MockNumerics struct {
    TSequence bool
    TSubImage bool

    PointCount int
}

// Check MockNumerics implements SequenceNumerics interface
var _ SequenceNumerics = (*MockNumerics)(nil)

func (mn *MockNumerics) Sequence() []base.PixelMember {
    mn.TSequence = true

    return []base.PixelMember{base.PixelMember{}}
}

func (mn *MockNumerics) SubImage(rect image.Rectangle) {
    mn.TSubImage = true
}
