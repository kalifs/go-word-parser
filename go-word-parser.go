package main

import (
	"github.com/kalifs/go-word-parser/internal/transformer"
)

func main() {
	transformer.TransformLatvianDefinitions()
	transformer.TransformEngilishDefinitions()
}
