package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"serve4/driver"
	"serve4/serve"

	"strconv"
)

var port = flag.Int("port", 3000, "port number to listen")

func main() {
	flag.Parse()

	path, err := filepath.Abs(filepath.Dir(filepath.Join(os.Args[0], "..")))
	if err != nil {
		log.Fatal(err)
		return
	}

	server := serve.NewServer(path, strconv.Itoa(*port))
	server.System = new(driver.FileSystem)
	server.Start()
}
