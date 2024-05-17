package config

import (
    "os"
    "strings"
)

var CTFDURL []string

func init() {
    loadenv()
    CTFDURL = strings.Split(os.Getenv("CTFDURL"), ",")
}
