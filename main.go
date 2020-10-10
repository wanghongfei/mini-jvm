package main

import (
	"flag"
	"fmt"
	"github.com/wanghongfei/mini-jvm/utils"
	"github.com/wanghongfei/mini-jvm/vm"
	"os"
	"strings"
)

func main() {
	// 命令行参数
	mainClass := flag.String("main", "", "主类全名")
	classpath := flag.String("classpath", "", "类路径,可以是目录也可以是jar包路径, 多个用逗号分隔")
	consoleLog := flag.Bool("consoleLog", false, "是否在控制台打印JVM日志")
	flag.Parse()

	if "" == *mainClass {
		fmt.Println("error: lack main class")
		os.Exit(1)
	}

	// 初始化日志
	utils.InitLog(*consoleLog)

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
		utils.LogErrorPrintf("%+v", err)
		os.Exit(1)
	}
	utils.LogInfoPrintf("JVM instance created")

	err = miniJvm.Start()
	if nil != err {
		utils.LogErrorPrintf("%+v", err)
		os.Exit(1)
	}
}
