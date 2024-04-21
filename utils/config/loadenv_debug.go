// +build !release

package config

import (
    "os"
    "path"
    "github.com/joho/godotenv"
)

func loadenv() {
    now, _ := os.Getwd()
    check := now
    _ = godotenv.Load()
    for check = path.Join(check, ".."); check != now; check = path.Join(check, "..") {
        _ = godotenv.Load(path.Join(check, ".env"))
        now = check
    }
}
