package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
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
var (
	//go:embed 01 Key.mp3
	track01 []byte
	//go:embed 02 Door.mp3
	track02 []byte
	//go:embed 03 Subwoofer Lullaby.mp3
	track03 []byte
	//go:embed 04 Death.mp3
	track04 []byte
	//go:embed 05 Living Mice.mp3
	track05 []byte
	//go:embed 06 Moog City.mp3
	track06 []byte
	//go:embed 07 Haggstrom.mp3
	track07 []byte
	//go:embed 08 Minecraft.mp3
	track08 []byte
	//go:embed 09 Oxygene.mp3
	track09 []byte
	//go:embed 10 Equinoxe.mp3
	track10 []byte
	//go:embed 11 Miceon Venus.mp3
	track11 []byte
	//go:embed 12 Dry Hands.mp3
	track12 []byte
	//go:embed 13 Wet Hands.mp3
	track13 []byte
	//go:embed 14 Clark.mp3
	track14 []byte
	//go:embed 15 Chris.mp3
	track15 []byte
	//go:embed 16 Thirteen.mp3
	track16 []byte
	//go:embed 17 Excuse.mp3
	track17 []byte
	//go:embed 18 Sweden.mp3
	track18 []byte
	//go:embed 19 Cat.mp3
	track19 []byte
	//go:embed 20 Dog.mp3
	track20 []byte
	//go:embed 21 Danny.mp3
	track21 []byte
	//go:embed 22 Beginning.mp3
	track22 []byte
	//go:embed 23 Droopy Likes Ricochet.mp3
	track23 []byte
	//go:embed 24 Droopy Likes Your Face.mp3
	track24 []byte
)

// Track information
type Track struct {
	name     string
	data     []byte
	duration time.Duration
}

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

// Player struct with enhanced features
type Player struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	currentTrack int
	tracks       []Track
	volume       float64
}

func initTracks() []Track {
	return []Track{
		{"01 Key.mp3", track01, 0},
		{"02 Door.mp3", track02, 0},
		{"03 Subwoofer Lullaby.mp3", track03, 0},
		{"04 Death.mp3", track04, 0},
		{"05 Living Mice.mp3", track05, 0},
		{"06 Moog City.mp3", track06, 0},
		{"07 Haggstrom.mp3", track07, 0},
		{"08 Minecraft.mp3", track08, 0},
		{"09 Oxygene.mp3", track09, 0},
		{"10 Equinoxe.mp3", track10, 0},
		{"11 Miceon Venus.mp3", track11, 0},
		{"12 Dry Hands.mp3", track12, 0},
		{"13 Wet Hands.mp3", track13, 0},
		{"14 Clark.mp3", track14, 0},
		{"15 Chris.mp3", track15, 0},
		{"16 Thirteen.mp3", track16, 0},
		{"17 Excuse.mp3", track17, 0},
		{"18 Sweden.mp3", track18, 0},
		{"19 Cat.mp3", track19, 0},
		{"20 Dog.mp3", track20, 0},
		{"21 Danny.mp3", track21, 0},
		{"22 Beginning.mp3", track22, 0},
		{"23 Droopy Likes Ricochet.mp3", track23, 0},
		{"24 Droopy Likes Your Face.mp3", track24, 0},
	}
}

// NewPlayer initializes a Player
func NewPlayer(audioContext *audio.Context) (*Player, error) {
	tracks := initTracks()
	p := &Player{
		audioContext: audioContext,
		currentTrack: 0,
		tracks:       tracks,
		volume:       1.0,
	}

	if err := p.loadTrack(); err != nil {
		return nil, err
	}

	return p, nil
}

// Load track with volume control
func (p *Player) loadTrack() error {
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
		return p.loadTrack()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		p.currentTrack = (p.currentTrack - 1 + len(p.tracks)) % len(p.tracks)
		return p.loadTrack()
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
		return p.loadTrack()
	}

	return nil
}

// Draw the UI with playlist and buttons
func (p *Player) draw(screen *ebiten.Image) {
	// Clear background
	screen.Fill(backgroundColor)

	// Draw current track info
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

	// Draw buttons
	// drawButtons(screen)

	// Draw track list
	for i, track := range p.tracks {
		var color color.Color
		if i == p.currentTrack {
			color = highlightColor
		} else {
			color = textColor
		}
		text.Draw(screen, track.name, face, 20, 100+i*20, color)
	}
}

// Draw buttons on the screen
// func drawButtons(screen *ebiten.Image) {
// 	face := basicfont.Face7x13
// 	text.Draw(screen, "Previous", face, 50, 430, textColor)
// 	text.Draw(screen, "Play/Pause", face, 200, 430, textColor)
// 	text.Draw(screen, "Next", face, 350, 430, textColor)
// }

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
	player, err := NewPlayer(audioContext)
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
