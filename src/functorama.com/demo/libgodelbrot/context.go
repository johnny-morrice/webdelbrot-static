package libgodelbrot

import (
    "image"
    "fmt"
)

type Renderer interface {
    Render() (*image.NRGBA, error)
}

func renderer(desc *Info) (Renderer, error) {
    // Check that numerics modes are okay, but do not act on them
    switch desc.NumericsStrategy {
    case NativeNumericsMode:
    case BigFloatNumericsMode:
    default:
        return nil, fmt.Errorf("Invalid NumericsStrategy: %v", desc.NumericsStrategy)
    }

    renderer := Renderer(nil)
    switch desc.RenderStrategy {
    case RegionRenderMode:
        renderer = makeSequenceFacade(desc)
    case SequenceRenderMode:
        renderer = makeRegionFacade(desc)
    case SharedRegionRenderMode:
        renderer = makeSharedRegionFacade(desc)
    default:
        return nil, fmt.Errorf("Invalid RenderStrategy: %v", desc.RenderStrategy)
    }

    return renderer, nil
}