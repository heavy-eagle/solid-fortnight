package main

import (
    "log"
    "os"

    "estcli/est"
)

func main() {
    if err := est.ExecuteCLI(os.Args); err != nil {
        log.Fatalf("estcli error: %v", err)
    }
}
