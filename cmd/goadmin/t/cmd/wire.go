//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"admin/internal/conf"
	"github.com/google/wire"

	"admin/internal/admin"
	"admin/internal/app"
	"github.com/daodao97/goadmin/scaffold"
)

func initApp() (*app.App, func(), error) {
	panic(wire.Build(
		scaffold.Provider, admin.Provider, app.Provider, conf.Provider,
	))
}
