package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hay-kot/yal"
	"github.com/mealie-recipes/recipes-server/pkgs/server"
	"github.com/urfave/cli/v2"
)

func init() {
	// Set the seed for the random number generator
	rand.Seed(time.Now().UnixNano())
}

type web struct {
	mux *http.ServeMux
}

func (w *web) GET(path string, handler http.Handler) {
	w.mux.Handle(path, MwLogger(handler))
	yal.Infof("GET %s", path)
}

type availableRecipes struct {
	Urls []string `json:"recipes"`
}

//go:embed html
var html embed.FS

func main() {
	app := &cli.App{
		Name:  "recipe-server",
		Usage: "A testing and development server for serving recipes from various sites",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "schema",
				Value: "http",
				Usage: "The schema to use for the server. Currently only used for constructing the urls",
			},
			&cli.StringFlag{
				Name:  "host",
				Value: "127.0.0.1",
				Usage: "The host to bind the server to",
			},
			&cli.StringFlag{
				Name:  "port",
				Value: "8080",
				Usage: "The port to listen on",
			},
			&cli.StringFlag{
				Name:  "latency",
				Value: "100-1000",
				Usage: "latency range randomly applied to requests (e.g. 0-100) in milliseconds",
			},
		},
		Action: func(c *cli.Context) error {
			// ============================================================
			// Setup

			var (
				schema       = c.String("schema")
				host         = c.String("host")
				port         = c.String("port")
				latencyRange = c.String("latency")
			)

			var totalRecipes int
			urls := []string{}

			file, err := html.ReadDir(("html"))
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range file {
				if f.IsDir() {
					continue
				}
				url := fmt.Sprintf("%s://%s:%s/%s", schema, host, port, f.Name())
				urls = append(urls, url)

				totalRecipes++
			}

			recipeHtml, err := fs.Sub(html, "html")
			if err != nil {
				log.Fatal(err)
			}

			// ============================================================
			// Web Server Setup

			mux := http.NewServeMux()
			w := web{mux}

			yal.Infof("Found %d recipes, mounting to root handler", totalRecipes)

			latencyMW := MwLatency(latencyRange)

			w.GET("/", http.StripPrefix("/", latencyMW(http.FileServer(http.FS(recipeHtml)))))
			w.GET("/api/v1/available", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(availableRecipes{Urls: urls})
			}))

			svr := server.NewServer(host, port)

			yal.Infof("Starting server on %s:%s", host, port)
			err = svr.Start(mux)

			yal.Info("Server shutdown")

			if err != nil {
				yal.Errorf("Error: %s", err)
			}

			return err
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func MwLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		yal.Infof("[%s] %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Mw Latency is a middleware factory that returns a middleware that adds latency to the request
func MwLatency(latencyRange string) func(next http.Handler) http.Handler {
	minMax := strings.Split(latencyRange, "-")

	if len(minMax) != 2 {
		yal.Errorf("Invalid latency range: %s", latencyRange)
	}

	min, err := strconv.Atoi(minMax[0])
	if err != nil {
		yal.Errorf("Invalid latency range: %s", latencyRange)
	}
	max, err := strconv.Atoi(minMax[1])
	if err != nil {
		yal.Errorf("Invalid latency range: %s", latencyRange)
	}

	// function to generate a random number between min and max
	randInt := func(min, max int) int {
		return min + rand.Intn(max-min)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			latency := randInt(min, max)
			yal.Infof("Applying latency of %dms", latency)
			time.Sleep(time.Duration(latency) * time.Millisecond)
			next.ServeHTTP(w, r)
		})
	}
}
