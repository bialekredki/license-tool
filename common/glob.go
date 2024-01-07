package common

import (
	"github.com/gobwas/glob"
)

func isMatchingPattern(value string, pattern glob.Glob) bool {
	return pattern.Match(value)
}

func checkPatternGroup(value string, patterns *[]glob.Glob, needsToMatch bool) bool {
	for _, pattern := range *patterns {
		if isMatchingPattern(value, pattern) == needsToMatch {
			return needsToMatch
		}
	}
	return !needsToMatch
}

func IsMatchingAnyPattern(value string, patterns *[]glob.Glob) bool {
	return checkPatternGroup(value, patterns, true)
}

func IsMatchingAllPatterns(value string, patterns *[]glob.Glob) bool {
	return checkPatternGroup(value, patterns, false)
}

func CompileListOfPatterns(values ...string) []glob.Glob {
	var compiledPatterns []glob.Glob
	for _, value := range values {
		compiledPatterns = append(compiledPatterns, glob.MustCompile(value))
	}
	return compiledPatterns
}
