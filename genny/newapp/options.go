package newapp

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gobuffalo/buffalo-docker/genny/docker"
	"github.com/gobuffalo/buffalo-pop/genny/newapp"
	"github.com/gobuffalo/buffalo/genny/assets/standard"
	"github.com/gobuffalo/buffalo/genny/assets/webpack"
	"github.com/gobuffalo/buffalo/genny/ci"
	"github.com/gobuffalo/buffalo/genny/refresh"
	"github.com/gobuffalo/buffalo/genny/vcs"
	"github.com/gobuffalo/buffalo/runtime"
	"github.com/gobuffalo/meta"
	"github.com/pkg/errors"
)

type Options struct {
	App            meta.App
	Docker         *docker.Options
	Pop            *newapp.Options
	CI             *ci.Options
	VCS            *vcs.Options
	Webpack        *webpack.Options
	Standard       *standard.Options
	Refresh        *refresh.Options
	Version        string
	ForbiddenNames []string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}

	if len(opts.Version) == 0 {
		opts.Version = runtime.Version
	}

	if opts.Docker != nil {
		if opts.Docker.App.IsZero() {
			opts.Docker.App = opts.App
		}
		if err := opts.Docker.Validate(); err != nil {
			return errors.WithStack(err)
		}
	}

	if opts.Pop != nil {
		if opts.Pop.App.IsZero() {
			opts.Pop.App = opts.App
		}
		if err := opts.Pop.Validate(); err != nil {
			return errors.WithStack(err)
		}
	}

	if opts.CI != nil {
		if opts.CI.App.IsZero() {
			opts.CI.App = opts.App
		}
		if err := opts.CI.Validate(); err != nil {
			return errors.WithStack(err)
		}
	}

	if opts.Webpack != nil {
		if opts.Webpack.App.IsZero() {
			opts.Webpack.App = opts.App
		}
		if err := opts.Webpack.Validate(); err != nil {
			return errors.WithStack(err)
		}
	}

	if opts.Standard != nil && opts.Webpack != nil {
		return errors.New("you can not use both webpack and standard generators")
	}

	if opts.Refresh != nil {
		if opts.Refresh.App.IsZero() {
			opts.Refresh.App = opts.App
		}
		if err := opts.Refresh.Validate(); err != nil {
			return errors.WithStack(err)
		}
	}

	if opts.VCS != nil {
		if opts.VCS.App.IsZero() {
			opts.VCS.App = opts.App
		}
		if err := opts.VCS.Validate(); err != nil {
			return errors.WithStack(err)
		}
	}

	if opts.App.WithModules && opts.App.WithDep {
		return ErrGoModulesWithDep
	}

	name := strings.ToLower(opts.App.Name.String())
	fb := append(opts.ForbiddenNames, "buffalo", "test", "dev")
	for _, n := range fb {
		rx, err := regexp.Compile(n)
		if err != nil {
			return errors.WithStack(err)
		}
		if rx.MatchString(name) {
			return fmt.Errorf("name %s is not allowed, try a different application name", opts.App.Name)
		}
	}

	if !nameRX.MatchString(name) {
		return fmt.Errorf("name %s is not allowed, application name can only contain [a-Z0-9-_]", opts.App.Name)
	}
	return nil
}

var nameRX = regexp.MustCompile(`^[\w-]+$`)
