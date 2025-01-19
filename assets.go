package main

import (
	_ "embed"
)

// Embed the image and audio assets.

//go:embed assets/audio/footstep_carpet_000.wav
var step1 []byte

//go:embed assets/audio/footstep_carpet_001.wav
var step2 []byte

//go:embed assets/audio/footstep_carpet_002.wav
var step3 []byte

//go:embed assets/img/1.png
var goLangImage []byte

//go:embed assets/img/2.png
var gopherImage []byte

//go:embed assets/img/3.png
var thesimpledevImage []byte

//go:embed assets/img/4.png
var aptImage []byte
