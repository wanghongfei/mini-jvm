package main

import (
	"github.com/wanghongfei/mini-jvm/vm"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	miniJvm, err := vm.NewMiniJvm("Hello", []string{"out/"})
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
