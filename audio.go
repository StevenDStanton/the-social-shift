package main

// import (
// 	"bytes"
// 	"log"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/audio"
// 	"github.com/hajimehoshi/ebiten/v2/audio/wav"
// )

// func (g *Game) initializeAudio() error {
// 	if g.audioContext != nil && g.audioPlayer != nil {
// 		return nil // Already initialized
// 	}

// 	// Initialize audio context
// 	g.audioContext = audio.NewContext(44100)
// 	log.Println("Audio context initialized")

// 	// Load the embedded audio
// 	audioBuffer := bytes.NewReader(stampWav)
// 	wavStream, err := wav.DecodeWithSampleRate(44100, audioBuffer)
// 	if err != nil {
// 		return err
// 	}

// 	g.audioPlayer, err = g.audioContext.NewPlayer(wavStream)
// 	if err != nil {
// 		return err
// 	}
// 	log.Println("Audio player created successfully")

// 	return nil
// }

// func mousePress(g *Game) error {
// 	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
// 	if mousePressed && !g.mousePressedLastFrame {
// 		x, y := ebiten.CursorPosition()

// 		log.Printf("Mouse clicked at (%d, %d)", x, y)
// 		// Initialize audio on first click
// 		if err := g.initializeAudio(); err != nil {
// 			log.Printf("Audio initialization failed: %v", err)
// 			return nil // Continue running even if audio fails
// 		}

// 		// Play the audio
// 		if g.audioPlayer != nil {
// 			g.audioPlayer.Rewind()
// 			g.audioPlayer.Play()
// 			log.Println("Audio played")
// 		}
// 	}
// 	g.mousePressedLastFrame = mousePressed

// 	return nil
// }
