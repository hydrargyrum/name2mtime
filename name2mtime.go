package main

import "errors"
import "fmt"
import "os"
import "path"
import "regexp"
import "strconv"
import "time"

var formatsRePatterns = []string{
	"(\\d{4})-(\\d{2})-(\\d{2})[T_](\\d{2}):(\\d{2}):(\\d{2})",
	"(\\d{4})-(\\d{2})-(\\d{2})-(\\d{2})-(\\d{2})-(\\d{2})",
	"(\\d{4})(\\d{2})(\\d{2})[_ ](\\d{2})(\\d{2})(\\d{2})",
	"(\\d{4})-(\\d{2})-(\\d{2})",
	"(\\d{4})(\\d{2})(\\d{2})",
}

func parseIntArray(strs []string) []int {
	ret := make([]int, len(strs))
	for i, s := range strs {
		ret[i], _ = strconv.Atoi(s)
	}
	return ret
}

func tryParse(filename string) (*time.Time, error) {
	// XXX can't have this as a global field?
	formatsRe := make([]*regexp.Regexp, len(formatsRePatterns))
	for i, pattern := range formatsRePatterns {
		formatsRe[i], _ = regexp.Compile(pattern)
	}

	for _, re := range formatsRe {
		subs := re.FindStringSubmatch(filename)
		numbers := parseIntArray(subs)

		switch len(subs) {
		case 4:
			t := time.Date(numbers[1], time.Month(numbers[2]), numbers[3], 0, 0, 0, 0, time.Local)
			return &t, nil
		case 7:
			t := time.Date(numbers[1], time.Month(numbers[2]), numbers[3], numbers[4], numbers[5], numbers[6], 0, time.Local)
			return &t, nil
		}
	}

	return nil, errors.New("unrecognized format")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: name2mtime FILE...")
		os.Exit(64)
	}

	failures := false

	for _, filepath := range os.Args[1:] {
		filename := path.Base(filepath)

		if parsed, err := tryParse(filename); err == nil {
			if err := os.Chtimes(filepath, *parsed, *parsed); err != nil {
				failures = true
				fmt.Fprintf(os.Stderr, "could not change times of %q: %s\n", filepath, err)
			} else {
				fmt.Printf("updated %q\n", filepath)
			}
		} else {
			failures = true
			fmt.Fprintf(os.Stderr, "could not parse date string in %q: %s\n", filepath, err)
		}
	}

	if failures {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
