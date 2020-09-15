package main

import (
	"urlShorten/urlshorten"
)

func main() {
	app := new(urlshorten.App)
	app.Run("0.0.0.0:8089", "http://localhost:8089/")
}
