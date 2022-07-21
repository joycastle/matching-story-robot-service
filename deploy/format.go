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

	fs, err := ioutil.ReadFile(*f)
	if err != nil {
		panic(err)
	}
	outs := string(fs)

	for _, bom := range boms {
		if strings.HasPrefix(outs, bom) {
			outs = strings.TrimPrefix(outs, bom)
			break
		}
	}

	if err := ioutil.WriteFile(*ff, []byte(outs), 0666); err != nil {
		panic(err)
	}
}
