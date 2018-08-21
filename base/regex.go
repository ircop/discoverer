package discoverer

import (
	"regexp"
	"strings"
	"fmt"
)

// ParseSingle parses output for single regex iteraction
func (p *Generic) ParseSingle(r *regexp.Regexp, str string) map[string]string {
	results := make(map[string]string)

	matches := r.FindStringSubmatch(str)
	count := len(matches)
	names := r.SubexpNames()

	for i, name := range names {
		if name != "" && count >= (i+1) {
			results[name] = strings.Trim(matches[i], " ")
		}
	}

	return results
}

// ParseMultiple parses output for multiple regex iteractions
func (p *Generic) ParseMultiple(r *regexp.Regexp, str string) []map[string]string {
	results := make([]map[string]string,0)

	matches := r.FindAllStringSubmatch(str, -1)
	names := r.SubexpNames()

	for _, match := range matches {
		curMap := make(map[string]string)
		for i, name := range names {
			if name != "" {
				curMap[name] = match[i]
			}
		}

		results = append(results, curMap)
	}

	return results
}

// CompileRegexs used for compiling 2+ regexs at the same time
func (p *Generic) CompileRegexps(patterns map[string]string) (map[string]*regexp.Regexp, error) {
	results := make(map[string]*regexp.Regexp)
	for name, str := range patterns {
		re, err := regexp.Compile(str)
		if nil != err {
			return results, fmt.Errorf("Cannot compile '%s' regexp: %s", name, err.Error())
		}
		results[name] = re
	}

	return results, nil
}