package main

import (
    "code.google.com/p/goncurses"
    _ "github.com/mattn/go-sqlite3"
//    "github.com/araddon/httpstream"
    "database/sql"
//    "fmt"
    "log"
    "os"
    "time"
)

func draw(scr goncurses.Window) {
    scr.Keypad(true)
    _, width := scr.Maxyx()
    title := "#sehackday"
    scr.MovePrint(0, (width - len(title))/2, title)
    scr.Move(1,0)
    scr.Print("Hello, Clarisse.")
    scr.Refresh()
    //scr.GetChar()

    time.Sleep(1 * time.Second)
}

func db_init(db *sql.DB) {
    sql := `
    create table if not exists tweets (id integer not null primary key, user text, tweet_body text, create_time text);
    `
    _, err := db.Exec(sql)
    if err != nil {
        log.Printf("%q: %s\n", err, sql)

        goncurses.End()
        os.Exit(1)
    }
}

func main() {
    // Initialize the screen.
    // WARNING: ncurses is not thread-safe.
    // TODO(cbhl): isolate ncurses code into a GUI goroutine
    scr, err := goncurses.Init()
    if err != nil {
        log.Fatal("init:", err)
        os.Exit(1)
    }
    defer goncurses.End() // Clean up when we're done.

    db, err := sql.Open("sqlite3", "./cms.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    db_init(db)

    draw(scr)
}
