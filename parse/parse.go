package parse

import (
	"fmt"
	"github.com/zavier/jclass-relation/classfile"
	"github.com/zavier/jclass-relation/classpath"
	"log"
)

func parseAll(cp string) []*classfile.ClassFile {
	log.Printf("classpath: %s\n", cp)
	load := classpath.LoadAll(cp)
	total := make([][]byte, 0)
	for i := 0; i < len(load); i++ {
		total = append(total, load[i])
	}
	fmt.Printf("total: %d\n", len(total))
	cfList := make([]*classfile.ClassFile, 0)
	for i := 0; i < len(total); i++ {
		cf, err := classfile.Parse(total[i])
		if err != nil {
			log.Printf("parse error: %s\n", err.Error())
		} else {
			cfList = append(cfList, cf)
		}
	}
	return cfList
}
