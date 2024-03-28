package toolkit

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateRandomPassword(t *testing.T)  {
	tools := ToolKit{}

	password := tools.RandomPasswordGenerator()

	n := len(password)
	assert.EqualValues(t, n > minL && n < maxL, true)
}

// another way of testing from stdout
func Test_CheckTimeElapsed(t *testing.T)  {
	stdOut := os.Stdout
	read, write, _ := os.Pipe()
	os.Stdout = write

	tools := ToolKit{}

	_ = tools.RandomPasswordGenerator()

	write.Close()

	data, _ := io.ReadAll(read)	
	result := string(data)

	os.Stdout = stdOut

	if !strings.Contains(result, "password") {
		t.Errorf("Exception failed running test cases #02")
	}
	
}