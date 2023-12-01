package aocdownloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func getAOCPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory:", err)
	}
	return filepath.Join(homeDir, ".config", "aoc")
}

func getCachePath(year string, day string) string {
	return filepath.Join(getAOCPath(), year, day)
}

func getSessionPath() string {
	return filepath.Join(getAOCPath(), "session")
}

func getCachedInput(year string, day string) (string, error) {
	path := filepath.Join(getCachePath(year, day), "input")

	inputFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer inputFile.Close()

	input, err := io.ReadAll(inputFile)
	if err != nil {
		return "", err
	}

	return string(input), nil

}

func getSessionID() string {
	sessionFile, err := os.Open(filepath.Join(getSessionPath(), "sessionID"))
	if err != nil {
		panic(err)
	}
	defer sessionFile.Close()

	sessionID, err := io.ReadAll(sessionFile)
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(sessionID))
}

func downloadInput(year string, day string) (string, error) {
	sessionID := getSessionID()

	url := fmt.Sprintf("https://adventofcode.com/%s/day/%s/input", year, day)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.AddCookie(&http.Cookie{Name: "session", Value: sessionID})
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	input, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	dayPath := getCachePath(year, day)
	err = os.MkdirAll(dayPath, 0755)
	if err != nil {
		return "", err
	}

	inputFile, err := os.Create(filepath.Join(dayPath, "input"))
	if err != nil {
		return "", err
	}
	defer inputFile.Close()

	_, err = inputFile.Write(input)
	if err != nil {
		return "", err
	}

	return string(input), nil
}

func GetInput(year string, day string) (string, error) {
	input, err := getCachedInput(year, day)

	if err != nil {
		fmt.Println("Downloading file...")
		input, err = downloadInput(year, day)
		if err != nil {
			return "", err
		}
	} else {
		fmt.Println("Using cached file...")
	}

	return input, nil
}
