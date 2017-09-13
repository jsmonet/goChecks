package main

import (
  "net/http"
  "fmt"
  "os"
)

func main() {
  resp, err := http.Post(os.ExpandEnv("http://www.google.com"), "", nil)
  if err != nil {
    fmt.Println("oh noes it died, you are a bad person")
  }
  fmt.Println(resp)
  defer resp.Body.Close()
}
