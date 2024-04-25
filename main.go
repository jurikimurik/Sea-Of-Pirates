package main

import (
	http "sea-of-pirates/HTTP"
	source "sea-of-pirates/source"
)

func main() {
	//source.DrawText("Let's begin our sea battle, arr!", false)
	//source.DrawText(http.TestPackage(), true)
	//source.JSONTest()

	resp := http.Call(http.GET, "stats", nil, nil)
	source.DrawText(string(resp.Body), true)
	//header, _, statusCode, _ := http.Call(http.POST, "game", nil, source.JSONGetDummy())
	//source.DrawText(header.Get("x-auth-token"), true)
	//source.DrawText(strconv.Itoa(statusCode), true)
}
