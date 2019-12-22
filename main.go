package main

import (
	"fmt"
	"github.com/zavier/jclass-relation/parse"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	// 这里需要使用者修改成自己要使用的classpath
	cp := getClassPath()
	// 这里需要指定要查找的类（全称）
	relation := parse.CalcRelation(cp, "org.springframework.beans.factory.support.DefaultListableBeanFactory")
	print2Graphviz(relation)
}

// todo 之后可以使用接口，实现不同方式的输出
func print2Graphviz(cinfo *parse.ClassInfo) {
	if cinfo == nil {
		log.Println("no find relation")
		return
	}
	// 去重使用
	strSet := make(map[string]interface{})

	graphviz := "digraph classRelation {\n"
	str := doPrint2Graphviz(cinfo)
	for _, v := range str {
		_, ok := strSet[v]
		if ok {
			continue
		}
		strSet[v] = nil
		graphviz += v + ";\n"
	}
	graphviz += "}\n"
	// 按照graphviz，  mac: brew install graphviz
	// 生成dot文件后， 执行： dot -Tpng classRelation.dot -o classRelation.png
	_ = ioutil.WriteFile("classRelation.dot", []byte(graphviz), 777)
}

// 输出关系
func doPrint2Graphviz(cinfo *parse.ClassInfo) []string {
	if cinfo == nil {
		return nil
	}
	cinfo.AddLevel()

	str := make([]string, 0)
	name := simpleName(cinfo.ClassName)
	parentClass := cinfo.ParentClass
	if parentClass != nil {
		str = append(str, name+"->"+simpleName(parentClass.ClassName))
	}
	infos := cinfo.Interfaces
	if infos != nil && len(infos) > 0 {
		for _, i := range infos {
			str = append(str, name+"->"+simpleName(i.ClassName))
		}
	}

	// 继续输出
	graphviz := doPrint2Graphviz(parentClass)
	if graphviz != nil && len(graphviz) > 0 {
		str = append(str, graphviz...)
	}
	for _, i := range infos {
		i2 := doPrint2Graphviz(i)
		if i2 != nil && len(i2) > 0 {
			str = append(str, i2...)
		}
	}
	return str
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
