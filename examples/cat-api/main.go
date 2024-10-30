package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	goatw "github.com/jwr1/goat/widget"

	"github.com/jwr1/goat"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type catApiResponse struct {
	ID     string `json:"id"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type app struct {
	goat.Widget
}

var _ goat.StateWidget = app{}

func (w app) Build() (goat.Widget, error) {
	cat, setCat := goat.UseState(catApiResponse{})
	fetchCatError, setFetchCatError := goat.UseState("")

	fetchCat := func() {
		r, err := http.Get("https://api.thecatapi.com/v1/images/search")
		if err != nil {
			setFetchCatError(err.Error())
			return
		}
		defer r.Body.Close()

		catSearchResponse := &[]catApiResponse{}
		err = json.NewDecoder(r.Body).Decode(catSearchResponse)
		if err != nil {
			setFetchCatError(err.Error())
			return
		}

		setCat((*catSearchResponse)[0])
	}

	goat.UseEffect(func() func() {
		fetchCat()
		return nil
	}, []any{})

	// If an error has occurred during fetchCat(), then display it.
	if fetchCatError != "" {
		return goatw.Center{
			Child: goatw.Text{
				Text: fetchCatError,
			},
		}, nil
	}

	// If still loading, then display loading indicator
	if cat.ID == "" {
		return goatw.Center{
			Child: goatw.Text{
				Text: "Loading cat...",
			},
		}, nil
	}

	return goatw.Column{
		// CrossAxisAlignment: goatw.CrossAxisAlignmentCenter,
		Children: []goat.Widget{
			goatw.Button{
				Label:         "[space] Press for new cat!",
				Padding:       goat.EdgeInsertsSymmetric(1, 2),
				OnActivate:    fetchCat,
				KeyActivators: []string{" "},
			},
			goatw.Text{Text: cat.Url},
			goatw.Text{Text: fmt.Sprintf("Original size: %dx%d", cat.Width, cat.Height)},
			goatw.ImageNetwork{
				Url: cat.Url,
				LoadingBuilder: func() goat.Widget {
					return goatw.Text{
						Text: "Loading cat...",
					}
				},
				ErrorBuilder: func(err error) goat.Widget {
					return goatw.Text{
						Text: err.Error(),
					}
				},
			},
		},
	}, nil
}

func main() {
	err := goat.RunApp(app{})
	if err != nil {
		panic(err.Error())
	}
}
