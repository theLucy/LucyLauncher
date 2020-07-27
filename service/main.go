package main

import (
	"fmt"
	"go_launcher_app/shared"
	"io/ioutil"
	_ "io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"sort"
)

type Service struct{}

func zipFileHander(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Query().Get("zipPath"))
}

func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}

	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

func (_ *Service) GetFiles(_ *shared.Arg, reply *[]shared.App) error {
	apps, err := readDirNames("apps")
	if err != nil {
		return err
	}

	for _, x := range apps {
		var TempApp shared.App

		versions, err := readDirNames("apps/" + x)
		if err != nil {
			return err
		}

		TempApp.Name = x

		for _, y := range versions {

			var TempVersion shared.Version

			if y == "icon.png" {
				dat, err := ioutil.ReadFile(filepath.Join("apps", x, y))
				if err != nil {
					panic(err)
				}
				TempApp.Icon = dat
			} else if y == "desc.txt" {
				dat, err := ioutil.ReadFile(filepath.Join("apps", x, y))
				if err != nil {
					panic(err)
				}
				TempApp.Description = string(dat)
			} else {
				files, err := readDirNames(filepath.Join("apps", x, y))
				if err != nil {
					return err
				}
				TempVersion.Name = y
				for _, z := range files {

					if z == "changes.txt" {
						dat, err := ioutil.ReadFile(filepath.Join("apps", x, y, z))
						if err != nil {
							panic(err)
						}
						TempVersion.Changelog = string(dat)
					} else if z == "app.zip" {
						TempVersion.ArchiveName = filepath.Join("apps", x, y, z)
					}
				}
				TempApp.Versions = append(TempApp.Versions, TempVersion)
			}
		}

		*reply = append(*reply, TempApp)
		fmt.Printf("\n")
	}

	return nil
}

func main() {
	rpc.Register(new(Service))
	rpc.HandleHTTP()
	listener, e := net.Listen("tcp", ":5090")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Printf("Listening on 5090..")

	http.HandleFunc("/", zipFileHander)
	http.Serve(listener, nil)
}
