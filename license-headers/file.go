package license_headers

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type CommentStyle uint8
type OperateOnFileCallback func(filename string) (bool, error)

type Language struct {
	name                string
	commentStyle        CommentStyle
	customCommentFormat string
}

type Templated struct {
	filenames     []string
	regexPattern  *regexp.Regexp
	templatedText string
}

const (
	Unknown   CommentStyle = iota
	SlashStar CommentStyle = iota
	Hashtag   CommentStyle = iota
	HTML      CommentStyle = iota
	Text      CommentStyle = iota
)

var languageExtensions = map[string]Language{
	"py":   {"Python", Hashtag, ""},
	"go":   {"Golang", SlashStar, ""},
	"js":   {"JavaScript", SlashStar, ""},
	"ts":   {"TypeScript", SlashStar, ""},
	"cpp":  {"C++", SlashStar, ""},
	"c":    {"C", SlashStar, ""},
	"java": {"Java", SlashStar, ""},
	"html": {"HTML", HTML, ""},
	"txt":  {"Text", Text, ""},
}

func (mcs CommentStyle) IsUnkown() bool {
	return mcs == Unknown
}

func (mcs CommentStyle) SupportsMultilineComments() bool {
	return !(mcs.IsUnkown() || mcs == Hashtag)
}

func getExtensionFromFile(filename string) (string, error) {
	path := filepath.Clean(filename)
	result := strings.Split(path, ".")
	if len(result) < 2 {
		return "", errors.New("failed to strip extension from file")
	}
	return result[len(result)-1], nil
}

func guessLanguage(filename string) (Language, error) {
	extension, err := getExtensionFromFile(filename)
	if err != nil {
		return Language{}, err
	}
	language, ok := languageExtensions[extension]
	if !ok {
		return Language{}, errors.New("language not defined for given extension")
	}
	return language, nil
}

func GetTemplatesForCollectedFiles(template string, filenames ...string) map[Language]Templated {
	var templateLanguages = make(map[Language]Templated, 0)
	for _, filename := range filenames {
		language, err := guessLanguage(filename)
		if err != nil {
			continue
		}
		log.Debugf("File %s was detected to be a %s file", filename, language.name)
		tl, ok := templateLanguages[language]
		if !ok {
			templateLanguages[language] = Templated{[]string{filename}, TemplateToRegExPattern(template), template}
		} else {
			tl.filenames = append(tl.filenames, filename)
			templateLanguages[language] = tl
		}
	}
	return templateLanguages
}
