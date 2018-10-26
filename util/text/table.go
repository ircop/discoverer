package text

import (
	"regexp"
	"strings"
)

// ParseTable - we will try to parse tables only with headers.
//
// Headers examples:
// `Port      Name               Status       Vlan       Duplex  Speed Type`
// `--------- ----------------- ------------- ----------------- ------------ -----`
//
// 1) Split header by first non-space char
// 2) Loop trough splitted parts. Calculate length between first and last (with spaces): this would be max. column len
// 3) Loop over next lines
// 4) If line is empty, or line is footer, break loop
// 5) for range columns { Split lines by column len. Last column should has unlimited len! }
//
// t: text to parse
// header: header matching regex
// footer: footer matching regex
func ParseTable(t string, header string, footer string, exceeding bool, dontTrimPrevLine bool) [][]string {
	rows := make([][]string, 0)

	// we NEED a header regex
	if header == "" {
		return rows
	}

	var reHeader, reFooter *regexp.Regexp

	reHeader, err := regexp.Compile(header)
	if err != nil {
		return rows
	}
	if footer != "" {
		reFooter, err = regexp.Compile(footer)
		if err != nil {
			return rows
		}
	}

	var started bool
	lines := strings.Split(t, "\n")
	var columnsLengths []int64

	for _, line := range lines {
		if !started {
			if reHeader.Match([]byte(line)) {
				started = true
				// parse header
				columnsLengths = parseHeader(line)
			} else {
				continue
			}
		} else {
			// We are started already.
			// Parse current line regarding to columns parameters.

			// But first - stop loop, if we have empty line or footer match
			if strings.Trim(line, " ") == "" || (footer != "" && reFooter.Match([]byte(line))) {
				break
			}

			cols := parseLine(line, columnsLengths, exceeding)
			// assuming first column is always non-empty...
			// If it's empty, we assuming this is continuous previous row
			cols[0] = strings.Trim(cols[0], " ")
			if cols[0] != "" {
				rows = append(rows, cols)
				continue
			}
			// cols[0] == "" empty first column. Loop trough all rows, append cur to prev row
			// if there is no prev.row, we don't know what to do
			if len(rows) == 0 {
				continue
			}
			lastRowIdx := len(rows)-1
			for i := range cols {
				//cols[i] = strings.TrimRight(cols[i], " ")
				if dontTrimPrevLine {
					rows[lastRowIdx][i] += " " + cols[i]
				} else {
					rows[lastRowIdx][i] += cols[i]
				}
			}
		}
	}

	// finally, just loop trough all rows/cols and trim them
	for i := range rows {
		for j := range rows[i] {
			rows[i][j] = strings.Trim(rows[i][j], " ")
		}

		// We have some exclusions :( Like cisco 'sh int status': 'trunk      a-full', where 'a' belongs to previous field :(
		if len(rows[i]) > 5 {
			c4 := rows[i][4]
			c5 := rows[i][5]
			if strings.HasSuffix(c4, " a") && strings.HasPrefix(c5, "-") {
				rows[i][4] = strings.Trim(c4[:len(c4)-2], " ")
				rows[i][5] = "a" + c5
			}
		}
	}

	return rows
}

func parseLine(line string, columnLengths []int64, exceeding bool) []string {
	cols := make([]string, 0)

	chars := strings.Split(line, "")
	for colIdx, colLen := range columnLengths {

		// todo: if LAST char of prev. column is NOT space + if FIRST char of current column is NOT space --> extending prev. column
		if len(chars) > 0 && colIdx > 0 && exceeding { // make multiple IF's to avoid unreadability :)
			prev := cols[colIdx-1]
			if len(prev) > 0 && prev[len(prev)-1:] != " " {
				//fmt.Printf("LAST NOT SPACE: '%s'\n", prev)
				//fmt.Printf("CUR CHARS: '%s'\n", chars)

				// last char of previous column is not space.
				// We should read current chars up to space and add this to prev. line
				for i := 0; i < len(chars); i++ {
					if chars[i] != " " {
						cols[colIdx-1] += chars[i]
						chars = chars[1:]
						i--
					} else {
						break
					}
				}
			}
		}

		// no chars left, add empty column
		if len(chars) == 0 {
			cols = append(cols, "")
			continue
		}

		// col lenght is more then chars left. Append existing chars and strip them.
		if colLen >= int64(len(chars)) {
			col := strings.Join(chars, "")
			//col = strings.TrimRight(col,  " ")
			cols = append(cols, col)
			chars = make([]string, 0)
			continue
		}

		// If this is last column, append all chars and break
		if colIdx == len(columnLengths) - 1 {
			col := strings.Join(chars, "")
			//col = strings.TrimRight(col, " ")
			cols = append(cols, col)
			break
		}

		// we have enough chars, append colLen chars to columns
		col := strings.Join(chars[:colLen], "")
		// we should NOT tim first spaces, but should trim last ones
		//col = strings.TrimRight(col, " ")
		cols = append(cols, col)
		chars = chars[colLen:]
	}

	for i, _ := range cols {
		cols[i] = strings.TrimRight(cols[i], " ")
	}

	return cols
}

func parseHeader(header string) []int64 {
	columns := make([]int64, 0)
	header = strings.TrimRight(header, " ")

	header = strings.Trim(header, "\n")
	var curColumn int64
	chars := strings.Split(header, "")
	columnTextStarted := false
	wasNonSpaces := false
	for _, c := range chars {
		if c == " " {
			// just add this char to column len
			curColumn++
			columnTextStarted = false
		} else {
			// If this is FIRST non-space character, this is probably next column
			// If not first, this is first part of column
			if columnTextStarted || !wasNonSpaces {
				curColumn++
				wasNonSpaces = true
				columnTextStarted =  true
			} else {
				// non-space: this is next column start. BUT only if this is not 1st non-space character in whole string.
				columns = append(columns, curColumn)
				curColumn = 1
				columnTextStarted = true
			}
		}
	}
	// if last column was stopped without spaces, append it also
	if curColumn > 0 && columnTextStarted {
		columns = append(columns, curColumn)
	}

	return columns
}
