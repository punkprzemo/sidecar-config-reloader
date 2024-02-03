package main

import (
    "flag"
    "fmt"
    "github.com/fsnotify/fsnotify"
    "os"
    "os/exec"
    "os/signal"
    "syscall"
)

var (
    watchDir   = flag.String("watch-dir", ".", "Directory to watch for changes")
    processName = flag.String("process-name", "", "Process name to reload")
)

func main() {
    flag.Parse()

    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer watcher.Close()

    done := make(chan bool)
    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                if event.Op&fsnotify.Write == fsnotify.Write {
                    fmt.Println("Modified file:", event.Name)
                    reloadProcess(*processName)
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                fmt.Println("Error:", err)
            }
        }
    }()

    err = watcher.Add(*watchDir)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    <-done
}

func reloadProcess(processName string) {
    cmd := exec.Command("pkill", "-HUP", processName)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    fmt.Println("Reloading process:", processName)
    if err := cmd.Run(); err != nil {
        fmt.Println("Failed to reload process:", err)
    }