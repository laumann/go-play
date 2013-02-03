// TODO Parse configuration file with (go-)yaml
package main


type __config struct {
	workDir string
}

// Default configuration
var config = __config{
	workDir: ".book", // The directory in which all the work is done...
}
