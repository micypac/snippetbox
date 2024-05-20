package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	/*
		Define a new command line flag with the name 'addr', a default value of ':4000', and a help text explaining 
		what the flag controls. The value of the flag will be stored in the variable at runtime.
	*/
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	
	/*
		Call the Parse() function to parse the command-line flag.
		This needs to be invoked before the use of variable otherwise it will just use the default value.
	*/
	flag.Parse()

	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" dir
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

