package text

import (
	"regexp"
	"strings"
)

// Since stupid golang regexps does not have negative lookahead, we should parse output manually =\
func SplitByParts(text string, reString string) ([]string, error) {
	result := make([]string, 0)

	re, err := regexp.Compile(reString)
	if err != nil {
		return result, err
	}


	lines := strings.Split(text, "\n")
	curPart := make([]string,0)
	for _, line := range lines {
		// if current string matches regex, this is next part
		if re.Match([]byte(line)) {
			if len(curPart) > 0 {
				result = append(result, strings.Join(curPart, "\n"))
			}
			curPart = []string{line}
			continue
		}
		curPart = append(curPart, line)
	}

	// and add last part
	result = append(result, strings.Join(curPart, "\n"))

	return result, nil
}
