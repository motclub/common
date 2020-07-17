package helpers

import (
	"log"
	"os"
)

func EnsureDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Println(err)
		}
	}
}
