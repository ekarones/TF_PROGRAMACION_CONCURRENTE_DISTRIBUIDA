package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
)

type GameData struct {
	NumPlayers int
	GameMap    [40]int
	NumTurno   int
}

var puertoRemoto string

func startGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	maxObstaculosStr := r.FormValue("maxObstaculos")
	maxObstaculos, err := strconv.Atoi(maxObstaculosStr)
	if err != nil {
		http.Error(w, "Invalid number of obstacles", http.StatusBadRequest)
		return
	}

	color := r.FormValue("opcion")
	if color == "rojo" {
		puertoRemoto = "8000"
	} else if color == "azul" {
		puertoRemoto = "8001"
	} else if color == "verde" {
		puertoRemoto = "8002"
	} else {
		puertoRemoto = "8003"
	}

	fmt.Println("nodo:" + puertoRemoto)
	direccionRemota := fmt.Sprintf("localhost:%s", puertoRemoto)

	var gameMap [40]int
	invalidPositions := []int{0, 39}
	initializeGameMap(&gameMap, invalidPositions, maxObstaculos)

	gameData := GameData{
		NumPlayers: 4,
		GameMap:    gameMap,
		NumTurno:   0,
	}

	jsonBytes, err := json.Marshal(gameData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonStr := string(jsonBytes)
	fmt.Println(jsonStr)
	fmt.Println(direccionRemota)

	// Aquí puedes realizar cualquier lógica adicional antes de enviar los datos al cliente

	con, _ := net.Dial("tcp", direccionRemota)
	defer con.Close()
	fmt.Fprintln(con, jsonStr)

	// Respuesta al cliente
	http.Redirect(w, r, "/show_game", http.StatusSeeOther)
}

func initializeGameMap(tabla *[40]int, invalidPositions []int, maxObstaculos int) {
	var contador int

	for contador < maxObstaculos {
		number := rand.Intn(40)
		found := false
		for _, v := range invalidPositions {
			if number == v {
				found = true
				break
			}
		}
		if !found {
			contador++
			(*tabla)[number] = -1
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func showGame(w http.ResponseWriter, r *http.Request) {
	// Leer el contenido del archivo de texto ROJO

	contentRojo, err := ioutil.ReadFile("archivo_ROJO.txt")
	if err != nil {
		fmt.Println("Error al leer el archivo ROJO:", err)
		http.Error(w, "Error al leer el archivo ROJO", http.StatusInternalServerError)
		return
	}

	// Convertir el contenido a string para ROJO
	fileContentRojo := string(contentRojo)

	// Leer el contenido del archivo de texto AZUL
	contentAzul, err := ioutil.ReadFile("archivo_AZUL.txt")
	if err != nil {
		fmt.Println("Error al leer el archivo AZUL:", err)
		http.Error(w, "Error al leer el archivo AZUL", http.StatusInternalServerError)
		return
	}

	// Convertir el contenido a string para AZUL
	fileContentAzul := string(contentAzul)

	// Leer el contenido del archivo de texto VERDE
	contentVerde, err := ioutil.ReadFile("archivo_VERDE.txt")
	if err != nil {
		fmt.Println("Error al leer el archivo VERDE:", err)
		http.Error(w, "Error al leer el archivo VERDE", http.StatusInternalServerError)
		return
	}

	// Convertir el contenido a string para VERDE
	fileContentVerde := string(contentVerde)

	contentAmarillo, err := ioutil.ReadFile("archivo_AMARILLO.txt")
	if err != nil {
		fmt.Println("Error al leer el archivo AMARILLO:", err)
		http.Error(w, "Error al leer el archivo", http.StatusInternalServerError)
		return
	}
	fileContentAmarillo := string(contentAmarillo)

	// Crear una estructura de datos para pasar al template
	data := struct {
		// FileContentMapa     template.HTML
		FileContentRojo     template.HTML
		FileContentAzul     template.HTML
		FileContentVerde    template.HTML
		FileContentAmarillo template.HTML
	}{
		// FileContentMapa:     template.HTML(fileContentMapa),
		FileContentRojo:     template.HTML(fileContentRojo),
		FileContentAzul:     template.HTML(fileContentAzul),
		FileContentVerde:    template.HTML(fileContentVerde),
		FileContentAmarillo: template.HTML(fileContentAmarillo),
	}

	// Parsear el template y pasar los datos
	tmpl, err := template.ParseFiles("templates/show_game.html")
	if err != nil {
		fmt.Println("Error al parsear el template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Ejecutar el template y manejar posibles errores
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("Error al ejecutar el template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Configurar el manejo de archivos estáticos
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Configurar las rutas de manejo
	http.HandleFunc("/", home)
	http.HandleFunc("/start_game", startGame)
	http.HandleFunc("/show_game", showGame)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
