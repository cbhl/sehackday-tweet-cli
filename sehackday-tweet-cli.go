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
    create table if not exists tweets (id integer not null primary key, twitter_id text, user text, tweet_body text, create_time text);
    `
    _, err := db.Exec(sql)
    if err != nil {
        log.Printf("%q: %s\n", err, sql)

        goncurses.End()
        os.Exit(1)
    }
}

type Tweet struct {
    twitter_id string
    user string
    tweet_body string
    create_time time.Time
}

func db_add_tweet(db *sql.DB, t Tweet) {
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }
    stmt, err := tx.Prepare("insert into tweets(twitter_id, user, tweet_body, create_time) values(?, ?, ?)")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    _, err = stmt.Exec(t.twitter_id, t.user, t.tweet_body, t.create_time.UTC().Format(time.RFC3339))
    if err != nil {
        log.Fatal(err)
    }
    tx.Commit()
}

func db_get_tweets(db *sql.DB) *[]Tweet {
    result := make([]Tweet, 0)
    rows, err := db.Query("select id, twitter_id, user, tweet_body, create_time from tweets")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        var id int
        var twitter_id string
        var user string
        var tweet_body string
        var create_time_string string
        rows.Scan(&id, &twitter_id, &user, &tweet_body, &create_time_string)
        var create_time, _ = time.Parse(time.RFC3339, create_time_string)
        var tweet = Tweet{twitter_id, user, tweet_body, create_time}
        result = append(result, tweet)
    }
    rows.Close()
    return &result
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
    //db_add_tweet(db, Tweet{"0", "cbhl", "This isn't a real tweet, Clarisse.", time.Now()})
    tweets := db_get_tweets(db)

    draw(scr)
}
