package main

import "errors"
import "log"
import "os"
import "path"
import "regexp"
import "strconv"
import "time"

// 2006-01-02T15:04:05.999999999Z07:00

var formats = []string{
	"2006-01-02T15:04:05",
	"20060102_150405",
	"IMG_20060102_150405",
	"VID_20060102_150405",
	"VID_20060102_150405",
}

var formatsRePatterns = []string{
	"(\\d{4})-(\\d{2})-(\\d{2})[T_](\\d{2}):(\\d{2}):(\\d{2})",
	"(\\d{4})-(\\d{2})-(\\d{2})-(\\d{2})-(\\d{2})-(\\d{2})",
	"(\\d{4})(\\d{2})(\\d{2})[_ ](\\d{2})(\\d{2})(\\d{2})",
	"(\\d{4})-(\\d{2})-(\\d{2})",
	"(\\d{4})(\\d{2})(\\d{2})",
}

func tryParse(filename string) (*time.Time, error) {
	for _, format := range formats {
		if parsed, err := time.Parse(format, filename); err == nil {
			return &parsed, nil
		}
	}
	return nil, errors.New("unrecognized format")
}

func parseIntArray(strs []string) []int {
	ret := make([]int, len(strs))
	for i, s := range strs {
		ret[i], _ = strconv.Atoi(s)
	}
	return ret
}

func tryParse2(filename string) (*time.Time, error) {
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
	for _, filepath := range os.Args[1:] {
		filename := path.Base(filepath)

		if parsed, err := tryParse2(filename); err == nil {
			if err := os.Chtimes(filepath, *parsed, *parsed); err != nil {
				log.Printf("could not change times of %q: %s\n", filepath, err)
			} else {
				log.Printf("successfully touched %q", filepath)
			}
		} else {
			log.Printf("could not parse %q: %s\n", filepath, err)
		}
	}
}
