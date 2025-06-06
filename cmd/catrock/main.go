package main

import (
    "catRock/cmd/catrock/internal/commands"
    "os"
)

const version = "0.2.1"

func main() {
    if err := commands.Execute(version); err != nil {
        os.Exit(1)
    }
}