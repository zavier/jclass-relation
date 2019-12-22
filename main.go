package main

import (
	"fmt"
	"github.com/zavier/jclass-relation/parse"
	"strings"
)

func main() {
	cp := getClassPath()
	relation := parse.CalcRelation(cp, "org.springframework.beans.factory.support.DefaultListableBeanFactory")
	printClassInfo(relation)
}

func printClassInfo(cinfo *parse.ClassInfo) {
	if cinfo == nil {
		return
	}
	cinfo.AddLevel()

	level := cinfo.Level
	for i := 1; i < level; i++ {
		fmt.Print("\t")
	}
	fmt.Println(simpleName(cinfo.ClassName))
	printClassInfo(cinfo.ParentClass)
	for _, i := range cinfo.Interfaces {
		printClassInfo(i)
	}
}

func simpleName(className string) string {
	index := strings.LastIndex(className, "/")
	return className[index+1:]
}

func getClassPath() string {
	// 以maven 仓库中的 spring的部分源码测试
	version := "4.3.10.RELEASE"
	template := "/Users/zhengwei/.m2/repository/org/springframework/%s/" + version + "/%s-" + version + ".jar"
	springPackages := []string{"spring-core", "spring-aop", "spring-beans", "spring-context", "spring-context-support", "spring-core", "spring-test", "spring-tx", "spring-web", "spring-webmvc"}
	paths := make([]string, 0)
	for _, pkg := range springPackages {
		path := fmt.Sprintf(template, pkg, pkg)
		paths = append(paths, path)
	}
	cp := strings.Join(paths, ":")
	return cp
}
