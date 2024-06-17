package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/color"
	"net"
	"os"
	"proto"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	LARGURA_JANELA = 1280
	ALTURA_JANELA  = 720
)

var (
	numeroDeSeriePacote uint32 = 0
)

func sequenciaInicianlizacao(conn *net.UDPConn, jogador proto.Jogador) {
	pacote := proto.PacoteInicial{
		Base: proto.PacoteBase{
			Tamanho: 0,
			Versao: 1,
			Tipo: proto.PACOTE_TIPO_INICIAL,
		},
		Jogador: jogador,
	}
	pacote.Base.Tamanho = uint16(binary.Size(pacote))
	binary.Write(conn, binary.BigEndian, pacote)
	fmt.Printf("enviei o pacote %+v\n", pacote)
	dados := make([]byte, 1)
	conn.ReadFromUDP(dados)
	fmt.Printf("Recebi pacote tamanho: %v pacote: %v\n", len(dados), dados)
	if dados[0] == 0x00 {
		fmt.Println("Erro se conectando no servidor")
		os.Exit(1)
	}
}

func mandarPacote(conn *net.UDPConn, jogador proto.Jogador) {
	numeroDeSeriePacote += 1

	pacote := proto.PacoteAtualiza{
		Base: proto.PacoteBase{
			Tamanho: 0,
			Versao: 1,
			Tipo: proto.PACOTE_TIPO_ATUALIZA,
		},
		NumeroDeSerie: numeroDeSeriePacote,
		Jogador:       jogador,
		Input:         proto.NovoInput(),
	}
	pacote.Base.Tamanho = uint16(binary.Size(pacote))

	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, pacote)
	conn.Write(buffer.Bytes())
}

func receberPacote(conn *net.UDPConn, jogador *proto.Jogador) {
	dados := make([]byte, proto.TAMANHO_MAXIMO_PACOTE)
	conn.ReadFromUDP(dados)
	buffer := bytes.NewBuffer(dados)
	pacote := &proto.PacoteAtualiza{}
	binary.Read(buffer, binary.BigEndian, pacote)
	*jogador = pacote.Jogador
}

func main() {
	fmt.Println("oi")

	upAddr, err := net.ResolveUDPAddr("udp", "192.168.0.12:9000")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, upAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	jogador := proto.NovoJogador(10, 10, 50, 50, 5, color.RGBA{255, 10, 255, 255})
	sequenciaInicianlizacao(conn, jogador)

	rl.InitWindow(LARGURA_JANELA, ALTURA_JANELA, "jogo")
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		mandarPacote(conn, jogador)

		receberPacote(conn, &jogador)

		rl.BeginDrawing()

		rl.ClearBackground(rl.Gray)

		jogador.Pintar()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}