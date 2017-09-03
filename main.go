package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/alumno", GetRandomAlumnus).Methods("GET")
	router.HandleFunc("/alumnolow", GetRandomAlumnusLowParticipation).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
	log.Fatal(http.ListenAndServe(":8000", router))
}

func loadCSVData() [][]string {
	in, err := os.Open("./database/Lista_logica.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(in)
	r.LazyQuotes = true

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	/* Debug
	fmt.Println("Tamaño de la matriz: " + strconv.Itoa(len(records)))
	fmt.Println("Primer registro: " + records[0][0])
	*/
	//	records[0][0] = "lol"
	//	writeCSV(records)
	return records
}

func getRandomAlumnusLowParticipation(records [][]string) Alumno {
	lowParticipationAlumnus = getLowParticipationsAlumnus(topParticipations(records), records)
	var max = len(lowParticipationAlumnus) - 1
	fmt.Println("Número de alumnos con baja participacion: ", max)
	return lowParticipationAlumnus[random(0, max)]
}

func getLowParticipationsAlumnus(topParticipations int, records [][]string) []Alumno {
	//var idString string = strconv.Itoa(id)
	var lowParticipationAlumnus []Alumno // == nil

	for _, value := range records {
		intValue, err := strconv.Atoi(value[2])
		checkError("Converting from int", err)
		if intValue < topParticipations {
			lowParticipationAlumnus = append(lowParticipationAlumnus,
				Alumno{ID: value[0],
					Nombre:          value[1],
					Participaciones: value[2]})

		}
	}
	return lowParticipationAlumnus
}

func topParticipations(records [][]string) int {
	var n, biggest int = 0, 0

	for _, value := range records {
		intValue, err := strconv.Atoi(value[2])
		checkError("Converting from int", err)
		if intValue > n {
			n = intValue
			biggest = n
		}
	}
	fmt.Println("\nEl valor más grande de asistencias es: " + strconv.Itoa(biggest) + " \n")
	return biggest
}

func getRandomAlumnus(records [][]string) Alumno {
	var max = len(records) - 1
	var randomIndex = random(0, max)
	return Alumno{ID: records[randomIndex][0], Nombre: records[randomIndex][1], Participaciones: records[randomIndex][2]}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func writeCSV(data [][]string) {
	file, err := os.Create("./database/alumnos.csv")
	checkError("Cannot create file", err)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		//fmt.Printf("data: %s\n", value)
		err := writer.Write(value)
		checkError("Cannot write to file", err)
		writer.Flush()

	}
}

func random(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// Alumno perteneciente a una Materia
type Alumno struct {
	ID              string `json:"id,omitempty"`
	Nombre          string `json:"nombre,omitempty"`
	Participaciones string `json:"participaciones,omitempty"`
}

// GetRandomAlumnusLowParticipation regresa un alumno con pocas participaciones aleatoriamente.
func GetRandomAlumnusLowParticipation(w http.ResponseWriter, r *http.Request) {
	var records = loadCSVData()
	alumno = getRandomAlumnusLowParticipation(records)
	fmt.Printf("El alumno es: %v", alumno)
	json.NewEncoder(w).Encode(alumno)
	// Todos los de baja participacion
	//json.NewEncoder(w).Encode(lowParticipationAlumnus)
}

// GetRandomAlumnus regresa un alumno al azar
func GetRandomAlumnus(w http.ResponseWriter, r *http.Request) {
	var records = loadCSVData()
	alumno = getRandomAlumnus(records)
	fmt.Printf("El alumno es: %v", alumno)
	json.NewEncoder(w).Encode(alumno)
}

var alumno Alumno
var lowParticipationAlumnus []Alumno
