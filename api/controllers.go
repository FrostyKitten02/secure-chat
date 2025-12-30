package api

import (
	"github.com/danielgtaylor/huma/v2"
	"secure-chat/env"
)

var controllers []Controller = []Controller{authController}

type Controller func(api huma.API, options *env.Options)

func RegisterAll(api huma.API, options *env.Options) {
	for _, c := range controllers {
		c(api, options)
	}
}
