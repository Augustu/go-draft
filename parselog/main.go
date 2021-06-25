package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	f := flag.String("f", "", "log file")

	flag.Parse()

	p := newParse(*f)

	var sqls []sqltime

	for {
		l, ok := p.readline()
		if !ok {
			fmt.Println("read done")
			break
		}

		sd, ok := parseline(l)
		if !ok {
			continue
		}

		st, ok := parsetype(sd)
		if ok {
			sqls = append(sqls, st)
		}
	}

	sort.Slice(sqls, func(i, j int) bool {
		return sqls[i].ms > sqls[j].ms
	})

	var total float64

	for _, s := range sqls {
		total += s.ms
	}
	fmt.Printf("%f querys: %d\n", total, len(sqls))

	for i, s := range sqls {
		fmt.Println(i, s)
		if i >= 100 {
			break
		}
	}

	p.close()
}

type parse struct {
	fileName string
	file     *os.File
	reader   *bufio.Reader
}

func newParse(fn string) *parse {
	f, e := os.Open(fn)
	if e != nil {
		panic(e)
	}

	r := bufio.NewReader(f)

	return &parse{
		fileName: fn,
		file:     f,
		reader:   r,
	}
}

func (p *parse) close() {
	p.file.Close()
}

func (p *parse) readline() (l string, ok bool) {
	lb, _, err := p.reader.ReadLine()
	if err != nil {
		return
	}

	l = string(lb)
	ok = true
	return
}

func parseline(l string) ([]string, bool) {

	// split color
	ss := strings.Split(l, "\033[0m\033[33m[")
	if len(ss) < 2 {
		return nil, false
	}
	m := ss[1]
	s2 := strings.Split(m, "]")

	// for i, s := range s2 {
	// 	fmt.Println(i, s)
	// }

	ms := s2[0][:len(s2[0])-2]
	sql := s2[2]

	return []string{ms, sql}, true
}

type sqltime struct {
	ms  float64
	sql string
}

func parsetype(ss []string) (sqltime, bool) {
	m, err := strconv.ParseFloat(ss[0], 32)
	if err != nil {
		return sqltime{}, false
	}

	return sqltime{
		ms:  m,
		sql: ss[1],
	}, true
}
