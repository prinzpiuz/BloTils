package server

import "fmt"

const version = "v1"

var ping = fmt.Sprintf("%s/ping", version)
var countClap = fmt.Sprintf("%s/count/{url}", version)
