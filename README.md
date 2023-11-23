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

#### a. Archivo main.go

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
    }

    jsonBytes, err := json.Marshal(gameData)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    jsonStr := string(jsonBytes)
    fmt.Println(jsonStr)
    fmt.Println(direccionRemota)


    con, _ := net.Dial("tcp", direccionRemota)
    defer con.Close()
    fmt.Fprintln(con, jsonStr)


    http.Redirect(w, r, "/show_game", http.StatusSeeOther)
}

```


