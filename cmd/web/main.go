package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}


func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	a := &application{
		logger: logger,
	}

	addr := flag.String("addr", ":4000","HTTP network address")
	flag.Parse()
	
		
	logger.Info("starting server", slog.Any("addr",*addr))

	err := http.ListenAndServe(*addr, a.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

