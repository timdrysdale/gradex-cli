package cmd

type Specification struct {
	Root    string `default:"/usr/local/gradex"`
	Verbose bool   `default:"false"`
}
