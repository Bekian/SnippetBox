package ui

import (
	"embed"
)

// comment directive to embed the static file into the binary

//go:embed "html" "static"
var Files embed.FS
