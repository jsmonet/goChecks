package main
import (
  "fmt"
  "os"
  // "os/exec"
)

func main() {
  fmt.Println(len(os.Args))
  var posit int
  posit = 1
  for posit < len(os.Args) {
    fmt.Println(os.Args[posit:])
  }
}
