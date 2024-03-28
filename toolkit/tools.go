package toolkit

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	minL = 8
	maxL = 30
)

type ToolKit struct {
}

func (t *ToolKit) RandomPasswordGenerator() string {
	currTime := time.Now()

	sourceLetter := []string{
		"abcdefghijklomnopqrstuvwxyz",
		"ABCDEFGHIJKLOMNOPQRSTUVWXYZ",
		"@#$%^&*",
		"01233456789",
	}

	source := rand.NewSource(time.Now().UnixNano())

	randn := rand.New(source)

	passwordLen := minL + randn.Intn(maxL)

	password := []byte{}

	for {
		if len(password) == passwordLen {
			break
		}

		indexI := randn.Intn(4)
		indexJ := randn.Intn(len(sourceLetter[indexI]))
		password = append(password, sourceLetter[indexI][indexJ])
	}

	totalTimeTaken := time.Since(currTime)
	fmt.Printf("Total time taken to generate password: %v\n", totalTimeTaken)

	return string(password)
}
