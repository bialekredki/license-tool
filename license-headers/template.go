package license_headers

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"

	log "github.com/sirupsen/logrus"
)

const TEMPLATE_READ_ERROR_FORMAT = "Failed to read template from file %s because of %s."

func GetTemplateFromFile(filename string) *template.Template {
	t, err := template.ParseFiles(filename)
	if err != nil {
		log.Fatalf(TEMPLATE_READ_ERROR_FORMAT, filename, err)
	}
	return t
}

func MakeTemplate(content string, name string) *template.Template {
	t := template.New(name)
	t.Parse(content)
	return t
}

func ParseTemplateIntoString(t *template.Template, data any) string {
	var buffer bytes.Buffer
	err := t.Execute(&buffer, data)
	if err != nil {
		log.Fatal(err)
	}
	return buffer.String()
}

func TemplateToRegExPattern(template string) *regexp.Regexp {
	// Escape special characters in the template
	escapedTemplate := regexp.QuoteMeta(template)

	// Replace template placeholders with regex capture groups
	regexPattern := regexp.MustCompile(`{{\.[^}]+}}`)
	escapedTemplate = regexPattern.ReplaceAllString(escapedTemplate, `(.+)`)

	// Create the final regex pattern
	finalPattern := fmt.Sprintf("^%s$", escapedTemplate)

	return regexp.MustCompile(finalPattern)
}
