// Package assets bundles static files directly into the shellsound binary
// so it has no dependency on any path outside itself once compiled.
package assets

import _ "embed"

// FailureSound is the default sound played when a command fails.
//
//go:embed sounds/error.mp3
var FailureSound []byte

// SuccessSound is the default sound played when a command succeeds.
//
//go:embed sounds/normal.wav
var SuccessSound []byte
