package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/agentgavel/gavel/internal/oracle"
)

func runOracle(args []string) int {
	fs := flag.NewFlagSet("oracle", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	listen := fs.String("listen", "127.0.0.1:0", "listen address (host:port; port 0 for ephemeral)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	ln, err := net.Listen("tcp", *listen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v\n", err)
		return 1
	}
	addr := ln.Addr().String()
	fmt.Println(addr)

	h := oracle.NewHandler()
	srv := &http.Server{Handler: h}
	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "serve: %v\n", err)
		return 1
	}
	return 0
}
