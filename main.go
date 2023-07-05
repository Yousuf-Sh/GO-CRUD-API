package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users []User

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.ID = len(users) + 1
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}

func readUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for _, user := range users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	json.NewEncoder(w).Encode(users)
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "User not found")
}

// func updateUser(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	id, _ := strconv.Atoi(params["id"])

// 	var updatedUser User
// 	_ = json.NewDecoder(r.Body).Decode(&updatedUser)

// 	for i := range users {
// 		if users[i].ID == id {
// 			// Update the user with the new data
// 			users[i].Name = updatedUser.Name
// 			users[i].Age = updatedUser.Age
// 			json.NewEncoder(w).Encode(users[i])
// 			return
// 		}
// 	}
// 	json.NewEncoder(w).Encode(users)

// }
func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Find the user with the given ID and update the data
	for i := range users {
		if users[i].ID == id {
			updatedUser.ID = id // Preserve the original ID
			users[i] = updatedUser
			json.NewEncoder(w).Encode(updatedUser)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(users)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", readUser).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
