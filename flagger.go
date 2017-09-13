package main

import (
	"flag"
	"fmt"
	//  "os"
)

func main() {
	stringArg := flag.String("vol", "", "enter any string here")
	float64Arg := flag.Float64("flo", 50.000, "enter any float64-compliant number")

	flag.Parse()

	fmt.Println(*stringArg)
	fmt.Println(*float64Arg)

	if *float64Arg > 20 {
		fmt.Println("it is gooood")
	} else {
		fmt.Println("it is baaaad")
	}
}
