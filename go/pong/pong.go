package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
	paddleWidth  = 10
	paddleHeight = 100
	ballSize     = 10
)

type PongGame struct {
	leftPaddle  float64
	rightPaddle float64
	ballX       float64
	ballY       float64
	ballVelX    float64
	ballVelY    float64
}

func (g *PongGame) Update() error {
	// Paddle movement
	if ebiten.IsKeyPressed(ebiten.KeyW) && g.leftPaddle > 0 {
		g.leftPaddle -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && g.leftPaddle < float64(screenHeight)-float64(paddleHeight) {
		g.leftPaddle += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) && g.rightPaddle > 0 {
		g.rightPaddle -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && g.rightPaddle < float64(screenHeight)-float64(paddleHeight) {
		g.rightPaddle += 5
	}

	// Ball movement
	g.ballX += g.ballVelX
	g.ballY += g.ballVelY

	// Ball collision with top and bottom
	if g.ballY <= 0 || g.ballY >= float64(screenHeight)-ballSize {
		g.ballVelY = -g.ballVelY
	}

	// Ball collision with paddles
	if (g.ballX <= float64(paddleWidth) && g.ballY+ballSize >= g.leftPaddle && g.ballY <= g.leftPaddle+float64(paddleHeight)) ||
		(g.ballX >= float64(screenWidth)-float64(paddleWidth)-ballSize && g.ballY+ballSize >= g.rightPaddle && g.ballY <= g.rightPaddle+float64(paddleHeight)) {
		g.ballVelX = -g.ballVelX
	}

	// Reset ball if it goes out of bounds
	if g.ballX < 0 || g.ballX > float64(screenWidth) {
		g.ballX = float64(screenWidth) / 2
		g.ballY = float64(screenHeight) / 2
		g.ballVelX = 5
		g.ballVelY = 5
	}

	return nil
}

func (g *PongGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	// Draw paddles
	drawRect(screen, 0, int(g.leftPaddle), paddleWidth, paddleHeight, color.White)
	drawRect(screen, screenWidth-paddleWidth, int(g.rightPaddle), paddleWidth, paddleHeight, color.White)

	// Draw ball
	drawRect(screen, int(g.ballX), int(g.ballY), ballSize, ballSize, color.White)
}

func drawRect(screen *ebiten.Image, x, y, width, height int, col color.Color) {
	rect := ebiten.NewImage(width, height)
	rect.Fill(col)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(rect, op)
}
