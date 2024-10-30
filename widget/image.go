package goatw

import (
	"image"
	"net/http"

	. "github.com/jwr1/goat"
	"golang.org/x/image/draw"
)

type Image struct {
	Widget

	image image.Image
}

var _ RenderWidget = Image{}

func (w Image) Layout(context LayoutContext) (Size, error) {
	return Size{Width: context.Constraints.Max.Width, Height: context.Constraints.Max.Height}, nil
}

func (w Image) Paint(context PaintContext) error {
	dst := image.NewRGBA(image.Rect(0, 0, context.Size.Width, context.Size.Height))
	draw.NearestNeighbor.Scale(dst, dst.Rect, w.image, w.image.Bounds(), draw.Over, nil)

	context.Canvas.OverlayImage(0, 0, dst)

	return nil
}

type ImageNetwork struct {
	Widget

	Url            string
	LoadingBuilder func() Widget
	ErrorBuilder   func(err error) Widget
}

var _ StateWidget = ImageNetwork{}

func (w ImageNetwork) Build() (Widget, error) {
	data, setData := UseState[*image.Image](nil)
	err, setErr := UseState[*error](nil)

	UseEffect(func() func() {
		cleanup := func() {
			setData(nil)
			setErr(nil)
		}

		resp, err := http.Get(w.Url)
		if err != nil {
			setErr(&err)
			return cleanup
		}

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			setErr(&err)
			return cleanup
		}

		setData(&img)

		return cleanup
	}, []any{w.Url})

	if err != nil {
		return w.ErrorBuilder(*err), nil
	}

	if data == nil {
		return w.LoadingBuilder(), nil
	}

	return Image{
		image: *data,
	}, nil
}
