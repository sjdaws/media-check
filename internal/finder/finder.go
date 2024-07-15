package finder

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/sjdaws/media-check/internal/config"
)

type Result struct {
	Found []Violation
}

type Violation struct {
	Path string
	Term string
}

func Check(cfg config.Find) (Result, error) {
	result := Result{}

	pattern := []string{"-type", "f"}
	for _, path := range cfg.Exceptions {
		pattern = append(pattern, "-not", "-iwholename", path)
	}

	unique := make(map[string]Violation)

	for _, path := range cfg.Paths {
		defaults := []string{path}
		defaults = append(defaults, pattern...)

		for _, term := range cfg.Block {
			parameters := make([]string, len(defaults))
			copy(parameters, defaults)
			parameters = append(parameters, "-iname", term)

			err := performCheck(unique, term, parameters...)
			if err != nil {
				return result, err
			}
		}
	}

	if len(unique) > 0 {
		for _, item := range unique {
			result.Found = append(result.Found, item)
		}
	}

	return result, nil
}

func performCheck(unique map[string]Violation, term string, parameters ...string) error {
	cmd := exec.Command("/usr/bin/find", parameters...)

	found := make([]string, 0)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdoutpipe: %v", err)
	}

	scanner := bufio.NewScanner(stdout)

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("unable to start find: %v", err)
	}

	for scanner.Scan() {
		found = append(found, scanner.Text())
	}

	if scanner.Err() != nil {
		err = cmd.Process.Kill()
		if err != nil {
			return fmt.Errorf("unable to kill scanner: %v", err)
		}

		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("unable to run find: %v", err)
		}

		return fmt.Errorf("unable to run scanner: %v", scanner.Err())
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("unable to run find: %v", err)
	}

	if len(found) > 0 {
		for _, path := range found {
			unique[path] = Violation{
				Path: path,
				Term: term,
			}
		}
	}

	return nil
}
