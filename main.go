package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var dir *string

func acl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "valid_token" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			return
		}
		next.ServeHTTP(w, r)
		log.Debug().Msg("acl middleware done")
	})
}

func upload(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Error().Err(err)
		return
	}

	defer file.Close()
	out, err := os.Create(fmt.Sprintf("%s/%s", *dir, header.Filename))
	if err != nil {
		log.Error().Err(err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintf(w, "File uploaded successfully: ")
	fmt.Fprintf(w, header.Filename)
}

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	dir = flag.String("dir", ".", "file output directory")
	lst := flag.String("listen", ":24999", "listen, default :24999")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	http.Handle("/upload", acl(http.HandlerFunc(upload)))

	log.Debug().Str("dir", *dir).Str("listen", *lst).Msg("listening...")
	err := http.ListenAndServe(*lst, nil)
	log.Fatal().Err(err).Msg("Done")
}
