package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	NFICHAS = 4
)

type GameData struct {
	NumPlayers int
	GameMap    [40]int
	NumTurno   int
}

type Ficha struct {
	id       int
	color    string
	posicion int
	estado   int
	meta     bool
}

type Lanzamiento struct {
	dadoA   int
	dadoB   int
	avanzar bool
}

var direccionRemota string
var fichas []Ficha
var mapa [40]int

func intArrayToString() string {
	strArr := make([]string, len(mapa))

	for i, v := range mapa {
		if v == -1 {
			strArr[i] = "X"
		} else {
			strArr[i] = strconv.Itoa(i)
		}
	}

	return "<p class='registro'>[" + fmt.Sprintf("%v", strArr) + "]</p>\n"
}

func guardarPosicionesEnArchivo(color string, turno int, jugo_jugador int) {
	ArchivoRegistro := fmt.Sprintf("archivo_%s.txt", color)
	file, err := os.OpenFile(ArchivoRegistro, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	if turno != 0 {
		messageT := fmt.Sprintf("<p class='registro'> TurnoActual: %d </p>\n", turno)
		_, err = file.WriteString(messageT)
		if err != nil {
			fmt.Println("Error al escribir en el archivo:", err)
			return
		}
	} else {
		_, err = file.WriteString(intArrayToString())
		if err != nil {
			fmt.Println("Error al escribir en el archivo:", err)
			return
		}
	}

	if jugo_jugador == -1 {

	} else if jugo_jugador == 1 {
		for _, f := range fichas {
			message := fmt.Sprintf("<p class='registro'> Id: %d - Color: %s - Posición: %d - Meta:%t </p>\n", f.id, f.color, f.posicion, f.meta)
			_, err := file.WriteString(message)
			if err != nil {
				fmt.Println("Error al escribir en el archivo:", err)
				return
			}
		}
	} else if jugo_jugador == 0 {
		_, err = file.WriteString("<p class='registro' style='color:black'> ESTE JUGADOR PERDIO SU TURNO </p>\n")
		if err != nil {
			fmt.Println("Error al escribir en el archivo:", err)
			return
		}
	} else {
		_, err = file.WriteString("<p class='registro' style='color:yellow'> FELICITACIONES GANASTE EL JUEGO</p>\n")
		if err != nil {
			fmt.Println("Error al escribir en el archivo:", err)
			return
		}
	}
	_, err = file.WriteString("--------------------------------------------------\n")
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		return
	}
}

func enviar(gameData GameData) {
	con, _ := net.Dial("tcp", direccionRemota)
	jsonBytes, _ := json.Marshal(gameData)
	jsonStr := string(jsonBytes)
	defer con.Close()
	fmt.Fprintln(con, jsonStr)
}

func lanzarDados() Lanzamiento {
	valor := rand.Intn(2)
	tiro := Lanzamiento{
		dadoA:   rand.Intn(6) + 1,
		dadoB:   rand.Intn(6) + 1,
		avanzar: valor == 1,
	}
	return tiro
}

func initialize_player(color string) {
	for j := 0; j < NFICHAS; j++ {
		ficha := Ficha{
			id:       j + 1,
			color:    color,
			posicion: 0,
			meta:     false,
		}
		fichas = append(fichas, ficha)
	}
}

func manejador(con net.Conn, color string, chFichas []chan bool) {
	var gameData GameData
	// time.Sleep(1 * time.Second)
	defer con.Close()
	fmt.Printf("Turno del Jugador: %s\n", color)
	defer con.Close()
	br := bufio.NewReader(con)
	msg, _ := br.ReadString('\n')
	msg = strings.TrimSpace(msg)
	json.Unmarshal([]byte(msg), &gameData)
	if gameData.NumPlayers > 0 {
		fmt.Println("Inicializando Fichas")
		initialize_player(color)
		gameData.NumPlayers = gameData.NumPlayers - 1
		mapa = gameData.GameMap
		fmt.Println(mapa)
		fmt.Println("------------------------")
		guardarPosicionesEnArchivo(color, gameData.NumTurno, -1)
		enviar(gameData)
	} else {
		gameData.NumTurno = gameData.NumTurno + 1
		fichasCompletadas := 0
		for _, f := range fichas {
			if f.meta == true {
				fichasCompletadas++
			}
		}
		if fichasCompletadas < 4 {
			jugo_jugador := turno_jugador(chFichas[0], chFichas[1], chFichas[2], chFichas[3])
			enviar(gameData)
			guardarPosicionesEnArchivo(color, gameData.NumTurno, jugo_jugador)
		} else {
			guardarPosicionesEnArchivo(color, gameData.NumTurno, 2)
			fmt.Printf("El jugador %s ha ganado el juego\n", color)
			fmt.Println(fichas)
		}
	}
}

