// SPDX-License-Identifier: MIT OR Unlicense
package processor

import (
	"fmt"
	file "github.com/boyter/go-code-walker"
	"io/ioutil"
)

var Version = "2.0.0 alpha"

// Set by user as command line arguments
var PossibleLicenceFiles = ""
var DirFilePaths []string
var PathBlacklist = ""

var FileOutput = ""
var ExtentionBlacklist = ""
var MaxSize = 50000
var DocumentName = ""
var PackageName = ""
var DocumentNamespace = ""
var Debug = false
var Trace = false

var IncludeBinaryFiles = false
var IgnoreIgnoreFile = false
var IgnoreGitIgnore = false
var IncludeHidden = false
var AllowListExtensions []string
var Format = ""

type Process struct {
	Directory string // What directory are we searching
	FindRoot  bool
}

func NewProcess(directory string) Process {
	return Process{
		Directory: directory,
	}
}

// Process is the main entry point of the command line output it sets everything up and starts running
func (process *Process) StartProcess() {
	lg := NewLicenceGuesser(true, true)
	lg.UseFullDatabase = true

	fileListQueue := make(chan *file.File, 1000)

	fileWalker := file.NewFileWalker(".", fileListQueue)
	fileWalker.IgnoreGitIgnore = true
	fileWalker.IgnoreIgnoreFile = true
	//fileWalker.AllowListExtensions = append(fileWalker.AllowListExtensions, "go")

	go fileWalker.Start()

	for f := range fileListQueue {
		data, err := ioutil.ReadFile(f.Location)
		if err == nil {

			isBinary := false
			// Check if this file is binary by checking for nul byte and if so bail out
			// this is how GNU Grep, git and ripgrep check for binary files
			for _, b := range data {
				if b == 0 {
					isBinary = true
					continue
				}
			}

			if !isBinary {
				fmt.Println()
				fmt.Println(f.Location)
				for _, x := range lg.SpdxIdentify(string(data)) {
					fmt.Println(x.LicenseId, x.ScorePercentage)
				}
				for _, x := range lg.VectorSpaceGuessLicence(data) {
					fmt.Println(x.LicenseId, x.ScorePercentage)
				}
			}
		}

	}
}
