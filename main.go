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

//trae un sólo animal (id:3)
func (h *Handlers) getAnimal(w http.ResponseWriter, r *http.Request) {
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
	/* -----------------OTHER WAY--------------------
	var auxAnimal Animales
	err := h.db.Get(&auxAnimal, "SELECT * FROM animales WHERE id=$1", 3)
	if err != nil {
		log.Default()
	}
	fmt.Fprintf(w, `%v`, auxAnimal)
	fmt.Printf("  %#v\n", auxAnimal)
	 ---------------------------------------------------*/
}

//trae todos los animales
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

//trae un animal por id
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

// crea un animal por id
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

// actualiza un animal por id
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

func (h *Handlers) deleteById(w http.ResponseWriter, r *http.Request) {
	captura := json.NewDecoder(r.Body)
	var formData Animales
	captura.Decode(&formData)
	_, err := h.db.NamedExec(`DELETE FROM animales WHERE id=:id`, formData)
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
	router.HandleFunc("/api", handlers.getAnimal).Methods("GET")
	//muestra todos los animales
	router.HandleFunc("/api/animales", handlers.getAllAnimals).Methods("GET")
	//muestra animales por id
	router.HandleFunc("/api/animales/{id}", handlers.getById).Methods("GET")
	//Elimina un animal por ID
	router.HandleFunc("/api/animales", handlers.deleteById).Methods("DELETE")
	//Crea animales
	router.HandleFunc("/api/animales", handlers.createAnimal).Methods("POST")
	// actualiza animales
	router.HandleFunc("/api/animales", handlers.updateAnimal).Methods("PUT")
	fmt.Print("Running on PORT:9041\n")
	err := http.ListenAndServe(":9041", router)
	if err != nil {
		panic("Error: " + err.Error())
	}
}
