package app

import (
	"log"
	"net/http"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/config"
)

func TestApp_init(t *testing.T) {
	s := &http.Server{}

	cfg := config.NewConfig()
	err := cfg.Open("../../configs/app/api/local.yaml")
	require.NoError(t, err)

	log.Println("Enabling containers with all environment...")

	cmd := exec.Command("docker-compose", "-f", "../../docker-compose.yml", "up", "-d")

	out, err := cmd.CombinedOutput()

	log.Println("\n" + string(out))

	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			assert.NoError(t, err, string(exitErr.Stderr))
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}

	log.Println("Successfully enabled")

	log.Println("Check running containers...")

	cmd = exec.Command("docker", "ps")

	out, err = cmd.CombinedOutput()

	log.Println("\n" + string(out))

	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			assert.NoError(t, err, string(exitErr.Stderr))
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}

	log.Println("Initializing server...")

	// Change local config for all services to start
	server := NewServer(s, cfg)
	assert.NoError(t, server.init())
	log.Println("Successfully initialized")

	log.Println("Removing all containers...")
	cmd = exec.Command("docker", "ps", "-aq")

	out, err = cmd.CombinedOutput()

	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			require.NoError(t, err, string(exitErr.Stderr))
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}

	containersToRemove := strings.Split(string(out), "\n")
	containersToRemove = containersToRemove[:len(containersToRemove)-1]
	cmd = exec.Command("docker", append(append([]string{"rm"}, "-f"), containersToRemove...)...)

	out, err = cmd.CombinedOutput()
	log.Println("\n" + string(out))

	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			require.NoError(t, err, string(exitErr.Stderr))
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}

	log.Println("Successfully removed")
}
