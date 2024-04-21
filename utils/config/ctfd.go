package config

import (
    "os"
)

var CTFDURL string

func init() {
    loadenv()
    CTFDURL = os.Getenv("CTFDURL")
}
