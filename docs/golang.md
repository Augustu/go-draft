### Effective Go

#### Formatting

#### Commentary

1. C-style /* */ block comments
2. C++-style // line comments

#### Names（命名）

##### Package names（包名）

```go
package main

import (
    // 导入自带包的形式
    "bytes"
    
    // 导入第三方包的形式
    "gorm.io/driver/mysql"
    "github.com/go-redis/redis/v8"
    // 最后的 v8 指定导入包的版本号，
    // 使用的时候还是用 redis.Client 方式
    rr "github.com/go-redis/redis/v8"
    // rr 是导入包的重命名，
    // 使用时用 rr.Client
    _ "github.com/go-redis/redis/v8"
    // _ 实际代码中没有直接引用这个包，
    // 只是使用了包中的初始化方法 init()
    // 代码中没有直接使用的包需要清理掉，
    // _ 告诉编译器，我们使用了它的 init() 方法

    // 导入自己私有的包和导入第三方包类似
    "github.com/myname/myproject/myversion"
    // 一般这三种类型的导入用空行分隔，
    // 从上到下一次是 自带包 第三方包 私有包
)
```



##### Getters（get 方法）

```go
owner  // 小写，私有
Owner  // 大写，共用，可以被外部访问
```

```go
package owner

// Owner 大写，导出的，外部可以访问，可以当做 geter 使用
type Owner struct {
    Name string		// 大写，导出，外部包可以直接使用
    Sex string
    age uint		// 小写，未导出，外部包不能使用，只能当前包使用
}

```

```go
package main

import (
    // 导入包的示意，只有导入自带的包是这种形式
    "owner"
)

func main() {
    // := 自动判断 o 的类型，Owner 大写开头，相当于是 GetOwner()
    // go 中不使用 GetOwner 这种命名
    o := owner.Owner {
        Name: "Augustu",
        Sex: "male",
        // age: 18		// 成员变量未导出，不能直接赋值
    }
    
    owner := obj.Owner()
    if owner != user {
        obj.SetOwner(user)
    }
}
```



##### Interface names（接口名）

```go
// interface 的方法以 er 结尾
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}
// ref: go/src/io/io.go

type Value struct {
    number int
}

// String 方法名使用 String，而不是 ToString
func (v Value) String() {
    return fmt.Sprintf("%d", v.number)
}
```



##### MixedCaps（命名法）

```go
// 使用驼峰命名法，而不是下划线
MixedCaps
mixedCaps
```



#### Semicolons

#### Control structures

#### Functions

#### Data（数据）

##### Allocation with new（new 内存分配）

new 内置函数分配内存，不初始化内存，而是将内存置为 “零值“ 

```go
type SyncedBuffer struct {
    lock	sync.Mutex		// 零值 是没有上锁的锁
    buffer	bytes.Buffer	// 零值 是空的缓冲区，可以直接使用
}

p := new(SyncedBuffer)	// type *SyncedBuffer
var v SyncedBuffer		// type SyncedBuffer

```



##### Constructors and composite literals（构造函数和组合器）

```go
func NewFile(fd int, name string) *File {
    if fd < 0 {
        return nil
    }

    // f := new(File)
    // f.fd = fd
    // f.name = name
    // f.dirinfo = nil
    // f.nepipe = 0

    // 避免上边繁琐的赋值语句
	f := File{fd, name, nil, 0}

    return &f

    // 写在一行上，不加标签的这种方式，需要把所有的值都写上
    // return &File{fd, name, nil, 0}

    // 加上标签，可以只写不是零值的，没有写出来的值默认是零值
    // return &File{fd: fd, name: name}
}
```

arrays slices maps 组合器

```go
const (
    Enone  = 0
    Eio    = 1
    Einval = 3
)

a := [...]string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
s := []string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
m := map[int]string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}

```



##### Allocation with make（make 内存分配）

