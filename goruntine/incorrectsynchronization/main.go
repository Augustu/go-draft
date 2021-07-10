package main

var a, b int

func f() {
	a = 1
	b = 2
}

func g() {
	print(b)
	print(a)
}

func main() {
	go f()
	g()
}

// var a string
// var done bool

// func setup() {
// 	a = "hello, world"
// 	done = true
// }

// func main() {
// 	go setup()
// 	for !done {
// 	}
// 	print(a)
// }
