package utils

import "fmt"

func GlobalCollection(collection string) string {
	return fmt.Sprintf("g_%s", collection)
}
