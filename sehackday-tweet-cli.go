package main

import (
    "code.google.com/p/goncurses"
    _ "github.com/mattn/go-sqlite3"
    _ "github.com/araddon/goauth"
    _ "github.com/araddon/httpstream"
    "database/sql"
//    "fmt"
    "log"
    "os"
    "time"
)

func draw(scr goncurses.Window, tweets *[]Tweet) {
    scr.Keypad(true)
    height, width := scr.Maxyx()
    title := "  #sehackday  "
    scr.Move(3,0)
    scr.ClearToBottom()
    for _, t := range *tweets {
        scr.Printf("%s\n    --%s\n\n", t.tweet_body, t.user)
    }
    scr.Move(height-1,0)
    scr.ClearToBottom()
    scr.Move(0,0)
    scr.ClearToEOL()
    scr.MovePrint(0, (width - len(title))/2, title)
    scr.Move(0,0)
    scr.Refresh()
    scr.GetChar()

    //time.Sleep(10 * time.Second)
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
    stmt, err := tx.Prepare("insert into tweets(twitter_id, user, tweet_body, create_time) values(?, ?, ?, ?)")
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

/*
func twitter() (stream <-chan, done <-chan) {

    // make a go channel for sending from listener to processor
    // we buffer it, to help ensure we aren't backing up twitter or else they cut us off
    stream := make(chan []byte, 1000)
    done := make(chan bool)

    _ = stream
    _ = done

    httpstream.OauthCon = &oauth.OAuthConsumer{
            Service:          "twitter",
            RequestTokenURL:  "http://twitter.com/oauth/request_token",
            AccessTokenURL:   "http://twitter.com/oauth/access_token",
            AuthorizationURL: "http://twitter.com/oauth/authorize",
            ConsumerKey:      os.Getenv("SE_CK"),
            ConsumerSecret:   os.Getenv("SE_CS"),
            CallBackURL:      "oob",
            UserAgent:        "go/httpstream",
    }

    at := oauth.AccessToken{Id: "",
            Token:    os.Getenv("SE_OT"),
            Secret:   os.Getenv("SE_OTS"),
            UserRef:  os.Getenv("SE_USER"),
            Verifier: "",
            Service:  "twitter",
    }

    // the stream listener effectively operates in one "thread"/goroutine
    // as the httpstream Client processes inside a go routine it opens
    // That includes the handler func we pass in here
    client := httpstream.NewOAuthClient(&at, httpstream.OnlyTweetsFilter(func(line []byte) {
            stream <- line
            // although you can do heavy lifting here, it means you are doing all
            // your work in the same thread as the http streaming/listener
            // by using a go channel, you can send the work to a
            // different thread/goroutine
    }))

}
*/

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
//    db_add_tweet(db, Tweet{"0", "hichamoaj", "#SEHackDay project with @RageAndQQ - Got it.", time.Now()})
//    db_add_tweet(db, Tweet{"0", "RangeAndQQ", "Pretty successful #SEHackDay . Props to my main man @hichamoaj for our awesome app Got it.", time.Now()})
//    db_add_tweet(db, Tweet{"0", "Nhieuton", "Sweetness. http://calvinhieu.github.io/BezierSnake/  #SEHackDay", time.Now()})
//    db_add_tweet(db, Tweet{"0", "Hatcrab", "seen at #sehackday:\nmuchhash[\"is\"] = \"such\"\nwow", time.Now()})
//    db_add_tweet(db, Tweet{"0", "TheRealKartik", "Just finished my project for #sehackday - an unofficial api for @mint https://t.co/4MHxoEwbhR", time.Now()})
//    db_add_tweet(db, Tweet{"0", "lizuqiliang", "#sehackday #make_it_rain #withcode", time.Now()})
    tweets := db_get_tweets(db)

    draw(scr, tweets)

}
