package main

import (
	"github.com/wanghongfei/mini-jvm/vm"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	if len(os.Args) < 3 {
		log.Printf("usage: [main-class-name] [classpath]")
		os.Exit(1)
	}

	miniJvm, err := vm.NewMiniJvm(os.Args[1], os.Args[2:])
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
