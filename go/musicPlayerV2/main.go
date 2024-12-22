package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

const (
	defaultMusicPath = "./mp3"
)

type Player struct {
	playlist    []string
	currentSong int
	streamer    beep.StreamSeeker
	ctrl        *beep.Ctrl
	format      beep.Format
	window      *ui.Window
	songList    *ui.MultilineEntry
	playing     bool
}

func newPlayer() *Player {
	return &Player{
		currentSong: -1,
		playing:     false,
	}
}

func (p *Player) loadDirectory(dir string) error {
	p.playlist = nil
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".mp3" {
			p.playlist = append(p.playlist, path)
		}
		return nil
	})

	// Update song list display
	p.songList.SetText("")
	for i, song := range p.playlist {
		p.songList.Append(fmt.Sprintf("%d. %s\n", i+1, filepath.Base(song)))
	}
	return err
}

func (p *Player) playSong(index int) error {
	if index < 0 || index >= len(p.playlist) {
		return fmt.Errorf("invalid song index")
	}

	if p.streamer != nil {
		speaker.Clear()
	}

	f, err := os.Open(p.playlist[index])
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	p.ctrl = &beep.Ctrl{Streamer: streamer}
	speaker.Play(p.ctrl)

	p.currentSong = index
	p.streamer = streamer
	p.format = format
	p.playing = true

	return nil
}

func main() {
	err := ui.Main(func() {
		player := newPlayer()

		window := ui.NewWindow("Music Player", 400, 300, true)
		window.SetMargined(true)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})

		vbox := ui.NewVerticalBox()
		vbox.SetPadded(true)

		hbox := ui.NewHorizontalBox()
		hbox.SetPadded(true)

		selectButton := ui.NewButton("Select Directory")
		selectButton.OnClicked(func(*ui.Button) {
			entry := ui.NewEntry()
			entry.SetText("Enter directory path:")
			dialog := ui.NewWindow("Select Directory", 300, 100, false)
			dialog.SetChild(entry)
			dialog.Show()

			entry.OnChanged(func(entry *ui.Entry) {
				path := entry.Text()
				if _, err := os.Stat(path); err == nil {
					player.loadDirectory(path)
					dialog.Hide()
				}
			})
		})
		hbox.Append(selectButton, false)

		playButton := ui.NewButton("Play")
		playButton.OnClicked(func(*ui.Button) {
			if len(player.playlist) > 0 {
				if player.currentSong == -1 {
					player.playSong(0)
				} else if !player.playing {
					speaker.Lock()
					player.ctrl.Paused = false
					player.playing = true
					speaker.Unlock()
				}
			}
		})
		hbox.Append(playButton, false)

		pauseButton := ui.NewButton("Pause")
		pauseButton.OnClicked(func(*ui.Button) {
			if player.playing {
				speaker.Lock()
				player.ctrl.Paused = true
				player.playing = false
				speaker.Unlock()
			}
		})
		hbox.Append(pauseButton, false)

		vbox.Append(hbox, false)

		songList := ui.NewMultilineEntry()
		songList.SetReadOnly(true)
		player.songList = songList
		vbox.Append(songList, true)

		window.SetChild(vbox)
		window.Show()
	})

	if err != nil {
		panic(err)
	}
}
