package web

import (
	"html/template"

	"github.com/gobuffalo/buffalo/genny/newapp/core"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Group, error) {
	if err := opts.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	gg, err := core.New(opts.Options)
	if err != nil {
		return gg, errors.WithStack(err)
	}

	g := genny.New()
	data := map[string]interface{}{
		"opts": opts,
	}

	helpers := template.FuncMap{}

	t := gotools.TemplateTransformer(data, helpers)
	g.Transformer(t)
	g.Box(packr.NewBox("../web/templates"))

	gg.Add(g)

	return gg, nil
}
