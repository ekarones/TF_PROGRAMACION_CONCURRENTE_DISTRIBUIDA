# Trabajo Final Programación Concurrente y Distribuida

### Integrantes:

- Jefferson Espinal Atencia (u201919607)
- Erick Aronés Garcilazo (u201924440)
- Ronaldo Cornejo Valencia (u201816502)

### Docente: 
Carlos Alberto Jara García

### Sección: CC65

### Ciclo: 2023-02

## Resumen

El objetivo es diseñar e implementar una versión modificada del juego de mesa Ludo utilizando programación concurrente, canales y algoritmos distribuidos. Asimismo, se busca simular el juego para un grupo grande de jugadores y personajes, tratando a cada jugador y al tablero de juego como entidades concurrentes separadas (nodos distribuidos) que se comunican a través de puertos y mensajes JSON. La simulación concurrente debe gestionar el movimiento de personajes a través de un laberinto con obstáculos, turnos sincronizados y actualizaciones de progreso en tiempo real. El proyecto también requiere el desarrollo de una interfaz web utilizando un marco de aplicación de página única (SPA) para configurar el juego, mostrar el progreso en tiempo real, presentar gráficamente estadísticas al final del juego y exhibir los créditos del juego. El alcance incluye el uso de terminales virtuales para la red de nodos distribuidos, la implementación de una interfaz gráfica con páginas web SPA, la adopción de un marco ágil (SCRUM) y el uso de herramientas de control de versiones como GitHub para una colaboración efectiva y la gestión del código. El resultado busca proporcionar una simulación concurrente y envolvente del juego de Ludo modificado, demostrando habilidades en sistemas distribuidos, desarrollo web y prácticas de gestión de proyectos.

## Índice

1. Objetivo del Estudiante
2. Capítulo I: Presentación 
3. Capítulo II: Marco Teórico 
4. Capítulo III: Implementación de solución 
5. Capítulo IV: Verificación de la solución 
6. Conclusiones 
7. Recomendaciones 
8. Glosario de términos 
9. Bibliografía 
10. Anexos 

## 1. Objetivo del Estudiante

| Objetivo  | Logrado |
| ------------- | ------------- |
| Demuestra ética profesional en el ejercicio de la ingeniería de software.  | Todos los miembros actuaron con integridad y respeto hacia sus compañeros y las normas éticas de la ingeniería de software. Se manejó la información de manera confidencial y de forma honesta. Asimismo, se consideró los valores éticos en la toma de decisiones.  |
| Demuestra Responsabilidad profesional para el logro de los objetivos  | Cada integrante cumplió con las responsabilidades asumidas y los planes comprometidos en el proyecto.  |
| Emite juicios considerando el impacto de las soluciones de ingeniería de software en el contexto global, impacto social, ambiental y económico  | Se evaluó críticamente la solución propuesta, considerando no sólo su viabilidad técnica, sino también su impacto en una escala más amplia. Esto implica considerar factores sociales, ambientales y económicos para tomar decisiones informadas y éticas.  |

## 2. Capítulo I: Presentación
El trabajo consiste en simular el juego de Ludo utilizando programación concurrente, canales y algoritmos distribuidos para la comunicación entre jugadores y el tablero del juego. Ludo es un juego en el que los jugadores compiten para guiar a sus fichas hasta la meta, a través de un laberinto lleno de obstáculos. 
La simulación debe ser capaz de manejar un grupo de jugadores de manera concurrente y usando algoritmo distribuido, donde la comunicación es a través de puertos y sincronización usando canales.
La simulación debe mostrar el progreso del juego en tiempo real, lo que significa que los jugadores deben recibir actualizaciones sobre el estado del juego.

## 3. Capítulo II: Marco Teórico
* Programación concurrente: Se refiere a la ejecución simultánea de múltiples tareas dentro de un programa. Este enfoque permite que varias operaciones progresen aparentemente al mismo tiempo, mejorando la eficiencia y la capacidad de respuesta de una aplicación. Se basa en la idea de la concurrencia, donde diferentes partes del programa pueden ejecutarse independientemente y de manera concurrente.
* Programación distribuida: Implica el diseño y la implementación de sistemas que operan en entornos distribuidos, donde los componentes de software se ejecutan en múltiples dispositivos interconectados. Este enfoque facilita la colaboración y la comunicación entre diferentes nodos de la red, permitiendo la construcción de sistemas escalables y resilientes.
* Canales: En el contexto de la programación concurrente y distribuida, los canales son mecanismos de comunicación que permiten la transferencia de datos entre diferentes partes de un programa. Los canales son esenciales para la sincronización y la coordinación entre procesos concurrentes, facilitando el intercambio de información de manera segura y eficiente.
* Puertos: En el ámbito de la programación distribuida, son puntos de conexión que permiten la comunicación entre nodos o entidades dentro de un sistema distribuido. Estos puertos actúan como interfaces a través de las cuales los componentes pueden intercambiar información. La gestión adecuada de los puertos es fundamental para garantizar una comunicación efectiva y segura en entornos distribuidos.

