package main

import (
	"log"

	"net/http"

	"fmt"

	"encoding/json"

	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Animales struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Handlers struct {
	db *sqlx.DB
}

func NewHandlers(db *sqlx.DB) *Handlers {
	return &Handlers{
		db: db,
	}
}

func (h *Handlers) getAnimals(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Queryx("SELECT id, name FROM animales WHERE id=$1", 3)
	if err != nil {
		log.Print(err)
	}
	var aux_animales Animales
	if rows.Next() {
		err := rows.StructScan(&aux_animales)
		if err != nil {
			log.Print(err)
		}
	}
	json.NewEncoder(w).Encode(aux_animales)
	// ---------------------------------------------------
	//
	// var auxAnimal Animales
	// err := h.db.Get(&auxAnimal, "SELECT * FROM animales WHERE id=$1", 3)
	// if err != nil {
	// 	log.Default()
	// }
	// fmt.Fprintf(w, `%v`, auxAnimal)
	// fmt.Printf("  %#v\n", auxAnimal)
	// ---------------------------------------------------
	//  FUNCIONA PARCIALMENTE
	// var nombre string
	// err := h.db.Get(&nombre, "SELECT name FROM animales WHERE id=$1", 2)
	// if err != nil {
	// 	log.Default()
	// }
	// fmt.Fprintf(w, `%v`, nombre)
	// fmt.Printf("  %#v\n", nombre)
}

func (h *Handlers) getAllAnimals(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Queryx("SELECT * FROM animales ORDER BY id ASC")
	if err != nil {
		log.Print(err)
	}
	var aux_animales Animales
	var consults []Animales
	for rows.Next() {
		err := rows.StructScan(&aux_animales)
		if err != nil {
			log.Print(err)
		}
		consults = append(consults, aux_animales)
	}
	json.NewEncoder(w).Encode(consults)
}

func (h *Handlers) getById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id_Params, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Print(err)
	}

	rows, err := h.db.Queryx("SELECT id, name FROM animales WHERE id=$1", id_Params)
	if err != nil {
		log.Print(err)
	}
	var aux_animales Animales
	if rows.Next() {
		err := rows.StructScan(&aux_animales)
		if err != nil {
			log.Print(err)
		}
	}
	json.NewEncoder(w).Encode(aux_animales)
}

func (h *Handlers) createAnimal(w http.ResponseWriter, r *http.Request) {
	//capturo el request body
	captura := json.NewDecoder(r.Body)
	//genero una variable de tipo Struct Animales
	var formData Animales
	// le doy formato a la captura
	captura.Decode(&formData)
	_, err := h.db.NamedExec(`INSERT INTO animales (id, name) VALUES (DEFAULT,:name)`, formData)
	if err != nil {
		log.Print(err)
	}
}
func (h *Handlers) updateAnimal(w http.ResponseWriter, r *http.Request) {
	captura := json.NewDecoder(r.Body)
	var formData Animales
	captura.Decode(&formData)
	log.Print("\n", formData)
	_, err := h.db.NamedExec(`UPDATE animales SET name=:name WHERE id=:id`, formData)
	if err != nil {
		log.Print(err)
	}
}

func initDb() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=ANIMALES2 password=WhySoSerious sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func main() {
	//conne w/ BD
	db := initDb()
	handlers := NewHandlers(db) // tira un return

	router := mux.NewRouter()

	//muestra un animales {id:3}
	router.HandleFunc("/", handlers.getAnimals).Methods("GET")
	//muestra todos los animales
	router.HandleFunc("/animales", handlers.getAllAnimals).Methods("GET")
	//muestra animales por id
	router.HandleFunc("/animales/{id}", handlers.getById).Methods("GET")
	//Crea animales
	router.HandleFunc("/animales", handlers.createAnimal).Methods("POST")
	// actualiza animales
	router.HandleFunc("/animales", handlers.updateAnimal).Methods("PUT")
	fmt.Print("Running on PORT:3000")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		panic("Error: " + err.Error())
	}
}