func pierdeTurno() bool {
	for i := 0; i < 4; i++ {
		for ind, valor := range mapa {
			if valor == -1 && ind == fichas[i].posicion {
				fichas[i].estado += 1
				if fichas[i].estado > 2 {
					fichas[i].estado = 2
				}
			}
			if fichas[i].estado == 2 && valor == 0 && ind == fichas[i].posicion {
				fichas[i].estado = 0
			}
		}
	}
	for i := 0; i < 4; i++ {
		if fichas[i].estado == 1 {
			return true
		}
	}
	return false
}

func turno_jugador(ficha1 chan bool, ficha2 chan bool, ficha3 chan bool, ficha4 chan bool) int {
	var tiro Lanzamiento = lanzarDados()
	var ind int
	if !pierdeTurno() {
		go func() {
			if fichas[0].meta == false {
				ficha1 <- true
			}
		}()
		go func() {
			if fichas[1].meta == false {
				ficha2 <- true
			}
		}()
		go func() {
			if fichas[2].meta == false {
				ficha3 <- true
			}
		}()
		go func() {
			if fichas[3].meta == false {
				ficha4 <- true
			}
		}()

		select {
		case <-ficha1:
			fmt.Printf("(JUEGA FICHA 1)\n")
			ind = 0

		case <-ficha2:
			fmt.Printf("(JUEGA FICHA 2)\n")
			ind = 1

		case <-ficha3:
			fmt.Printf("(JUEGA FICHA 3)\n")
			ind = 2

		case <-ficha4:
			fmt.Printf("(JUEGA FICHA 4)\n")
			ind = 3
		}

		go func() {
			for {
				select {
				case <-ficha1:
				case <-ficha2:
				case <-ficha3:
				case <-ficha4:
					// Descartar elementos del canal
				default:
					// El canal está vacío
					return
				}
			}
		}()

		if tiro.avanzar {
			fmt.Println("RESULTADO LANZAMIENTO: ", tiro.dadoA+tiro.dadoB)
			fichas[ind].posicion += tiro.dadoA + tiro.dadoB
			if fichas[ind].posicion > 39 {
				fichas[ind].posicion = 39 - (fichas[ind].posicion - 39)
			}
		} else {
			fmt.Println("RESULTADO LANZAMIENTO: ", tiro.dadoA-tiro.dadoB)
			fichas[ind].posicion += tiro.dadoA - tiro.dadoB
			if fichas[ind].posicion < 0 {
				fichas[ind].posicion = 0
			}
		}
		fmt.Println("POSCION ACTUAL DE LA FICHA: ", fichas[ind].posicion)
		// gano?
		for i := 0; i < 4; i++ {
			if fichas[i].posicion == 39 {
				fichas[i].meta = true
			}
		}
		fmt.Println("-------------------------")
		return 1

	} else {
		fmt.Println("ESTE JUGADOR PERDIO SU TURNO")
		fmt.Println("-------------------------")
		return 0
	}
}

func main() {

	br := bufio.NewReader(os.Stdin)
	fmt.Print("Ingresa el color del jugador: ")
	color, _ := br.ReadString('\n')
	color = strings.TrimSpace(color)

	fmt.Print("Puerto Actual: ")
	strPuertoLocal, _ := br.ReadString('\n')
	strPuertoLocal = strings.TrimSpace(strPuertoLocal)
	direccionLocal := fmt.Sprintf("localhost:%s", strPuertoLocal)

	fmt.Print("Puerto Destino: ")
	strPuertoRemoto, _ := br.ReadString('\n')
	strPuertoRemoto = strings.TrimSpace(strPuertoRemoto)
	direccionRemota = fmt.Sprintf("localhost:%s", strPuertoRemoto)

	chFichas := make([]chan bool, NFICHAS)

	for i := range chFichas {
		chFichas[i] = make(chan bool)
	}

	ln, _ := net.Listen("tcp", direccionLocal)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go manejador(con, color, chFichas)
	}
}
