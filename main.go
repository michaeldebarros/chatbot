package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/", handleTalk)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}

func handleTalk(w http.ResponseWriter, r *http.Request) {
	//ler o input vindo do cliente
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var stringBody string
	json.Unmarshal(body, &stringBody)

	//pegar a matrícula no header
	matricula := r.Header.Get("matricula")

	//passar matrícula e input para a função que retornará uma resposta a ser mandada novamente para o cliente
	response := marcacao(matricula, stringBody)

	//mandar mensagem para o cliente
	b, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error while marshaling response")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
	w.Write(b)
}
