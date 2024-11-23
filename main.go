package main

import (
	"fmt"
	"os"
	"strings"
	"bufio"
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var LogFile string = "/var/log/passwords.log"

// if binary is called without password value, that's because it was a successful auth
// replace the last entry with success result instead of default failure result
func editLine(username string) {
	file, err := os.OpenFile(LogFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil{
		os.Exit(0)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	lastLine := lines[len(lines)-1]
	searchStr := fmt.Sprintf(`"PAM_USER=%s"`, username)
	if strings.Contains(lastLine, searchStr) {
		editLine := strings.Replace(lastLine, `"result":"failure"`, `"result":"success"`, 1)
		writer := bufio.NewWriter(file)
		lastLineOffset := -len(lastLine) - 1
		if _, err := file.Seek(int64(lastLineOffset), io.SeekEnd); err != nil {
			log.Error().Err(err).Msg("seek")
		}
		if _, err := writer.WriteString(editLine); err != nil {
			log.Error().Err(err).Msg("write")
		}
		if err := writer.Flush(); err != nil {
			log.Error().Err(err).Msg("flush")
		}
	}
}

func main() {
	file, err := os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		os.Exit(0)
	}
	defer file.Close()

	vars := os.Environ()

	log.Logger = zerolog.New(file).With().Timestamp().Logger()

	// IppSec's original code:
	// var password string
	// fmt.Scanln(&password)
	// Didn't keep that because I found out the hard way
	// that it fails to handle passwords with spaces. Fun!

	// This handles passwords with spaces like you'd expect:
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		os.Exit(0)
	}
	password := string(input)

	if len(password) == 0 {
		// edit the line (if remote user uses a null password, might cause problems)
		editLine(os.Getenv("PAM_USER"))
	} else {
		log.Info().
			Interface("vars", vars).
			Str("result", "failure").
			Str("password", password).
			Msg("pam_exec dump pw")
	}
}
