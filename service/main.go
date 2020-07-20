package main

import (
	"fmt"
	"go_launcher_app/shared"
	"log"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	_ "io/ioutil"
)

type App struct {
	Name string
	Icon []byte
	Description string
	Changelong string
}

type Service struct{}

func (_ *Service) Multiply(args shared.Args, reply *int32) error {
	*reply = args.A * args.B
	return nil
}

func (_ *Service) GetFiles(_ *shared.Arg, reply *[]string) error {
	var Apps []App
	
	/*dat, err := ioutil.ReadFile("/tmp/dat")
    if err != nil {
            panic(err)
        }
    fmt.Print(string(dat))*/

	err := filepath.Walk("apps", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		*reply = append(*reply, info.Name())
		return nil
	})
	if err != nil {
		return err
	}

	if len(*reply) > 0 {
		(*reply)[0] = (*reply)[len(*reply)-1]
		*reply = (*reply)[:len(*reply)-1]
	}
	return nil
}

func main() {
	rpc.Register(new(Service))
	listener, e := net.Listen("tcp", ":5090")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Printf("Listening on 5090..")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go rpc.ServeConn(conn)
	}
}
