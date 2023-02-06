package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	Id       int       `json:"id"`
	Isbn     string    `json:"isbn"` // unique number for the movie
	Title    string    `json:"title"`
	Dierctor *Dierctor `json:"dierctor"`
}

type Dierctor struct {
	Firstname string `json:"fname"`
	Lastname  string `json:"lname"`
}

// lcoal db
var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	// set the headers
	w.Header().Set("Content-Type", "application/json")

	// encode the data
	json.NewEncoder(w).Encode(movies)
}

func getMovieById(w http.ResponseWriter, r *http.Request) {
	// set the headers
	w.Header().Set("Content-Type", "application/json")

	// get all params
	params := mux.Vars(r)

	moviewId, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Error while parsing the passed movie id", http.StatusInternalServerError)
		return
	}
	for _, movie := range movies {
		if movie.Id == moviewId {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	// set the headers
	w.Header().Set("Content-Type", "application/json")

	// read the data
	newMovie := Movie{}
	err := json.NewDecoder(r.Body).Decode(&newMovie)
	if err != nil {
		http.Error(w, "error while parsing the given movie object", http.StatusBadRequest)
		return
	}
	newMovie.Id = len(movies) + 1
	movies = append(movies, newMovie)
	json.NewEncoder(w).Encode(movies)
}

func updateMoview(w http.ResponseWriter, r *http.Request) {
	// set the headers
	w.Header().Set("Content-Type", "application/json")

	// get the id
	params := mux.Vars(r)

	// parse the movie id
	movieId, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Error while parsing the passed movie id", http.StatusInternalServerError)
		return
	}
	for idx, movie := range movies {
		if movie.Id == movieId {
			// decode the given new movie data and replace this old movie object with new one
			json.NewDecoder(r.Body).Decode(&movies[idx])
			return
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	moviewId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Error while parsing the passed movie id", http.StatusInternalServerError)
		return
	}

	for idx, movie := range movies {
		if movie.Id == moviewId {
			movies = append(movies[:idx], movies[idx+1:]...)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	dierctor1 := Dierctor{Firstname: "fady", Lastname: "gamil"}

	movies = append(movies, Movie{
		Id:       1,
		Isbn:     "movie_1",
		Title:    "sun will raise again",
		Dierctor: &dierctor1,
	}, Movie{
		Id:       2,
		Isbn:     "movie_2",
		Title:    "superman 2",
		Dierctor: &dierctor1,
	})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovieById).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMoview).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Server is running on port 8000")

	log.Fatal(http.ListenAndServe("localhost:8000", r))
}
