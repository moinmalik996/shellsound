// Package assets bundles static files directly into the shellsound binary
// so it has no dependency on any path outside itself once compiled.
package assets

import _ "embed"

// FailureSound is the default sound played when a command fails.
//
//go:embed sounds/shellsound.mp3
var FailureSound []byte
