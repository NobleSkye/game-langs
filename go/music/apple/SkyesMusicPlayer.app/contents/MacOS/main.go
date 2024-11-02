// main.go

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/ncruces/zenity"
)

// Constants
const (
	screenWidth  = 640
	screenHeight = 480
	sampleRate   = 48000
	volumeStep   = 0.1
	buttonWidth  = 150
	buttonHeight = 30
)

// Colors
var (
	backgroundColor  = color.RGBA{0x1A, 0x1A, 0x1A, 0xFF}
	highlightColor   = color.RGBA{0x00, 0x88, 0xFF, 0xFF}
	volumeBarColor   = color.RGBA{200, 200, 200, 255}
	playerBarColor   = color.RGBA{200, 200, 200, 255}
	progressColor    = color.RGBA{100, 200, 200, 255}
	playButtonColor  = color.RGBA{0, 0, 255, 255}
	pauseButtonColor = color.RGBA{255, 0, 0, 255}
	textColor        = color.RGBA{255, 255, 255, 255}
	changeDirColor   = color.RGBA{0, 128, 0, 255}
)

// Load your custom font
var myFont font.Face

func loadFont() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting executable path: %w", err)
	}

	bundlePath := filepath.Dir(filepath.Dir(filepath.Dir(execPath)))
	fontPath := filepath.Join(bundlePath, "Resources", "font", "Poppins-SemiBold.ttf")

	f, err := os.Open(fontPath)
	if err != nil {
		return fmt.Errorf("error opening font file: %w", err)
	}
	defer f.Close()

	fontBytes, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("error reading font file: %w", err)
	}

	fnt, err := opentype.Parse(fontBytes)
	if err != nil {
		return fmt.Errorf("error parsing font: %w", err)
	}

	myFont, err = opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    12, // Adjust size as needed
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("error creating font face: %w", err)
	}

	return nil
}

// Track information
type Track struct {
	name     string
	data     []byte
	duration time.Duration
}

// Player struct
type Player struct {
	audioContext     *audio.Context
	audioPlayer      *audio.Player
	currentTrack     int
	tracks           []Track
	volume           float64
	volumeFeedback   string
	currentDirectory string
}

// Game struct
type Game struct {
	player *Player
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

// NewPlayer initializes a Player
func NewPlayer(audioContext *audio.Context, directory string) (*Player, error) {
	tracks, err := initTracks(directory)
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(directory)
	if err != nil {
		absPath = directory
	}

	p := &Player{
		audioContext:     audioContext,
		currentTrack:     0,
		tracks:           tracks,
		volume:           1.0,
		currentDirectory: absPath,
	}

	if len(tracks) > 0 {
		if err := p.loadTrackData(); err != nil {
			return nil, err
		}
	}

	return p, nil
}

// Reload tracks from a new directory
func (p *Player) reloadTracks(directory string) error {
	tracks, err := initTracks(directory)
	if err != nil {
		return err
	}
	p.tracks = tracks
	p.currentTrack = 0
	p.currentDirectory = directory
	if len(tracks) > 0 {
		return p.loadTrackData()
	}
	return nil
}

// Toggle play/pause state
func (p *Player) togglePlayPause() {
	if p.audioPlayer.IsPlaying() {
		p.audioPlayer.Pause()
	} else {
		p.audioPlayer.Play()
	}
}

// Update player state with playlist navigation
func (p *Player) update() error {
	// Volume controls
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		p.volume = min(p.volume+volumeStep, 1.0)
		if p.audioPlayer != nil {
			p.audioPlayer.SetVolume(p.volume)
		}
		p.volumeFeedback = fmt.Sprintf("Volume: %.0f%%", p.volume*100)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		p.volume = max(p.volume-volumeStep, 0.0)
		if p.audioPlayer != nil {
			p.audioPlayer.SetVolume(p.volume)
		}
		p.volumeFeedback = fmt.Sprintf("Volume: %.0f%%", p.volume*100)
	}

	// Track navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && len(p.tracks) > 0 {
		p.currentTrack = (p.currentTrack + 1) % len(p.tracks)
		if err := p.loadTrackData(); err != nil {
			return err
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && len(p.tracks) > 0 {
		p.currentTrack = (p.currentTrack - 1 + len(p.tracks)) % len(p.tracks)
		if err := p.loadTrackData(); err != nil {
			return err
		}
	}

	// Playback controls via space
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.togglePlayPause()
	}

	// Handle mouse clicks for changing directory
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()

		// Button bounds for directory change
		buttonX := 500
		buttonY := 450
		if buttonX <= mouseX && mouseX <= buttonX+buttonWidth &&
			buttonY <= mouseY && mouseY <= buttonY+buttonHeight {
			newDirectory, err := zenity.SelectFile(
				zenity.Title("Choose Music Directory"),
				zenity.Directory(),
			)
			if err == nil && newDirectory != "" {
				if err := p.reloadTracks(newDirectory); err != nil {
					fmt.Println("Error loading tracks:", err)
				}
			}
		}

		// Play button bounds
		playButtonX := 20
		playButtonY := 400
		if playButtonX <= mouseX && mouseX <= playButtonX+40 &&
			playButtonY <= mouseY && mouseY <= playButtonY+30 {
			p.togglePlayPause()
		}

		// Pause button bounds
		pauseButtonX := playButtonX + 50
		if pauseButtonX <= mouseX && mouseX <= pauseButtonX+40 &&
			playButtonY <= mouseY && mouseY <= playButtonY+30 {
			p.togglePlayPause()
		}
	}

	// Auto-advance
	if p.audioPlayer != nil && !p.audioPlayer.IsPlaying() &&
		len(p.tracks) > 0 &&
		p.audioPlayer.Current() >= p.tracks[p.currentTrack].duration {
		p.currentTrack = (p.currentTrack + 1) % len(p.tracks)
		return p.loadTrackData()
	}

	// Clear volume feedback after a short time
	if p.volumeFeedback != "" {
		go func() {
			time.Sleep(2 * time.Second)
			p.volumeFeedback = ""
		}()
	}

	return nil
}

