package main

import (
	http "sea-of-pirates/http"
	source "sea-of-pirates/source"
	"strconv"
)

func main() {
	//source.DrawText("Let's begin our sea battle, arr!", false)
	//source.DrawText(http.TestPackage(), true)
	//source.JSONTest()

	//_, body, statusCode, _ := http.Call(http.GET, "stats", nil, nil)
	//source.DrawText(string(body), true)
	header, _, statusCode, _ := http.Call(http.POST, "game", nil, source.JSONGetDummy())
	source.DrawText(header.Get("x-auth-token"), true)
	source.DrawText(strconv.Itoa(statusCode), true)
}
