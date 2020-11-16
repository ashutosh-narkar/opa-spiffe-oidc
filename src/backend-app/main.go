package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

var (
	addrFlag = flag.String("addr", ":8080", "address to bind the backend server to")
	logFlag  = flag.String("log", "", "path to log to (empty=stderr)")
)

type Claim struct {
	ClaimStatus                  string `json:"claimStatus,omitempty"`
	Provider                     string `json:"provider,omitempty"`
	EnrolleeResponsibilityAmount int    `json:"enrolleeResponsibilityAmount,omitempty"`
	Procedure                    string `json:"procedure,omitempty"`
}

func serveClaims(w http.ResponseWriter, r *http.Request) {
	claims := []Claim{}

	claims = append(claims, Claim{
		ClaimStatus:                  "Claim Paid",
		Provider:                     "Zara Medical Center",
		EnrolleeResponsibilityAmount: 100,
		Procedure:                    "Cataract surgery",
	})

	claims = append(claims, Claim{
		ClaimStatus:                  "Claim Not Paid",
		Provider:                     "Chanel Medical Center",
		EnrolleeResponsibilityAmount: 1200,
		Procedure:                    "Low back pain surgery",
	})

	claimMap := make(map[string][]Claim)
	claimMap["claims"] = claims
	json.NewEncoder(w).Encode(claimMap)
}

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) (err error) {
	flag.Parse()
	log.SetPrefix("backend> ")
	log.SetFlags(log.Ltime)
	if *logFlag != "" {
		logFile, err := os.OpenFile(*logFlag, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("unable to open log file: %v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	} else {
		log.SetOutput(os.Stdout)
	}
	log.Printf("starting backend server...")

	ln, err := net.Listen("tcp", *addrFlag)
	if err != nil {
		return fmt.Errorf("unable to listen: %v", err)
	}
	defer ln.Close()

	r := chi.NewRouter()
	r.Use(noCache)
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from the backend service !\n"))
	})
	r.Get("/claims", http.HandlerFunc(serveClaims))

	log.Printf("listening on %s...", ln.Addr())
	server := &http.Server{
		Handler: r,
	}
	return server.Serve(ln)
}

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Expires", "0")
		h.ServeHTTP(w, r)
	})
}
