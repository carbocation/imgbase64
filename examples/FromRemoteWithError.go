package main

import (
	"fmt"

	"github.com/carbocation/imgbase64"
)

// For this example a DefaultImage will not be set
func main() {
	img := imgbase64.FromRemote("https://www.example.com/uploads.hipchat.com/46533/311662/7rcjJ1KHy128Dxv/transparent.png")

	fmt.Println(img)
}
