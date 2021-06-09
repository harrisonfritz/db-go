package main

import (
	"database/sql"
	"fmt"
    // "net/http"
	_ "github.com/lib/pq"
)

const (
	// TODO fill this in directly or through environment variable
    host = "localhost" 
    port = 5432
    user = "postgres"
    dbname = "myblog"
    password = "Trombone88!"
)

type sandbox struct {
    post_id int
    tlte string 
    content string
    author string
    description string
    publication_date string
}
    var psqlconn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

    var db, err = sql.Open("postgres", psqlconn)
func main(){
    db_connect()
    getEntry()
    // http.HandleFunc("/retrieve", retrieveRecord) // (1)
    // http.ListenAndServe(":8080", nil) // (2)
    
}
func db_connect() {
    // connection string
    // psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    
        // open database
    // db, err := sql.Open("postgres", psqlconn)
    CheckError(err)
        
        // close database
    
        // check db
    err = db.Ping()
    CheckError(err)
    
    fmt.Println("Connected!")
}
     


func getEntry(){
    rowsRs, err := db.Query(`SELECT * FROM "posts"`)
    defer rowsRs.Close()
    for rowsRs.Next() {
        var post_id int
        var title string 
        var content string
        var author string
        var description string
        var publication_date string
        err = rowsRs.Scan(&post_id, &title, &content, &author, &description, &publication_date)
        CheckError(err)
        fmt.Println(title, author, content, author, description)
    }

}


func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}

// func retrieveRecord(w http.ResponseWriter, r *http.Request) {

//     fmt.Println("shnoodle")
//     // checks if the request is a "GET" request
//     if r.Method != "GET" {
//         fmt.Println("not get")
//         http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
//     return
//     }

//     // We assign the result to 'rows'
//     fmt.Println(db.Query("SELECT * FROM posts"))

//     rowsRs, err := db.Query("SELECT * FROM posts")
//     if err != nil {
//         fmt.Println("nono")
//         fmt.Println(err)
//         http.Error(w, http.StatusText(500), http.StatusInternalServerError)
//     return
//     }
//     defer rowsRs.Close()
//     for rowsRs.Next() {
//         // var post_id int
//         var title string 
//         // var content string
//         var author string
//         // var description string
//         // var publication_date string


//         err = rowsRs.Scan(&title, &author)
//         fmt.Println(title, author)
//     }


    // creates placeholder of the sandbox
    // snbs := make([]sandbox, 0)

    // we loop through the values of rows
    // for rowsRs.Next() {
        // snb := sandbox{}

    //     err := rowsRs.Scan(&snb.post_id, &snb.title, &snb.content, &snb.content &snb.author, &snb.description, &snb.publication_date)
    //     if err != nil {
    //     log.Println(err)
    //     http.Error(w, http.StatusText(500), 500)
    //     return
    //     }
    //     snbs = append(snbs, snb)
    // }

    // if err = rowsRs.Err(); err != nil {
    //     http.Error(w, http.StatusText(500), 500)
    //     return
    // }

    // // loop and display the result in the browser
    // for _, snb := range snbs {
    //     fmt.Fprintf("Hello")
    // }

// }