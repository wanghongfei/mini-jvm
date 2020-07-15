package main

import (
	"flag"
	"github.com/wanghongfei/mini-jvm/vm"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetOutput(os.Stdout)

	mainClass := flag.String("main", "", "主类全名")
	classpath := flag.String("classpath", "", "类路径,可以是目录也可以是jar包路径, 多个用逗号分隔")
	flag.Parse()

	if "" == *mainClass {
		log.Printf("lack main class")
		os.Exit(1)
	}

	var path []string
	if "" != *classpath {
		path = strings.Split(*classpath, ",")
	} else {
		path = []string {*classpath}
	}

	cmdArgs := flag.Args()

	miniJvm, err := vm.NewMiniJvm(*mainClass, path, cmdArgs...)
	if nil != err {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	err = miniJvm.Start()
	if nil != err {
		log.Printf("%+v", err)
		os.Exit(1)
	}
}
