package main

import (
    "code.google.com/p/goncurses"
//    "github.com/mattn/go-sqlite3"
//    "github.com/araddon/httpstream"
//    "database/sql"
    "fmt"
    "log"
    "os"
)

func main() {
    // Initialize the screen.
    // WARNING: ncurses is not thread-safe.
    // TODO(cbhl): isolate ncurses code into a GUI goroutine
    _, err := goncurses.Init()
    if err != nil {
        log.Fatal("init:", err)
        os.Exit(1)
    }
    defer goncurses.End() // Clean up when we're done.

    fmt.Printf("Hello, world.\n")
}
