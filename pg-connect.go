package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
    "sort"
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
	// timex := (Blogposts[0].Publication_date)
	// t, err := time.Parse(time.RFC3339, timex)
	// timey := t.Format(time.Stamp)
	// fmt.Println("UNIX: ", timey)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("date is of type %T \n", t)
//	fmt.Printf("date is %v", t)
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

func clearBlogposts() {
	Blogposts = nil
	fmt.Println("")
	fmt.Println("Clearing Blogs")
	fmt.Println("")
}

func refreshBlogposts() {
	clearBlogposts()
	fmt.Println("Refreshing ")
	getAll()
}

func getAll() {
	Blogposts = Blogposts[:0]
	rowsRs, err := db.Query(`SELECT * FROM "posts"`)
	defer rowsRs.Close()
	for rowsRs.Next() {
		var post_id string
		var title string
		var content string
		var author string
		var description string
		var publication_date time.Time 
		err = rowsRs.Scan(&post_id, &title, &content, &author, &description, &publication_date)
		CheckError(err)

        // formatting the current timestamp in unixDate format
        publication_date_string := publication_date.Format(time.UnixDate)
        // generating a string of the unixdate format
        unix_time, err := time.Parse(time.UnixDate, publication_date_string)
        // formatting that string to be nice
        pubdate := unix_time.Format("January 2, 2006") 
		blogpost := Blogpost{post_id, title, content, author, description, pubdate}

		Blogposts = append(Blogposts, blogpost)
		fmt.Println(err)
	}
    SortByDate(Blogposts)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func SortByDate(blogs []Blogpost) {
    sort.Slice(blogs, func(i, j int) bool {
        t1, err1 := time.Parse("January 2, 2006", blogs[i].Publication_date)
        t1a := t1.Format(time.RFC3339)
        t2, err2 := time.Parse("January 2, 2006", blogs[j].Publication_date)
        t2a := t2.Format(time.RFC3339)
        if err1  != nil || err2 != nil {
            fmt.Println("error with SortByDateFunction")
        }
        return t1a > t2a 

        })

    fmt.Println("SORTED")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Blogposts[0].Publication_date)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")

}

func returnAllPosts(w http.ResponseWriter, r *http.Request) {
	getAll()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Blogposts)
}

func returnSinglePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	dashTitle := strings.ToLower(strings.Replace(vars["title"], " ", "-", -1))

	// Loop over all of our Articles
	// if the article.Id equals the key we pass in
	// return the article encoded as JSON
	for _, post := range Blogposts {
		title := strings.ToLower(strings.Replace(post.Title, " ", "-", -1))
		if title == dashTitle {
			json.NewEncoder(w).Encode(post)
		}
	}
}

func handleRequests() {
	refreshBlogposts()
	fmt.Println("SHNOODIEDOODLE")
	// http.HandleFunc("/", homePage)
	// http.HandleFunc("/posts", returnAllPosts)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/posts", returnAllPosts)
	myRouter.HandleFunc("/posts/{title}", returnSinglePost)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
	fmt.Println(Blogposts)
}
