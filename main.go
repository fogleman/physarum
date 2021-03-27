package main

import (
	"math/rand"
	_ "net/http/pprof"
	"physarum/physarum"
	"time"
)

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	rand.Seed(time.Now().UTC().UnixNano())

	physarum.Run()
}
