package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Article{}
var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"/Users/nurtileu/Documents/snippetbox/ui/static/index.html",
		"/Users/nurtileu/Documents/snippetbox/ui/static/header.html",
		"/Users/nurtileu/Documents/snippetbox/ui/static/footer.html")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}
	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"/Users/nurtileu/Documents/snippetbox/ui/static/create.html",
		"/Users/nurtileu/Documents/snippetbox/ui/static/header.html",
		"/Users/nurtileu/Documents/snippetbox/ui/static/footer.html")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")
	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Ne vse zapolneny")
	}
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES('%s', '%s', '%s')", title, anons, full_text))

	if err != nil {
		panic(err)
	}
	defer insert.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, err := template.ParseFiles(
		"/Users/nurtileu/Documents/snippetbox/ui/static/show.html",
		"/Users/nurtileu/Documents/snippetbox/ui/static/header.html",
		"/Users/nurtileu/Documents/snippetbox/ui/static/footer.html")

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}
	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		showPost = post
	}

	t.ExecuteTemplate(w, "show", showPost)
}

func home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"/Users/nurtileu/Documents/snippetbox/ui/static/index.html",
		"/Users/nurtileu/Documents/snippetbox/ui/static/header.html",
		"/Users/nurtileu/Documents/snippetbox/ui/static/footer.html")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, "index", posts)
}
func contacts(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"/Users/nurtileu/Documents/snippetbox/ui/static/contacts.html",
		"/Users/nurtileu/Documents/snippetbox/ui/static/header.html",
		"/Users/nurtileu/Documents/snippetbox/ui/static/footer.html")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, "contacts", nil)
}
func aboutus(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"/Users/nurtileu/Documents/snippetbox/ui/static/index.html",
		"/Users/nurtileu/Documents/snippetbox/ui/static/header.html",
		"/Users/nurtileu/Documents/snippetbox/ui/static/footer.html")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, "index", posts)
}
func handleFunc() {
	rtr := mux.NewRouter()
	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("/Users/nurtileu/Documents/snippetbox/ui/static"))))
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
	rtr.HandleFunc("/home", home).Methods("GET")
	rtr.HandleFunc("/contacts", contacts).Methods("GET")
	//rtr.HandleFunc("/aboutus", aboutus).Methods("GET")

	http.Handle("/", rtr)
	http.ListenAndServe(":5501", nil)
}

// /Users/nurtileu/Documents/snippetbox/ui/static
func main() {
	handleFunc()
}
