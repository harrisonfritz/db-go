package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	// "net/http"
	_ "github.com/lib/pq"
)

const (
	// TODO fill this in directly or through environment variable
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	dbname   = "myblog"
	password = "Trombone88!"
)

type Blogpost struct {
	Post_id          string
	Title            string
	Content          string
	Author           string
	Description      string
	Publication_date string
}

var Blogposts []Blogpost
var psqlconn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

var db, err = sql.Open("postgres", psqlconn)

func main() {
	db_connect()
	getAll()
	handleRequests()

}
func db_connect() {
	CheckError(err)

	// close database

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")
}

func getAll() {
	// returnEntry := []string{"a", "a", "a", "a", "a", "a"}
	rowsRs, err := db.Query(`SELECT * FROM "posts"`)
	defer rowsRs.Close()
	for rowsRs.Next() {
		var post_id string
		var title string
		var content string
		var author string
		var description string
		var publication_date string
		err = rowsRs.Scan(&post_id, &title, &content, &author, &description, &publication_date)
		CheckError(err)
		blogpost := Blogpost{post_id, title, content, author, description, "testdate"}
		Blogposts = append(Blogposts, blogpost)
		fmt.Println(err)
		fmt.Println(Blogposts)
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")

}

func returnAllPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Blogposts)
}

func returnSinglePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	dashTitle := strings.Replace(vars["title"], " ", "-", -1)

	// Loop over all of our Articles
	// if the article.Id equals the key we pass in
	// return the article encoded as JSON
	for _, post := range Blogposts {
		title := strings.Replace(post.Title, " ", "-", -1)
		if title == dashTitle {
			json.NewEncoder(w).Encode(post)
		}
	}
}

func handleRequests() {
	// http.HandleFunc("/", homePage)
	// http.HandleFunc("/posts", returnAllPosts)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/posts", returnAllPosts)
	myRouter.HandleFunc("/posts/{title}", returnSinglePost)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
