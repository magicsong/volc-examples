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
}

