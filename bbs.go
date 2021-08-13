package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var workdir string
var port string

type Comment struct {
	Name string `json:"Name"`
	Date string `json:"Date"`
	Cont string `json:"Cont"`
}

type Thread struct {
	Title    string    `json:"Title"`
	Id       string    `json:"Id"`
	Comments []Comment `json:"Comments"`
}

func check(err error) {
	// TODO: Gracefully handle errors
	if err != nil {
		panic(err)
	}
}

func threadListHandler(writer http.ResponseWriter, request *http.Request) {
	_, err := os.Stat(workdir + "/dat")
	if err != nil {
		os.Mkdir(workdir+"/dat", os.FileMode(0755))
	}

	file, err := os.OpenFile(workdir+"/dat/threads.csv", os.O_CREATE, os.FileMode(0600))
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
	html, err := template.New("top.html").ParseFiles(workdir + "/template/top.html")
	check(err)
	err = html.Execute(writer, threads)
	check(err)
}

func threadCreateHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.FormValue("title")
	id := uuid.NewString()

	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE

	dat, err := json.Marshal(Thread{Title: title, Id: id})
	check(err)
	err = ioutil.WriteFile(workdir+"/dat/"+id, dat, 0644)
	check(err)

	threadfile, err := os.OpenFile(workdir+"/dat/threads.csv", options, os.FileMode(0600))
	check(err)
	defer threadfile.Close()
	fmt.Fprintf(threadfile, "%s,%s\n", title, id)

	http.Redirect(writer, request, "/", http.StatusFound)
}

func threadViewHandler(writer http.ResponseWriter, request *http.Request) {

	id := request.URL.Query().Get("id")
	datfile := workdir + "/dat/" + id
	_, err := os.Stat(datfile)
	check(err)

	raw, err := ioutil.ReadFile(datfile)
	check(err)
	var thread Thread
	json.Unmarshal(raw, &thread)

	funcMap := template.FuncMap{
		"crlf2br": func(s string) template.HTML {
			return template.HTML(strings.Replace(s, "\r\n", "<br>", -1))
		},
	}
	html, err := template.New("view.html").Funcs(funcMap).ParseFiles(workdir + "/template/view.html")
	check(err)
	check(html.Execute(writer, thread))
}

func addCommentHandler(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	name := request.FormValue("name")
	if name == "" {
		name = "名無しさん"
	}
	comment := request.FormValue("comment")
	wdays := [...]string{"日", "月", "火", "水", "木", "金", "土"}
	t := time.Now()
	date := t.Format("2006-01-02 15:04:05") + " (" + wdays[t.Weekday()] + ")"
	datfile := workdir + "/dat/" + id

	// read dat and add comment
	raw, err := ioutil.ReadFile(datfile)
	check(err)
	var thread Thread
	json.Unmarshal(raw, &thread)
	thread.Comments = append(thread.Comments, Comment{Name: name, Date: date, Cont: comment})

	// write back to dat file
	dat, err := json.Marshal(thread)
	check(err)
	err = ioutil.WriteFile(datfile, dat, 0644)
	check(err)

	http.Redirect(writer, request, "/view?id="+id, http.StatusFound)
}

func main() {
	workdir = os.Getenv("GOLANG_BBS_WORKDIR")
	if workdir == "" {
		workdir, _ = os.Getwd()
	}
	port = os.Getenv("GOLANG_BBS_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Start Server")
	fmt.Println("Port     ", port)
	fmt.Println("Workdir  ", workdir)

	http.HandleFunc("/", threadListHandler)
	http.HandleFunc("/create", threadCreateHandler)
	http.HandleFunc("/view", threadViewHandler)
	http.HandleFunc("/comment", addCommentHandler)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatal(err)
}
