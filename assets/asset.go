package assets

import (
	_ "embed"
)

//go:embed index.min.html
var indexHTML []byte

//go:embed main.min.js
var mainJS []byte

func IndexHTML() []byte {
	return indexHTML
}

func MainJS() []byte {
	return mainJS
}
