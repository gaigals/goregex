package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const (
	classHighlight = "highlight"

	formatCharID = "char-%d"
)

type Builder struct {
	text        string
	regexp      *regexp.Regexp
	regexErr    string
	matchingIDs [][]int
	result      Result
}

func NewBuilder(text string) *Builder {
	return &Builder{
		text: text,
	}
}

func (builder *Builder) CompileRegex(regexPattern string) *Builder {
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		builder.regexErr = err.Error()
		return builder
	}

	builder.regexp = regex
	return builder
}

func (builder *Builder) MatchString() *Builder {
	if builder.regexErr != "" || builder.text == "" {
		return builder
	}

	matches := builder.regexp.FindStringIndex(builder.text)
	if len(matches) == 0 {
		return builder
	}

	// builder.matchingIDs = builder.regexp.FindAllStringSubmatchIndex(builder.text, -1)

	builder.matchingIDs = [][]int{matches}
	return builder
}

func (builder *Builder) BuildResponse() ([]byte, error) {
	result := Result{
		HTML:  builder.compileAsHTML(),
		Error: builder.regexErr,
	}

	jsonBytes, err := json.Marshal(&result)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}

func (builder *Builder) compileAsHTML() string {
	html := ""
	fmt.Printf("%s|%s|text=>'%s'\n",
		builder.regexErr, builder.matchingIDs, builder.text)
	for idx, char := range builder.text {
		matching := false
		if char == '\n' {
			fmt.Printf("[%d], '\\n'\n", idx)
		} else {
			fmt.Printf("[%d], '%c'\n", idx, char)
		}

		if builder.regexErr == "" {
			for _, matchingIdx := range builder.matchingIDs {
				if len(matchingIdx) != 2 {
					continue
				}

				if idx >= matchingIdx[0] && idx < matchingIdx[1] {
					html += builder.highlightedSPAN(idx, char)
					matching = true
				}

				if idx == matchingIdx[1] {
					break
				}
			}
		}

		if matching {
			fmt.Println("matching, highligthing")
			continue
		}

		html += builder.defaultSPAN(idx, char)
	}

	// return builder.appendEmptyStringOnBreak(html)
	return html
}

func (builder *Builder) appendEmptyStringOnBreak(html string) string {
	if !strings.HasSuffix(builder.text, "\n") {
		return html
	}

	return html + builder.defaultSPAN(len(builder.text), 0)
}

func (builder *Builder) highlightedSPAN(id int, char rune) string {
	return fmt.Sprintf(`<span id="%s" class=%s>%c</span>`,
		fmt.Sprintf(formatCharID, id),
		classHighlight,
		char,
	)
}

func (builder *Builder) defaultSPAN(id int, char rune) string {
	return fmt.Sprintf(`<span id="%s">%c</span>`,
		fmt.Sprintf(formatCharID, id),
		char,
	)
}
