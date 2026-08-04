// Package assets bundles static files directly into the fahh_alert binary
// so it has no dependency on any path outside itself once compiled.
package assets

import _ "embed"

// FailureSound is the default sound played when a command fails.
//
//go:embed sounds/fahhh.mp3
var FailureSound []byte
