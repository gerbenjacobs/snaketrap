package core

import (
	"fmt"
	"path/filepath"

	"regexp"
	"strings"
)

const (
	StateFileFolder       = "cache"
	StateFileCharReplacer = "-"
)

// partly based on https://github.com/parshap/node-sanitize-filename/blob/master/index.js
var (
	StateFileRegex       = regexp.MustCompile(`[\/\?<>\\:\*\|":!\s.]`)
	StateFileRegexDashes = regexp.MustCompile(`--+`)
)

func hashBotStateFile(name string) string {
	lower := strings.ToLower(name)
	dashes := StateFileRegex.ReplaceAllString(lower, StateFileCharReplacer)
	singles := StateFileRegexDashes.ReplaceAllString(dashes, StateFileCharReplacer)
	trimmed := strings.Trim(singles, "-")
	return filepath.Join(StateFileFolder, fmt.Sprintf("%s.json", trimmed))
}
