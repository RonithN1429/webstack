package main

import (
    "log"
    "fmt"
    "bytes"
    "strings"
    "io/ioutil"
    "net/http"
    "net/url"
    "database/sql"

    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func index(w http.ResponseWriter, r *http.Request) {
    // load index.html from static/
    http.ServeFile(w, r, "./static/index.html")
}

func summerProgGrade(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    // read the form data from the request
    err := r.ParseForm()
    if err != nil {
        log.Fatal(err)
    }

    // iterate over values
    htmlToInsert := ""

    v := r.Form["subjects"]
    for i := 0; i < len(v); i++ {
        htmlToInsert += fmt.Sprintf("<input type='hidden' name='subjects' value='%s'>", v[i])
    }

    if htmlToInsert == "" {
        // redirect back to summerSubjects.html
        http.Redirect(w, r, "/static/summerSubjects.html", http.StatusFound)
        return
    }

    // read the file manually first
    content, err := ioutil.ReadFile("./static/grades-summer-programs.html")
    if err != nil {
        log.Fatal(err)
    }

    // replace the first two instances of '{}' with the html
    // one for next; one for back
    content = bytes.Replace(content, []byte("{}"), []byte(htmlToInsert), 2)
    fmt.Fprintf(w, string(content))

}

func resultsSummerProg(w http.ResponseWriter, r *http.Request) {
    // read the form data from the request
    err := r.ParseForm()
    if err != nil {
        // redirect back to summerSubjects.html
        http.Redirect(w, r, "/static/summerSubjects.html", http.StatusFound)
        return
    }

    subjects := r.Form["subjects"]
    grades := r.Form["grades"]
    costs := r.Form["cost"]

    // check if subject is defined in form
    if subjects == nil || len(subjects) == 0 {
        // redirect back to summerSubjects.html
        http.Redirect(w, r, "/static/summerSubjects.html", http.StatusFound)
        return
    }

    // check if grade and cost is defined in form
    if grades == nil || len(grades) == 0 || costs == nil || len(costs) == 0 {
        // TODO: pass form data to make this work
        params := url.Values{}

        for _, subject := range subjects {
            params.Add("subjects", subject)
        }
        http.Redirect(w, r, 
            "/static/grades-summer-programs.html?" + params.Encode(), http.StatusFound)
        return
    }


    /* dynamically construct the query */

    // Create placeholders for the 'IN' clause
    subjectPlaceholders := make([]string, len(subjects))
    for i := range subjects {
        subjectPlaceholders[i] = "?" // Each ? is a placeholder
    }

    query := "SELECT * FROM summerProgs WHERE " + 
                strings.Repeat(" startGrade <= ? AND  ? <= endGrade AND ", len(grades)) +
                " subject IN (" + strings.Join(subjectPlaceholders, ",") + ") " +
                "OR startGrade IS NULL OR endGrade IS NULL ORDER BY name ASC;"
    if len(costs) == 1 && costs[0] == "paid" {
        query += " AND cost > 0"
    }

    // Prepare arguments
    args := []interface{}{}
    for _, grade := range grades {
        args = append(args, grade, grade)
    }
    for _, subject := range subjects {
        args = append(args, subject)
    }

    // query the database for the subject

    rows, err := db.Query(query, args...)
    if err != nil {
        log.Print("Query error:", err)
    }

    htmlToInsert := ""

    // Iterate over the result set
    for rows.Next() {
        var name string
        var startGrade sql.NullInt32
        var endGrade sql.NullInt32
        var deadline sql.NullString
        var link string
        var cost int
        var scholarship sql.NullString
        var notes sql.NullString
        var subject string

        // Scan the result into variables
        err := rows.Scan(&name, &startGrade, &endGrade, 
                        &deadline, &link, &cost, &scholarship, 
                        &notes, &subject)

        if err != nil {
            log.Print("Error scanning row:", err)
            return
        }


        htmlToInsert += `
            <button onclick="launchModal('` + name + `')" class="program-cards2">
              <div class="nyu-applied-research" 
                data-link="` + link + `"
                data-cost="` + fmt.Sprint(cost) + `"
                data-subject="` + subject + `"`

        // add the grade if its not null
        if startGrade.Valid {
            htmlToInsert += ` data-start-grade="` + fmt.Sprint(startGrade.Int32) + `"`
        }
        if endGrade.Valid {
            htmlToInsert += ` data-end-grade="` + fmt.Sprint(endGrade.Int32) + `"`
        }

        // add the deadline if its not null
        if deadline.Valid {
            htmlToInsert += ` data-deadline="` + deadline.String + `"`
        }

        // add the scholarship if its not null
        if scholarship.Valid {
            htmlToInsert += ` data-scholarship="` + scholarship.String + `"`
        }

        // add the notes if its not null
        if notes.Valid {
            htmlToInsert += ` data-notes="` + notes.String + `"`
        }

        htmlToInsert += `        
              id="` + name + `">` + name + `</div>
              <div class="program-cards-child1"></div>
              <img
                class="lab-items-icon"
                loading="lazy"
                alt=""
                src="./public/` + subject + `@2x.png"
              />
              </button>
        `
    }

    // read the file manually first
    content, err := ioutil.ReadFile("./static/results-page-summer-programs.html")
    if err != nil {
        log.Fatal(err)
    }

    // replace the first instance of '{}' with the html
    content = bytes.Replace(content, []byte("{}"), []byte(htmlToInsert), 1)
    fmt.Fprintf(w, string(content))
}

func main() {
    var err error
    db, err = sql.Open("sqlite3", "./main.db")
    if err != nil {
        log.Fatal(err)
    }

    // create table
    _, err = db.Exec(
        `CREATE TABLE IF NOT EXISTS summerProgs (
            name TEXT PRIMARY KEY NOT NULL,
            start_grade INTEGER,
            end_grade INTEGER,
            deadline DATE,
            link TEXT NOT NULL,
            cost INTEGER NOT NULL,
            scholarship TEXT,
            notes TEXT,
            subject TEXT NOT NULL
        );`,
    )
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/", index)
    http.HandleFunc("/static/grades-summer-programs.html", summerProgGrade)
    http.HandleFunc("/static/results-page-summer-programs.html", resultsSummerProg)

    // fallback to static/
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

    defer db.Close()
}
