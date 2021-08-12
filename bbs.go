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
	"time"
)

type Comment struct {
	Name string
	Date string
	Cont string
}

type Thread struct {
	Title    string
	Id       string
	Comments []Comment
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
	os.Create("dat/" + id)
	if err == nil {
		fmt.Fprintf(file, "%s,%s\n", title, id)
	} else {
		fmt.Println("Oops! Failed to create dat file!")
	}
	http.Redirect(writer, request, "/", http.StatusFound)
}

func threadViewHandler(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	title := request.URL.Query().Get("title")
	datfile := "dat/" + id
	_, err := os.Stat(datfile)
	if id == "" || title == "" || err != nil {
		fmt.Println("Invalid request id:", id, "title:", title)
		http.Redirect(writer, request, "/", http.StatusFound)
		return
	}

	file, err := os.Open(datfile)
	check(err)
	defer file.Close()

	reader := csv.NewReader(file)
	var comments []Comment

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		check(err)
		comments = append(comments, Comment{Name: record[0], Date: record[1], Cont: record[2]})
	}

	html, err := template.ParseFiles("template/view.html")
	check(err)

	err = html.Execute(writer, Thread{Title: title, Id: id, Comments: comments})
	check(err)
}

func addCommentHandler(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	name := request.FormValue("name")
	title := request.FormValue("title")
	comment := request.FormValue("comment")
	wdays := [...]string{"日","月","火","水","木","金","土"}
	t := time.Now()
	date := t.Format("2006-01-02 15:04:05") + " (" + wdays[t.Weekday()] + ")"

	fmt.Println("Title", title)
	fmt.Println("Id", id)
	fmt.Println("Name", name)
	fmt.Println("Date", date)
	fmt.Println("Comment", comment)

	datfile := "dat/" + id
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile(datfile, options, os.FileMode(0600))
	check(err)
	defer file.Close()
	fmt.Fprintf(file, "%s,%s,%s\n", name, date, comment)

	http.Redirect(writer, request, "/view?id="+id+"&title="+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/", threadListHandler)
	http.HandleFunc("/create", threadCreateHandler)
	http.HandleFunc("/view", threadViewHandler)
	http.HandleFunc("/comment", addCommentHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
