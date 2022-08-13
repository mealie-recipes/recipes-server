# Mealie Recipes Server

## This is for testing purposes ONLY, this should not be used to browser or view any of the URLs exposed by the API.

Mealie Recipes Server is a testing server utilized for load testing and memory profiling to determine the performance of Mealie and/or any other Recipe Parser.

It uses the html data provided by the [recipe-scrapers](https://github.com/hhursev/recipe-scrapers) python package and serves it on the root URL. It also provides a single API endpoint that returns a JSON list of all the endpoints available. It has a few different options for serving the HTML.

## CLI Help

```shell
NAME:
   recipe-server - A testing and development server for serving recipes from various sites

USAGE:
   recipe-server [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h       show help (default: false)
   --host value     The host to bind the server to (default: "127.0.0.1")
   --latency value  latency range randomly applied to requests (e.g. 0-100) in milliseconds (default: "100-1000")
   --port value     The port to listen on (default: "8080")
   --schema value   The schema to use for the server. Currently only used for constructing the urls (default: "http")
```