// Draw the UI with playlist and buttons
func (p *Player) draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	face := myFont

	// Draw current directory at the top
	text.Draw(screen, "Current Directory: "+p.currentDirectory, face, 10, 15, textColor)

	if len(p.tracks) > 0 {
		text.Draw(screen, p.tracks[p.currentTrack].name, face, 20, 35, textColor)

		// Draw progress bar
		currentTime := p.audioPlayer.Current()
		progress := float64(currentTime) / float64(p.tracks[p.currentTrack].duration)
		if progress > 1 {
			progress = 1
		}

		x, y := 10, 50
		w, h := screenWidth-20, 10
		draw.Draw(screen, image.Rect(x, y, x+w, y+h), &image.Uniform{C: playerBarColor}, image.Point{}, draw.Src)
		progressWidth := int(float64(w) * progress)
		draw.Draw(screen, image.Rect(x, y, x+progressWidth, y+h), &image.Uniform{C: progressColor}, image.Point{}, draw.Src)

		// Draw times
		currentTimeStr := formatDuration(currentTime)
		totalTimeStr := formatDuration(p.tracks[p.currentTrack].duration)
		text.Draw(screen, currentTimeStr, face, x, y+h+15, textColor)
		text.Draw(screen, totalTimeStr, face, x+w-50, y+h+15, textColor)
	} else {
		text.Draw(screen, "No tracks loaded", face, 20, 35, textColor)
	}

	// Draw track list
	for i, track := range p.tracks {
		var color color.Color
		if i == p.currentTrack {
			color = highlightColor
		} else {
			color = textColor
		}
		xPos := 100
		if (i+1)%2 == 0 {
			xPos = 350
		}
		yPos := 100 + (i/2)*20
		text.Draw(screen, track.name, face, xPos, yPos, color)
	}

	// Draw volume bar
	volumeX, volumeY := 10, 400
	volumeW, volumeH := 10, screenHeight-420
	draw.Draw(screen, image.Rect(volumeX, volumeY, volumeX+volumeW, volumeY+volumeH), &image.Uniform{C: volumeBarColor}, image.Point{}, draw.Src)
	volumeProgress := int(float64(volumeH) * (1 - p.volume))
	draw.Draw(screen, image.Rect(volumeX, volumeY+volumeProgress, volumeX+volumeW, volumeY+volumeH), &image.Uniform{C: progressColor}, image.Point{}, draw.Src)

	// Draw volume text
	volumeStr := fmt.Sprintf("Volume: %.0f%%", p.volume*100)
	text.Draw(screen, volumeStr, face, volumeX+15, volumeY+volumeH-15, textColor)

	// Draw play button
	playButtonY := 400
	playButtonX := volumeX + volumeW + 10
	draw.Draw(screen, image.Rect(playButtonX, playButtonY, playButtonX+40, playButtonY+30), &image.Uniform{C: playButtonColor}, image.Point{}, draw.Src)
	text.Draw(screen, ">", face, playButtonX+5, playButtonY+20, textColor)

	// Draw pause button
	pauseButtonX := playButtonX + 50
	draw.Draw(screen, image.Rect(pauseButtonX, playButtonY, pauseButtonX+40, playButtonY+30), &image.Uniform{C: pauseButtonColor}, image.Point{}, draw.Src)
	text.Draw(screen, "||", face, pauseButtonX+5, playButtonY+20, textColor)

	// Draw directory change button
	buttonX := 500
	buttonY := 450
	draw.Draw(screen, image.Rect(buttonX, buttonY, buttonX+buttonWidth, buttonY+buttonHeight), &image.Uniform{C: changeDirColor}, image.Point{}, draw.Src)
	text.Draw(screen, "Change Directory", face, buttonX+5, buttonY+20, textColor)

	// Draw volume feedback if available
	if p.volumeFeedback != "" {
		text.Draw(screen, p.volumeFeedback, face, 20, 370, textColor)
	}
}

// Format duration to string
func formatDuration(duration time.Duration) string {
	minutes := int(duration.Seconds()) / 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// Utility function to get minimum value
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Utility function to get maximum value
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Update method for Game
func (g *Game) Update() error {
	return g.player.update()
}

// Draw method for Game
func (g *Game) Draw(screen *ebiten.Image) {
	g.player.draw(screen)
}

// Layout method for Game
func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Load the custom font
	if err := loadFont(); err != nil {
		fmt.Println("Error loading font:", err)
		return
	}

	audioContext := audio.NewContext(sampleRate)
	player, err := NewPlayer(audioContext, ".")
	if err != nil {
		fmt.Println("Error initializing player:", err)
		return
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Skye's Music Player")

	if err := ebiten.RunGame(&Game{player}); err != nil {
		fmt.Println("Error running Player:", err)
	}
}
