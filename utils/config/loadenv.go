// +build release

package config

import (
    "os"
    "path/filepath"
    "github.com/joho/godotenv"
)

func loadenv() {
    ex, err := os.Executable()
    if err == nil {
        exPath := filepath.Dir(ex)
        os.Chdir(exPath)
    }
    _ = godotenv.Load()
}
