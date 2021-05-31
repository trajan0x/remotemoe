package main

import (
	"flag"
	"log"

	"github.com/fasmide/remotemoe/http"
	"github.com/fasmide/remotemoe/router"
	"github.com/fasmide/remotemoe/services"
	"github.com/fasmide/remotemoe/ssh"
)

func main() {
	routerData := "routerdata"

	if os.Getenv("STATE_DIRECTORY") != "" {
		routerData = path.Join(os.Getenv("STATE_DIRECTORY"), "routerdata")
	}

	err := os.Mkdir(routerData, 0700)

	// we are not going to be stopping on ErrExists errors
	if errors.Is(err, os.ErrExist) {
		err = nil
	}

	if err != nil {
		log.Fatalf("unable to make directory for router data: %s", err)
	}

	proxy := &http.Proxy{}
	proxy.Initialize()

	server, err := http.NewServer()
	if err != nil {
		panic(err)
	}

	server.Handler = proxy

	services.Serve("http", server)
	services.ServeTLS("https", server)

	sshConfig, err := ssh.DefaultConfig()
	if err != nil {
		log.Fatalf("cannot get default ssh config: %s", err)
	}

	sshServer := &ssh.Server{Config: sshConfig}

	services.Serve("ssh", sshServer)

	// we shall be dealing with shutting down in the future :)
	select {}
}