make(T, args) 用于创建 slices、maps、channels，这三种类型需要初始化一些数据结构才可以使用

make 返回一个初始化的类型，而不是指针

```go
v := make([]int, 10, 100)
```



##### Arrays

##### Slices

##### Two-dimensional slices

##### Maps

##### Printing

##### Append

#### Initialization

##### Constants

##### Variables

##### The init function（init 函数）

每个源文件都可以定义它自己的无参数无返回值（niladic）的 init 函数，用来设置一些初始状态，实际上每个文件可以有多个 init 函数。

初始化的过程如下：1. 导入的包，2. 包中的变量求值，3. 调用 init 函数

#### Methods（方法）

除了指针和接口类型，其他的类型都可以定义方法，而不只局限于结构体

```go
type ByteSlice []byte

func (slice ByteSlice) Append(data []byte) []byte {
    // Body exactly the same as the Append function defined above.
}

// Append 使用 ByteSlice 指针作为接收器，Append 方法可以覆盖接收器（调用者）的 slice
func (p *ByteSlice) Append(data []byte) {
    slice := *p
    // Body as above, without the return.
    *p = slice
}

// Write 实现 Writer 方法，满足 io.Writer 接口，就可以把 ByteSlice 当作 io.Writer 的方式使用了
func (p *ByteSlice) Write(data []byte) (n int, err error) {
    slice := *p
    // Again as above.
    *p = slice
    return len(data), nil
}

var b ByteSlice
fmt.Fprintf(&b, "This hour has %d days\n", 7)
// 值的方法可以在指针和值中被调用（invoked）
// 指针的方法只能在指针中被调用（invoked）
// 这里的 Writer 方法是定义在指针（*ByteSlice）上的，
// 所以，Fprintf 中要取 b 的地址，用指针的方式调用 Write 方法

// 指针的方法可以用来修改接收器（receiver，在上面的例子中就是 b 本身）
// 值的方法修改的是值的拷贝，调用后修改就被丢弃了
// 当值是可以取地址的时候，编译器会自动取地址
// 例如：b.Write，编译器会自动重写成 (&b).Write
```



#### Interfaces and other types（接口和其他类型）

##### Interfaces（接口）

如果一个类型实现了这里的接口（interface）的方法，那么它就可以被用在这里（if something can do this, then it can be used here）

一个类型可以实现多个接口（例如：实现了 Read 和 Write 方法，就是同时实现了 io.Reader 和 io.Writer 接口，具体要参照接口的定义）

```go
// Sequence 包含了 Len()、Less(i, j int) bool 和 Swap(i, j int) 方法，那么它就实现了 sort.Interface 接口
type Sequence []int

// Methods required by sort.Interface.
func (s Sequence) Len() int {
    return len(s)
}
func (s Sequence) Less(i, j int) bool {
    return s[i] < s[j]
}
func (s Sequence) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

// Copy returns a copy of the Sequence.
func (s Sequence) Copy() Sequence {
    copy := make(Sequence, 0, len(s))
    return append(copy, s...)
}

// Method for printing - sorts the elements before printing.
func (s Sequence) String() string {
    s = s.Copy() // Make a copy; don't overwrite argument.
    sort.Sort(s)
    str := "["
    for i, elem := range s { // Loop is O(N²); will fix that in next example.
        if i > 0 {
            str += " "
        }
        str += fmt.Sprint(elem)
    }
    return str + "]"
}
```

##### Conversions

##### Interface conversions and type assertions

##### Generality

##### Interfaces and methods

##### 

#### The blank identifier

##### The blank identifier in multiple assignment

##### Unused imports and variables

##### Import for side effect

##### Interface checks



#### Embedding

#### Concurrency

##### Share by communicating

##### Goroutines

##### Channels

##### Channels of channels

##### Parallelization

##### A leaky buffer

#### Errors

##### Panic

##### Recover

#### A web server


ref: https://golang.org/doc/effective_go

