package main

import (
	"fmt"
	"os"

	"github.com/magicsong/volc-examples/irsa/internal/service"
)

func main() {
    svc,err := service.NewService()
    if err != nil {
        fmt.Println("Error initializing service:", err)
        os.Exit(1)
    }
    result := svc.DoSomething()
    fmt.Println(result)
    // loop to keep the program running
    for {
        // do nothing, just keep the program running
    }
}

