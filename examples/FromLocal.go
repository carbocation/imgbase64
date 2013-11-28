package main

import (
	"fmt"

	"github.com/carbocation/imgbase64"
)

// Open and convert a local file.
func main() {
	img := imgbase64.FromLocal("test.png")

	fmt.Println(img)
}