## 4. Capítulo III: Implementación de la solución

![image](https://github.com/ekarones/TF_PROGRAMACION_CONCURRENTE_DISTRIBUIDA/assets/66271146/5fced8f4-7a88-4d3d-a36c-b84f7d16748a)

### a. Archivo main.go

initializeGameMap(): En esta función se agregarán los obstáculos en el mapa del juego. Estos obstáculos serán puestos de manera aleatoria y estarán representado por -1.

```go
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
```

start_game(): Primero, se crea un objeto llamado GameData que almacenará la variable de número de jugadores y una lista que simulará el mapa. Luego se transforma el objeto a un archivo json y posteriormente a un string, para luego enviarlo por la conexión. También, se establecerá el número de obstáculos que tendrá el mapa y el jugador que empezará primero por medio de la interfaz gráfica. Finalmente, se inicializa el mapa del juego llamando a la función “initializeGameMap”.

```go
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

```
showGame(): Esta función mostrará lo que sucedió en cada turno de cada jugador, mostrar que ficha movió y si llegó o no a la meta. Esto será posible, por los archivos de texto que se generarán en el archivo player.go.

```go
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
		FileContentRojo     template.HTML
		FileContentAzul     template.HTML
		FileContentVerde    template.HTML
		FileContentAmarillo template.HTML
	}{
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
```

main(): Esta función configura el manejo de archivos estáticos y define rutas para manejar las solicitudes. Finalmente, Inicia el servidor HTTP en el puerto 8080.

```go
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

```

### b. Archivo player.go

Se definen dos estructuras: "Ficha" y "Lanzamiento". La estructura "Ficha" representa las fichas de los jugadores, y "Lanzamiento" representa los resultados de lanzar dos dados. Uno de los datos más importantes es “estado” de la estructura “Ficha” que podemos saber si una ficha se va entrar en una casilla donde hay obstáculo (1) o si ya estuvo en zona obstáculo (2). Además, crean la variable para el nodo remoto, un arreglo del objeto fichas y una lista para guardar el mapa.

```go
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

```

guardarPosicionesEnArchivo(): Se utiliza para guardar las posiciones de las fichas en un archivo de texto. La función recibe el color del jugador como parámetro.

```go
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
```

enviar():  Transformación del objeto GameData ha formato json y luego a tipo string, con el objetivo de enviar el objeto hacia el nodo siguiente.

```go
func enviar(gameData GameData) {
	con, _ := net.Dial("tcp", direccionRemota)
	jsonBytes, _ := json.Marshal(gameData)
	jsonStr := string(jsonBytes)
	defer con.Close()
	fmt.Fprintln(con, jsonStr)
}

```

lanzarDados(): Genera un lanzamiento de dados aleatorio y devuelve un objeto de tipo "Lanzamiento".
```go
func lanzarDados() Lanzamiento {
	valor := rand.Intn(2)
	tiro := Lanzamiento{
		dadoA:   rand.Intn(6) + 1,
		dadoB:   rand.Intn(6) + 1,
		avanzar: valor == 1,
	}
	return tiro
}

```

manejador():  Almacenamiento del mensaje y transformación en un objeto de tipo GameData. En la primera ronda se inicializarán las fichas de los jugadores, en las rondas siguientes se validarán si el jugador actual ha ganado el juego. En caso de que el jugador aún no haya ganado, jugará su turno y enviará el objeto GameData al nodo siguiente. En caso contrario, se imprimirá un mensaje de victoria del jugador ganador.
```go
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
```

pierdeTurno(): Verifica si un jugador pierde su turno debido a un obstáculo en el tablero.

```go
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

```

turnoJugador(): Representa el turno de un jugador. Utiliza canales para coordinar los movimientos de las fichas y realiza cálculos para avanzar las fichas en función de los resultados de los dados. También verifica si alguna ficha ha llegado a la meta.
```go
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

```

main(): Se ingresa los valores del color del jugador,  el puerto actual y el puerto destino .Se crea una lista de canales para las fichas del jugador y se ejecuta de manera concurrente la función manejador().

```go
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

```

## 5. Capítulo IV: Verificación de la solución

En la pantalla de inicio el usuario puede seleccionar la ficha que iniciará el juego y el número de obstáculos que estarán en el mapa. Después de dar al botón “Iniciar Juego”, se genera un estructura “GameData” que al enviarse a los otros nodos dará inicio a las rondas. 

![image](https://github.com/ekarones/TF_PROGRAMACION_CONCURRENTE_DISTRIBUIDA/assets/66271146/89687c03-9790-45e8-ad77-57e4c013fb40)

Después de dar inicio al juego, se generarán 4 cuadros que mostrarán las estadísticas del juego sobre los 4 jugadores.

![image](https://github.com/ekarones/TF_PROGRAMACION_CONCURRENTE_DISTRIBUIDA/assets/66271146/f3f8d2bd-de4d-4efb-b416-67c570dc823a)

En la primera ronda se inicializan las fichas de los jugadores y se comprueban que la información del GameData haya llegado a todos los nodos.

![image](https://github.com/ekarones/TF_PROGRAMACION_CONCURRENTE_DISTRIBUIDA/assets/66271146/981e97dd-3017-4dee-ba36-4f81b37a41bf)

En las siguientes rondas se muestran las posiciones de las fichas y se notifica si el ha perdido un turno.

![image](https://github.com/ekarones/TF_PROGRAMACION_CONCURRENTE_DISTRIBUIDA/assets/66271146/f4d69d29-2002-4694-86fe-9d79312b8f2f)

Finalmente el juego termina cuando un jugador tiene sus 4 fichas en la última casilla (pos 39), se detiene el juego y se imprime un mensaje de victoria.

![image](https://github.com/ekarones/TF_PROGRAMACION_CONCURRENTE_DISTRIBUIDA/assets/66271146/bc89d676-560b-4f37-8d27-70a893fe70b3)

## 6. Conclusiones

* La comprensión y aplicación de la programación concurrente y distribuida son fundamentales en el desarrollo de sistemas modernos. La capacidad de ejecutar tareas simultáneamente y coordinar la comunicación entre nodos distribuidos es esencial para construir aplicaciones eficientes y escalables.
* La utilización de canales como mecanismos de comunicación, junto con la gestión adecuada de puertos, ha demostrado ser crucial en entornos concurrentes y distribuidos. La implementación eficiente de estos elementos facilita una comunicación segura y coordinada entre los diferentes componentes de un sistema.
* La adopción de marcos ágiles, como SCRUM, y herramientas de gestión de versiones, como GitHub, ha demostrado ser efectiva en la coordinación de equipos y el seguimiento del progreso del proyecto. La planificación iterativa y la colaboración a través de estas herramientas contribuyen a un desarrollo más eficiente y organizado.

## 7. Recomendaciones

* Fomentar una comunicación efectiva y una colaboración estrecha entre los miembros del equipo es esencial. El uso de herramientas colaborativas en línea y reuniones regulares puede mejorar la coordinación y la resolución eficiente de problemas.

## 8. Glosario de términos

* Concurrencia: Concepto en programación que se refiere a la ejecución simultánea de múltiples tareas en un programa. Se logra mediante el uso de goroutines y canales.
* Host: Se refiere a una máquina o dispositivo físico que forma parte de una red. Participa en la ejecución de tareas o proporciona un servicio. a través de la red.
* Framework: Conjunto estructurado de herramientas, bibliotecas y convenciones que proporciona una base para el desarrollo de software.

## 9. Bibliografía

* Código Facilito. ¿Qué es la programación concurrente? Recuperado de https://codigofacilito.com/articulos/programacion-concurrente [Consulta:23 de noviembre del 2023]
* Platzi. ¿Cómo funciona la metodología Scrum? Qué es y sus 5 fases. Recuperado de https://platzi.com/blog/metodologia-scrum-fases/ [Consulta: 23 de noviembre del 2023]

## 10. Anexos

* Código
https://github.com/ekarones/TF_PROGRAMACION_CONCURRENTE_DISTRIBUIDA
* Video
https://www.youtube.com/watch?v=xpcmkFCy1Ds




