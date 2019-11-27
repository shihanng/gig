package integration

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		log.Fatal(err)
	}

	if err := exec.Command("make", "build").Run(); err != nil {
		log.Fatalf("Fail to run make: %v", err)
	}

	code := m.Run()

	if err := exec.Command("make", "clean").Run(); err != nil {
		log.Fatalf("Fail to run make: %v", err)
	}

	os.Exit(code)
}

func TestCli(t *testing.T) {
	out, err := exec.Command("./gi", "go").CombinedOutput()
	assert.NoError(t, err)
	fmt.Println(string(out))
}
