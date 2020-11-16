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
	addrFlag = flag.String("addr", ":5000", "address to bind the invoice server to")
	logFlag  = flag.String("log", "", "path to log to (empty=stderr)")
)

type Invoice struct {
	Claims            Claim
	AdjustmentAmount  int
	BillToGroupName   string
	InvoiceStatus     string
	NumberOfEnrollees int
}

type Claim struct {
	ClaimStatus                  string `json:"claimStatus,omitempty"`
	Provider                     string `json:"provider,omitempty"`
	EnrolleeResponsibilityAmount int    `json:"enrolleeResponsibilityAmount,omitempty"`
	Procedure                    string `json:"procedure,omitempty"`
}

func serveInvoices(w http.ResponseWriter, r *http.Request) {
	// this path does not query OPA so return hardcoded invoices

	invoices := []Invoice{}

	invoices = append(invoices, Invoice{
		Claims: Claim{
			ClaimStatus:                  "Claim Paid",
			Provider:                     "Zara Medical Center",
			EnrolleeResponsibilityAmount: 100,
			Procedure:                    "Cataract surgery",
		},
		AdjustmentAmount:  1385,
		BillToGroupName:   "COSTCO",
		InvoiceStatus:     "OPEN",
		NumberOfEnrollees: 100,
	})

	invoices = append(invoices, Invoice{
		Claims: Claim{
			ClaimStatus:                  "Claim Not Paid",
			Provider:                     "Chanel Medical Center",
			EnrolleeResponsibilityAmount: 1200,
			Procedure:                    "Low back pain surgery",
		},
		AdjustmentAmount:  11914,
		BillToGroupName:   "CITY OF PASADENA",
		InvoiceStatus:     "OPEN",
		NumberOfEnrollees: 1474,
	})

	invoiceMap := make(map[string][]Invoice)
	invoiceMap["invoices"] = invoices
	json.NewEncoder(w).Encode(invoiceMap)
}

func serveInvoicesWithOPA(w http.ResponseWriter, r *http.Request) {
	invoices := []*Invoice{}

	invoice1 := Invoice{
		AdjustmentAmount:  1385,
		BillToGroupName:   "COSTCO",
		InvoiceStatus:     "OPEN",
		NumberOfEnrollees: 100,
	}

	invoice2 := Invoice{
		AdjustmentAmount:  11914,
		BillToGroupName:   "CITY OF PASADENA",
		InvoiceStatus:     "OPEN",
		NumberOfEnrollees: 1474,
	}

	// Fetch claims from the Backend service
	req, err := http.NewRequest("GET", "http://backend:8001/claims", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	// TODO: Generate the JWT-SVID for the invoice service. The "aud" claim should contain
	// the token received from the Frontend service. The token received from the
	// Frontend service contains the end-user identity as well as the service SPIFFE ID
	req.Header.Set("serviceId", "invoice_service")
	req.Header.Set("token", r.URL.Query().Get("token")) // TODO: this will be replaced by the JWT-SVID

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var claims []Claim
		json.NewDecoder(resp.Body).Decode(&claims)

		if len(claims) != 0 {
			invoice1.Claims = claims[0]
			invoice2.Claims = claims[1]
		}
	}

	invoices = append(invoices, &invoice1)
	invoices = append(invoices, &invoice2)

	invoiceMap := make(map[string][]*Invoice)
	invoiceMap["invoices"] = invoices
	json.NewEncoder(w).Encode(invoiceMap)
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
		w.Write([]byte("Hello from the invoice service !\n"))
	})
	r.Get("/invoices", http.HandlerFunc(serveInvoices))
	r.Get("/invoices/opa", http.HandlerFunc(serveInvoicesWithOPA))

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
