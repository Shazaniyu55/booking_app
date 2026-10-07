package main
import (
	"fmt"
	"context"
	"application"

)

func main() {
 app := application.New()

 err := app.Start(context.TODO())
 if err != nil {
  fmt.Printf("Error starting the application: %v\n", err)
 }

}

