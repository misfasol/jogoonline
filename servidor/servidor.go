package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"proto"
)

type Conexao struct {
	Endereco *net.UDPAddr
	Jogador  proto.Jogador
}

func tratarCliente(dados []byte, conn *net.UDPConn, addr *net.UDPAddr, conexoes *map[string]proto.Jogador) {
	pacote := proto.PacoteBase{}
	buf := bytes.NewBuffer(dados)
	binary.Read(buf, binary.BigEndian, &pacote)

	// fmt.Printf("Novo pacote inicial: %v\n", pacote)

	switch pacote.Tipo {

	case proto.PACOTE_TIPO_INICIAL:
		pacoteInicial := proto.PacoteInicial{
			Base:    pacote,
			Jogador: proto.Jogador{},
		}
		binary.Read(buf, binary.BigEndian, &pacoteInicial.Jogador)
		(*conexoes)[fmt.Sprintf("%v:%v", (*addr).IP, (*addr).Port)] = pacoteInicial.Jogador
		// fmt.Println(conexoes)
		fmt.Printf("Nova conexao: %v:%v\n", addr.IP, addr.Port)
		// fmt.Printf("Novo pacote: %+v\n", pacoteInicial)
		novoBuffer := new(bytes.Buffer)
		binary.Write(novoBuffer, binary.BigEndian, uint8(0x01))
		conn.WriteToUDP(novoBuffer.Bytes(), addr)

	case proto.PACOTE_TIPO_ATUALIZA:
		pacoteAtualiza := proto.PacoteAtualiza{
			Base:          pacote,
			NumeroDeSerie: 0,
			Jogador:       proto.Jogador{},
			Input:         proto.Input{},
		}
		binary.Read(buf, binary.BigEndian, &pacoteAtualiza.NumeroDeSerie)
		binary.Read(buf, binary.BigEndian, &pacoteAtualiza.Jogador)
		binary.Read(buf, binary.BigEndian, &pacoteAtualiza.Input)
		// fmt.Printf("Novo pacote: %+v\n", pacoteAtualiza)
		if pacoteAtualiza.Input.Up {
			pacoteAtualiza.Jogador.PosY -= pacoteAtualiza.Jogador.Velocidade
		}
		if pacoteAtualiza.Input.Down {
			pacoteAtualiza.Jogador.PosY += pacoteAtualiza.Jogador.Velocidade
		}
		if pacoteAtualiza.Input.Left {
			pacoteAtualiza.Jogador.PosX -= pacoteAtualiza.Jogador.Velocidade
		}
		if pacoteAtualiza.Input.Right {
			pacoteAtualiza.Jogador.PosX += pacoteAtualiza.Jogador.Velocidade
		}
		fmt.Println(addr)
		(*conexoes)[fmt.Sprintf("%v:%v", (*addr).IP, (*addr).Port)] = pacoteAtualiza.Jogador
		novoBuffer := new(bytes.Buffer)
		binary.Write(novoBuffer, binary.BigEndian, pacoteAtualiza)
		conn.WriteToUDP(novoBuffer.Bytes(), addr)

	}
}

func main() {
	upADDR, err := net.ResolveUDPAddr("udp", "0.0.0.0:9000")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", upADDR)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Ouvindo")
	conexoes := make(map[string]proto.Jogador)

	for {
		buf := make([]byte, proto.TAMANHO_MAXIMO_PACOTE)
		_, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		go tratarCliente(buf, conn, addr, &conexoes)
		fmt.Println(conexoes)
	}
}
