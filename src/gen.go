// +build ignore

package main

import (
	"github.com/shurcooL/vfsgen"
	"log"
	"net/http"
)

func main() {
	var fs http.FileSystem = http.Dir("resources")

	err := vfsgen.Generate(fs, vfsgen.Options{
		PackageName:  "main",
		Filename:     "resources.go",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
