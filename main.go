package main

import (
	"PanUpload/cmd"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()
	cmd.Execute()
}
