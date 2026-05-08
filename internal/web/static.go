package web

import (
	"embed"
	"io/fs"
)

//go:embed all:static
var staticEmbed embed.FS

func StaticFS() fs.FS {
	sub, err := fs.Sub(staticEmbed, "static")
	if err != nil {
		panic(err)
	}
	return sub
}
