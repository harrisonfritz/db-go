package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"

	// "net/http"
	_ "github.com/lib/pq"
)

const (
	// TODO fill this in directly or through environment variable
	host   = "0.0.0.0"
	port   = 5432
	user   = "user1"
	dbname = "blog"
)

var password string = "1ll10o0oI10I1lll1l1lILIL01!0O00O            \n\n\n \\t\n\tt\t                   "

type Blogpost struct {
	Post_id          string
	Title            string
	Content          string
	Author           string
	Description      string
	Publication_date string
	Card_image_url   string
}

var Blogposts []Blogpost
var psqlconn string
var db *sql.DB
var err error = nil
var defaultImageUrl = "https://images3.alphacoders.com/165/thumb-1920-165265.jpg"

func main() {
	setupEnv()
	db_connect()
	getAll()
	if err != nil {
		log.Fatal(err)
	}
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
	//what does this do?
	defer rowsRs.Close()
	for rowsRs.Next() {
		var post_id string
		var title string
		var content string
		var author string
		var description string
		var publication_date time.Time
		var card_image_url string
		if card_image_url == "" {
			card_image_url = defaultImageUrl
		}
		//why is Scan checking the pointers to those variables?
		err = rowsRs.Scan(&post_id, &title, &content, &author, &description, &publication_date, &card_image_url)
		CheckError(err)

		// formatting the current timestamp in unixDate format
		publication_date_string := publication_date.Format(time.UnixDate)
		// generating a string of the unixdate format
		unix_time, err := time.Parse(time.UnixDate, publication_date_string)
		// formatting that string to be nice
		pubdate := unix_time.Format("January 2, 2006")
		blogpost := Blogpost{post_id, title, content, author, description, pubdate, card_image_url}

		Blogposts = append(Blogposts, blogpost)
		fmt.Println(err)
	}
	SortByDate(Blogposts)
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func SortByDate(blogs []Blogpost) {
	//how doe sthis sort function work?
	sort.Slice(blogs, func(i, j int) bool {
		t1, err1 := time.Parse("January 2, 2006", blogs[i].Publication_date)
		t1a := t1.Format(time.RFC3339)
		t2, err2 := time.Parse("January 2, 2006", blogs[j].Publication_date)
		t2a := t2.Format(time.RFC3339)
		if err1 != nil || err2 != nil {
			fmt.Println("error with SortByDateFunction")
		}
		return t1a > t2a

	})

	fmt.Println("SORTED")
}

// what is r *http.Request is this a pointer to an http.Request object?
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
	//what does mux do?
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

func loadEnv() {
	password = os.Getenv("pgConnP")
	if err != nil {
		fmt.Println("pgConnP environment variable does not exist")
		fmt.Println(err)
	}
}

func setupEnv() {
	loadEnv()
	psqlconn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlconn)

}
