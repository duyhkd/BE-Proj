package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)
	http.HandleFunc("/signup", signUp)

	err := http.ListenAndServe(":3333", nil)
	log.Printf("starting server on port 3333")
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "Hello World, This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func signUp(w http.ResponseWriter, r *http.Request) {
	var newUser User
	body, _ := io.ReadAll(r.Body)
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		fmt.Printf("User %s\n", body)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users := make(map[string]User)
	fileContent, err := os.ReadFile("users.json")
	if err != nil {
		os.Create("users.json")
	}
	json.Unmarshal(fileContent, &users)

	_, ok := users[newUser.UserName]

	if !ok {
		user := User{UserName: newUser.UserName, Password: newUser.Password}
		users[newUser.UserName] = user
		jsonString, _ := json.Marshal(users)
		os.WriteFile("users.json", jsonString, os.ModePerm)
	} else {
		resp := make(map[string]string)
		resp["message"] = "User already signed up"
		jsonResp, _ := json.Marshal(resp)
		w.Write(jsonResp)
	}
}
