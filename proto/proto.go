package proto

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Jogador struct {
	PosX, PosY, Width, Height int32
	Velocidade                int32
	Col                       color.RGBA
}

func NovoJogador(posX, posY, width, height, velocidade int32, col color.RGBA) Jogador {
	return Jogador{PosX: posX, PosY: posY, Width: width, Height: height, Velocidade: velocidade, Col: col}
}

func (j *Jogador) Pintar() {
	rl.DrawRectangle(j.PosX, j.PosY, j.Width, j.Height, j.Col)
}

func (j *Jogador) AtualizarPosicao(LARGURA_JANELA, ALTURA_JANELA int32) {
	if rl.IsKeyDown(rl.KeyLeft) {
		j.PosX -= j.Velocidade
	}
	if rl.IsKeyDown(rl.KeyRight) {
		j.PosX += j.Velocidade
	}
	if rl.IsKeyDown(rl.KeyUp) {
		j.PosY -= j.Velocidade
	}
	if rl.IsKeyDown(rl.KeyDown) {
		j.PosY += j.Velocidade
	}
	if j.PosX < 0 {
		j.PosX = 0
	} else if j.PosX+j.Width > LARGURA_JANELA {
		j.PosX = LARGURA_JANELA - j.Width
	}
	if j.PosY < 0 {
		j.PosY = 0
	} else if j.PosY+j.Height > ALTURA_JANELA {
		j.PosY = ALTURA_JANELA - j.Height
	}
}

type Input struct {
	Up, Down, Left, Right bool
}

func NovoInput() Input {
	return Input{
		Up:    rl.IsKeyDown(rl.KeyUp),
		Down:  rl.IsKeyDown(rl.KeyDown),
		Left:  rl.IsKeyDown(rl.KeyLeft),
		Right: rl.IsKeyDown(rl.KeyRight),
	}
}

const (
	PACOTE_TIPO_INICIAL  = 0x01
	PACOTE_TIPO_ATUALIZA = 0x02
)

type PacoteBase struct {
	Tamanho uint16
	Versao  uint8
	Tipo    uint8
}

type PacoteInicial struct {
	Base    PacoteBase
	Jogador Jogador
}

type PacoteAtualiza struct {
	Base          PacoteBase
	NumeroDeSerie uint32
	Jogador       Jogador
	Input         Input
}

const TAMANHO_MAXIMO_PACOTE = 508
