package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

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

	/*
		Use log.New() to create a logger for writing information messages. This takes 3 parms:
		destination to write logs to, a string prefix for message, and flags to indicate addl information to include.
	*/
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	/*
		Initialize new instance of our app struct containing the dependencies.
	*/
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}


	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" dir
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	/*
		Initialize a new http server struct. 
		Set the network address, handler, and errorLog fields. 
		Server now uses the custom errorLog logger in the event of any problems.
	*/

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	// err := http.ListenAndServe(*addr, mux)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

