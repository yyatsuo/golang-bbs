package main

import (
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type Thread struct {
	Title string
	Id    string
}

func check(err error) {
	// TODO: Gracefully handle errors
	if err != nil {
		log.Fatal(err)
	}
}

func threadListHandler(writer http.ResponseWriter, request *http.Request) {
	_, err := os.Stat("dat")
	if err != nil {
		os.Mkdir("dat", os.FileMode(0755))
	}

	file, err := os.OpenFile("dat/threads.csv", os.O_CREATE, os.FileMode(0600))
	check(err)
	defer file.Close()
	reader := csv.NewReader(file)
	var threads []Thread

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		check(err)
		threads = append(threads, Thread{Title: record[0], Id: record[1]})
	}
	html, err := template.ParseFiles("template/top.html")
	check(err)
	err = html.Execute(writer, threads)
	check(err)
}

func threadCreateHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.FormValue("title")
	id := uuid.NewString()
	fmt.Println(id)
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("dat/threads.csv", options, os.FileMode(0600))
	check(err)
	defer file.Close()
	fmt.Fprintf(file, "%s,%s\n", title, id)
	http.Redirect(writer, request, "/", http.StatusFound)
}

func main() {
	http.HandleFunc("/", threadListHandler)
	http.HandleFunc("/create", threadCreateHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
