// Command langmuir-iso is the entry point for the Langmuir monolayer adsorption
// isotherm toolkit. It wires the solver packages (internal/iso, internal/comp,
// internal/inverse) to either an HTTP service with a small web UI or a
// one-shot command-line computation over a bundled example dataset.
//
// Subcommands:
//
//	serve    start the HTTP server (default when run with no subcommand)
//	compute  evaluate the example isotherm and print a coverage/amount table
//
// All numerical work lives in the internal packages; this file only parses
// flags, loads data and prints results.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"langmuir-iso/internal/iso"
	"langmuir-iso/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		runServe([]string{"-http", ":8080"})
		return
	}
	switch os.Args[1] {
	case "serve":
		runServe(os.Args[2:])
	case "compute":
		runCompute(os.Args[2:])
	default:
		// No recognised subcommand: treat the whole argument list as serve flags
		// so `go run . -http :8080` still works.
		runServe(os.Args[1:])
	}
}

// runServe starts the HTTP service that exposes /api/* and serves web/.
func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	httpAddr := fs.String("http", ":8080", "HTTP listen address")
	examplePath := fs.String("example", "example/n2-77k.json", "path to the bundled example dataset")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "serve:", err)
		os.Exit(2)
	}
	srv := server.New(*examplePath, http.Dir("web"))
	if err := srv.Start(*httpAddr); err != nil {
		fmt.Fprintln(os.Stderr, "serve error:", err)
		os.Exit(1)
	}
}

// runCompute loads the example dataset and prints a table of coverage and
// adsorbed amount at every listed pressure. It is the offline, no-browser path
// to the same numbers the web UI fetches from /api/curve.
func runCompute(args []string) {
	fs := flag.NewFlagSet("compute", flag.ExitOnError)
	examplePath := fs.String("example", "example/n2-77k.json", "path to the example dataset")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "compute:", err)
		os.Exit(2)
	}
	ex, err := iso.LoadExample(*examplePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "compute error:", err)
		os.Exit(1)
	}
	m, err := ex.Model()
	if err != nil {
		fmt.Fprintln(os.Stderr, "compute error:", err)
		os.Exit(1)
	}
	pts, err := ex.Curve()
	if err != nil {
		fmt.Fprintln(os.Stderr, "compute error:", err)
		os.Exit(1)
	}

	fmt.Printf("# %s\n", ex.Name)
	fmt.Printf("# %s at %.0f K, driver=%s, unit=%s\n", ex.Adsorbate, ex.TemperatureK, ex.Mode, ex.Unit)
	fmt.Printf("# %s\n", m.Summary())
	fmt.Println()
	fmt.Printf("%-14s %-14s %-14s\n", "pressure", "theta", "q")
	fmt.Println("----------------------------------------------")
	for _, p := range pts {
		fmt.Printf("%-14.6g %-14.6g %-14.6g\n", p.X, p.Theta, p.Q)
	}

	// Sanity lines that make the defining limits explicit for the reader.
	fmt.Println()
	fmt.Printf("# half-coverage pressure x_{1/2} = %.6g\n", m.HalfPressure())
	fmt.Printf("# coverage at largest pressure = %.6g (-> 1 as p -> inf)\n", iso.MaxTheta(pts))
}
