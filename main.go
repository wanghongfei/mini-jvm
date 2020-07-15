package main

import (
	"flag"
	"github.com/wanghongfei/mini-jvm/vm"
	"io"
	"log"
	"os"
	"strings"
)

var vmErrorLog io.Writer

func main() {
	// 命令行参数
	mainClass := flag.String("main", "", "主类全名")
	classpath := flag.String("classpath", "", "类路径,可以是目录也可以是jar包路径, 多个用逗号分隔")
	flag.Parse()

	if "" == *mainClass {
		log.SetOutput(os.Stdout)
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

	// 启动jvm
	miniJvm, err := vm.NewMiniJvm(*mainClass, path, cmdArgs...)
	if nil != err {
		log.SetOutput(createVmErrorLog())
		log.Printf("%+v", err)
		os.Exit(1)
	}

	err = miniJvm.Start()
	if nil != err {
		log.SetOutput(createVmErrorLog())
		log.Printf("%+v", err)
		os.Exit(1)
	}
}

func createVmErrorLog() io.Writer {
	if nil != vmErrorLog {
		return vmErrorLog
	}

	logFile, err := os.OpenFile("vm-error.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0766)
	if nil != err {
		panic(err)
	}

	vmErrorLog = logFile

	return logFile
}
