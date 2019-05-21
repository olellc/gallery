package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

type Opts struct {
	Addr       string `long:"addr" default:":8080" description:"TCP network address to listen on"`
	DBPath     string `long:"dbpath" description:"path to the SQLite database file (default: temp file)"`
	KeepSchema bool   `long:"keep-schema" description:"keep existing schema"`
}

func main() {
	var opts Opts
	parser := flags.NewParser(&opts, flags.HelpFlag)

	_, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if opts.DBPath == "" {
		file, err := ioutil.TempFile("", "gallery-")
		if err != nil {
			log.Fatalf("[ERROR] unable to create database file: %v", err)
		}
		opts.DBPath = file.Name()
		file.Close()
		defer os.Remove(opts.DBPath)
	}

	store, err := NewStore(opts.DBPath)
	if err != nil {
		log.Fatalf("[ERROR] unable to create store: %v", err)
	}
	defer store.Close()

	if !opts.KeepSchema {
		err := store.CreateSchema()
		if err != nil {
			log.Fatalf("[ERROR] failed to create schema: %v", err)
		}
	}

	router := Routes(store)

	log.Printf("[INFO] starting server on the address: %s", opts.Addr)
	err = RunUntilSignal(router, opts.Addr)
	if err != nil {
		log.Fatalf("[ERROR] server stopped with the error: %v", err)
	}

	log.Println("[WARN] server stopped successfully")
}
