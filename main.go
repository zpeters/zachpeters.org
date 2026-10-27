// main is the entry point. The site itself is in internal/site.
package main

import (
	"log"

	"zachpeters.org/website/internal/site"
)

// version is the release version, set at build time with -ldflags
var version = "dev"

func main() {
	log.Fatalln(site.Run(version, ":8080"))
}
