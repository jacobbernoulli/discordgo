package discordgo

import "encoding/json"

var (
	Marshal   func(v interface{}) ([]byte, error)   = json.Marshal
	Unmarshal func(src []byte, v interface{}) error = json.Unmarshal
)
