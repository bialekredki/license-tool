package cmd

import (
	"os"
	"path/filepath"
	"slices"

	log "github.com/sirupsen/logrus"

	"github.com/gobwas/glob"
	"github.com/spf13/cobra"

	"bialekredki/license-tool/common"
	license_headers "bialekredki/license-tool/license-headers"
)

// headersCmd represents the headers command

var disableCommonDirectories = false
var dryRun = false
var ignorePatterns []string

type copyrightsTemplate struct {
	Holder string
	Year   uint16
	path   string
}

func getFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("File %s coudln't be opened.", path)
		return nil
	}
	return file
}

func isPathToBeProcessed(
	path string,
	ignorePatterns *[]glob.Glob,
	includePatterns *[]glob.Glob,
) bool {

	paths := [2]string{path, filepath.Base(path)}

	var isToBeIncluded = true

	var isToBeIgnored = false

	if includePatterns != nil && len(*includePatterns) > 0 {
		for _, path := range paths {
			isToBeIncluded = common.IsMatchingAnyPattern(path, includePatterns)
			if isToBeIncluded {
				break
			}
		}
	}

	if ignorePatterns != nil && len(*ignorePatterns) > 0 {
		for _, path := range paths {
			isToBeIgnored = common.IsMatchingAnyPattern(path, ignorePatterns)
			if isToBeIgnored {
				break
			}
		}
	}

	return isToBeIncluded && !isToBeIgnored
}

func recursivelyCollectFiles(
	path string,
	ignorePatterns *[]glob.Glob,
	includePatterns *[]glob.Glob,
) []string {

	files, directories := common.ListContentOfDirectory(path)

	var collectedFiles []string

	for _, file := range files {
		path, _ := filepath.Abs(filepath.Join(path, file))
		if !isPathToBeProcessed(path, ignorePatterns, includePatterns) {
			continue
		}
		collectedFiles = append(collectedFiles, path)
	}

	for _, directory := range directories {
		path, _ := filepath.Abs(filepath.Join(path, directory))
		if !isPathToBeProcessed(path, ignorePatterns, nil) || slices.Contains(common.DirectoriesToBeIgnored[:], directory) {
			continue
		}
		collectedFiles = append(collectedFiles, recursivelyCollectFiles(path, ignorePatterns, includePatterns)...)
	}
	return collectedFiles
}

var headersCmd = &cobra.Command{
	Use:   "headers",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		common.SetupLogging(Verbosity)

		if !common.IsExistingDirectory(Path) {
			log.Fatalf("Working directory %s doesn't exist.", Path)
		}
		includePatterns := common.CompileListOfPatterns(args...)
		ignorePatterns := common.CompileListOfPatterns(ignorePatterns...)

		collectedFiles := recursivelyCollectFiles(
			Path,
			&ignorePatterns,
			&includePatterns,
		)
		log.Infof("Collected %d files.", len(collectedFiles))
		copyrights := copyrightsTemplate{
			common.GetStringFlagOr(cmd, "holder"),
			common.GetUInt16FlagOr(cmd, "year", common.GetCurrentYear()),
			common.GetStringFlagOr(cmd, "template"),
		}
		if !common.IsExisitingFile(copyrights.path) {
			log.Fatalf("License file %s not found", copyrights.path)
		}
		licenseTemplateContent := common.GetFileContent(copyrights.path)
		template := license_headers.MakeTemplate(licenseTemplateContent, "license")
		header := license_headers.ParseTemplateIntoString(template, copyrights)
		log.Debugf("License header: %s", header)

		templatedLanguages := license_headers.GetTemplatesForCollectedFiles(header, collectedFiles...)
		log.Debug(templatedLanguages)
	},
}

func init() {
	rootCmd.AddCommand(headersCmd)

	headersCmd.Flags().StringArrayVarP(&ignorePatterns, "ignore", "i", nil, "Patterns to be ignored.")
	headersCmd.Flags().BoolVar(&disableCommonDirectories, "disable-common-directories", false, "Disables commonly ignored values like .git or node_modules.")
	headersCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Performs checks only.")
	headersCmd.Flags().Uint16P("year", "y", common.GetCurrentYear(), "Year for copyright template.")
	headersCmd.Flags().StringP("holder", "H", "", "License holder for copyright template.")
	headersCmd.Flags().StringP("template", "t", "license.txt", "License template.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// headersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// headersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
