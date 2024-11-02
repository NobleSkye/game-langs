package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"
	"time"

	_ "embed"

	"golang.org/x/image/font/basicfont"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// Embed MP3 audio files
var trackFiles []string // Hold track filenames

// Constants
const (
	screenWidth  = 640
	screenHeight = 480
	sampleRate   = 48000
	volumeStep   = 0.1
)

// Colors
var (
	backgroundColor = color.RGBA{0x1A, 0x1A, 0x1A, 0xFF}
	playerBarColor  = color.RGBA{0x40, 0x40, 0x40, 0xFF}
	progressColor   = color.RGBA{0x00, 0x88, 0xFF, 0xFF}
	textColor       = color.White
	highlightColor  = color.RGBA{0x00, 0x88, 0xFF, 0xFF}
)

// Track information
type Track struct {
	name     string
	data     []byte
	duration time.Duration
}

// Player struct with enhanced features
type Player struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	currentTrack int
	tracks       []Track
	volume       float64
}

// Load track data
func (p *Player) loadTrackData() error {
	track := p.tracks[p.currentTrack]

	reader, err := mp3.DecodeF32(bytes.NewReader(track.data))
	if err != nil {
		return err
	}

	player, err := p.audioContext.NewPlayerF32(reader)
	if err != nil {
		return err
	}

	if p.audioPlayer != nil {
		p.audioPlayer.Close()
	}

	p.audioPlayer = player
	p.tracks[p.currentTrack].duration = time.Duration(reader.Length()) * time.Second / 8 / sampleRate
	p.audioPlayer.SetVolume(p.volume)
	p.audioPlayer.Play()

	return nil
}

// Initialize tracks list by scanning a directory
func initTracks(directory string) ([]Track, error) {
	var tracks []Track

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".mp3" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			tracks = append(tracks, Track{name: info.Name(), data: data, duration: 0})
		}
		return nil
	})
	return tracks, err
}

// NewPlayer initializes a Player with audio context and directory of tracks
func NewPlayer(audioContext *audio.Context, directory string) (*Player, error) {
	tracks, err := initTracks(directory) // Load tracks from directory
	if err != nil {
		return nil, err // Handle error
	}

	p := &Player{
		audioContext: audioContext,
		currentTrack: 0,
		tracks:       tracks,
		volume:       1.0,
	}

	if len(tracks) > 0 {
		if err := p.loadTrackData(); err != nil { // Load first track data
			return nil, err
		}
	}

	return p, nil // Return initialized player
}

// Update player state with playlist navigation
func (p *Player) update() error {
	// Volume controls
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		p.volume = min(p.volume+volumeStep, 1.0)
		p.audioPlayer.SetVolume(p.volume)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		p.volume = max(p.volume-volumeStep, 0.0)
		p.audioPlayer.SetVolume(p.volume)
	}

	// Track navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		p.currentTrack = (p.currentTrack + 1) % len(p.tracks)
		return p.loadTrackData()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		p.currentTrack = (p.currentTrack - 1 + len(p.tracks)) % len(p.tracks)
		return p.loadTrackData()
	}

	// Playback controls
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if p.audioPlayer.IsPlaying() {
			p.audioPlayer.Pause()
		} else {
			p.audioPlayer.Play()
		}
	}

	// Auto-advance
	if !p.audioPlayer.IsPlaying() && p.audioPlayer.Current() >= p.tracks[p.currentTrack].duration {
		p.currentTrack = (p.currentTrack + 1) % len(p.tracks)
		return p.loadTrackData()
	}

	return nil
}

// Draw the UI with playlist and buttons
func (p *Player) draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	face := basicfont.Face7x13
	text.Draw(screen, p.tracks[p.currentTrack].name, face, 20, 35, textColor)

	// Draw progress bar
	currentTime := p.audioPlayer.Current()
	progress := float64(currentTime) / float64(p.tracks[p.currentTrack].duration)
	if progress > 1 {
		progress = 1
	}

	x, y := 10, 50
	w, h := screenWidth-20, 10
	draw.Draw(screen, image.Rect(x, y, x+w, y+h), &image.Uniform{playerBarColor}, image.Point{}, draw.Src)
	progressWidth := int(float64(w) * progress)
	draw.Draw(screen, image.Rect(x, y, x+progressWidth, y+h), &image.Uniform{progressColor}, image.Point{}, draw.Src)

	// Draw times
	currentTimeStr := formatDuration(currentTime)
	totalTimeStr := formatDuration(p.tracks[p.currentTrack].duration)
	text.Draw(screen, currentTimeStr, face, x, y+h+15, textColor)
	text.Draw(screen, totalTimeStr, face, x+w-50, y+h+15, textColor)

	// Draw track list in two columns based on track number
	for i, track := range p.tracks {
		var color color.Color
		if i == p.currentTrack {
			color = highlightColor
		} else {
			color = textColor
		}

		// Determine x position based on odd/even track index
		var xPos int
		if (i+1)%2 == 0 { // Even track (1-based index)
			xPos = 220 // Right side
		} else { // Odd track (1-based index)
			xPos = 20 // Left side
		}

		yPos := 100 + (i/2)*20 // New row every two tracks
		text.Draw(screen, track.name, face, xPos, yPos, color)
	}
}

// Format duration as MM:SS
func formatDuration(d time.Duration) string {
	seconds := int(d.Seconds())
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// Game struct to encapsulate player
type Game struct {
	player *Player
}

// Update game state
func (g *Game) Update() error {
	return g.player.update()
}

// Draw game frame
func (g *Game) Draw(screen *ebiten.Image) {
	g.player.draw(screen)
}

// Set the game layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Main function to initialize and run the game
func main() {
	audioContext := audio.NewContext(sampleRate)

	// Specify the directory where your MP3 files are located
	directory := "./mp3" // Change this to your directory
	player, err := NewPlayer(audioContext, directory)
	if err != nil {
		panic(err)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Skye's Music Player")
	if err := ebiten.RunGame(&Game{player: player}); err != nil {
		panic(err)
	}
}

// Helper functions for min/max
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
