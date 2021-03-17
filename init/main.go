package main

import (
	"github.com/Augustu/go-draft/init/a"
	"github.com/Augustu/go-draft/init/b"
)

// init function only run once,
// no matter how many time it is imported
func main() {
	a.Print()
	b.Print()
}
