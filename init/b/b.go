package b

import (
	"fmt"

	"github.com/Augustu/go-draft/init/a"
)

func init() {
	fmt.Println("init b")
}

func Print() {
	a.Print()
	fmt.Println("bb")
}
