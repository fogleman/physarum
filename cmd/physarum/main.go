package main

import (
	_ "net/http/pprof"

	"github.com/fogleman/physarum/pkg/physarum"
)

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	// rand.Seed(time.Now().UTC().UnixNano())

	physarum.Run()
}
