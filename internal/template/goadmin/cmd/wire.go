//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"

	"admin/internal/app"
	"admin/internal/service"
	"github.com/daodao97/goadmin/scaffold"
)

func initApp() (*app.App, func(), error) {
	panic(wire.Build(
		scaffold.Provider, service.Provider, app.Provider,
	))
}
