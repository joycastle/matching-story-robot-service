package main

import (
	"flag"
	"io/ioutil"
	"strings"
)

func main() {
	f := flag.String("f", "", "文件路径")
	ff := flag.String("ff", "", "文件路径")
	flag.Parse()

	boms := []string{string('\uFEFF')}
	for _, bom := range boms {
		fs, err := ioutil.ReadFile(*f)
		if err != nil {
			panic(err)
		}
		s := string(fs)
		if strings.HasPrefix(s, bom) {
			s = strings.TrimPrefix(s, bom)
			if err := ioutil.WriteFile(*ff, []byte(s), 0666); err != nil {
				panic(err)
			}
		}
	}
}
