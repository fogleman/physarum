package physarum

import (
	"time"
)

type Settings map[string]interface{}

var DefaultSettings = Settings{
	"width":       4096,
	"height":      2048,
	"initType":    "random_circle_random",
	"particles":   1 << 23,
	"fps":         60,
	"seed":        time.Now().UTC().UnixNano(),
	"output_file": "out",
}

func NewSettings() Settings {
	// TODO: Stub, replace with json read/write
	return DefaultSettings
	// return Settings{}
}
