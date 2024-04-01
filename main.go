package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "time"
)

func main() {
    // Initialize the ticker to trigger every minute
    ticker := time.NewTicker(time.Minute / 10)
    defer ticker.Stop()

    // Create a channel to receive ticker events
    tickerCh := ticker.C

    logFolderPath := "./log/"
    // Initialize the first log file
    filename := generateFilename()
    logFile, err := os.OpenFile(logFolderPath+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error opening log file:", err)
        return
    }
    defer logFile.Close()

    // Commands to execute sequentially
    cmd := exec.Command("sh", "-c", "./ledplayer")

    // Connect command's stdout and stderr to pipes
    cmd.Stdout = logFile
    cmd.Stderr = logFile

    // Start the command
    err = cmd.Start()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Command started successfully. Log stored in", filename)

    // Loop to handle ticker events
    for range tickerCh {
        // Close the current log file
        err := logFile.Close()
        if err != nil {
            fmt.Println("Error closing log file:", err)
            return
        }

        log.Printf("%#v\n", logFile)

        // Generate a new filename for the next minute
        filename := generateFilename()

        // Open the new log file
        logFile, err = os.OpenFile(logFolderPath+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            fmt.Println("Error opening log file:", err)
            return
        }

        // Print the new log file name
        fmt.Println("Log file switched to", filename)

        // Update command's stdout and stderr to write to the new log file
        cmd.Stdout = logFile
        cmd.Stderr = logFile
    }

    // Wait for the command to finish before exiting
    err = cmd.Wait()
    if err != nil {
        fmt.Println("Error waiting for command to finish:", err)
    }
}

func generateFilename() string {
    // Generate filename with the current date and time
    return time.Now().Format("2006-01-02_15-04_05") + "_ledplayer.log"
}