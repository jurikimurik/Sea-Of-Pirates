package main

import (
	http "sea-of-pirates/HTTP"
	source "sea-of-pirates/Source"
)

func main() {
	source.DrawText("Let's begin our sea battle, arr!")
	source.DrawText(http.TestPackage())
}
