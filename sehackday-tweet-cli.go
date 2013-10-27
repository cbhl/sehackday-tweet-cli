package main

import (
    "code.google.com/p/goncurses"
//    "github.com/mattn/go-sqlite3"
//    "github.com/araddon/httpstream"
//    "database/sql"
//    "fmt"
    "log"
    "os"
)

func draw(scr goncurses.Window) {
    scr.Keypad(true)
    _, width := scr.Maxyx()
    title := "#sehackday"
    scr.MovePrint(0, (width - len(title))/2, title)
    scr.Move(1,0)
    scr.Print("Hello, world.")
    scr.Refresh()
    scr.GetChar()
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

    draw(scr)
}
