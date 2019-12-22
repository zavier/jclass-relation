package parse

import (
	"fmt"
	"github.com/zavier/jclass-relation/classfile"
	"strings"
)

type ClassInfo struct {
	ClassName   string
	IsInterface bool
	Level       int
	ParentClass *ClassInfo
	Interfaces  []*ClassInfo
}

func (cinfo *ClassInfo) String() string {
	return fmt.Sprintf("%s : level: %d", cinfo.ClassName, cinfo.Level)
}

func (cinfo *ClassInfo) AddLevel() {
	cinfo.Level = cinfo.Level + 1
	info := cinfo.ParentClass
	if info != nil {
		info.Level = cinfo.Level
	}
	infos := cinfo.Interfaces
	for _, i := range infos {
		if i != nil {
			i.Level = cinfo.Level
		}
	}
}

func CalcRelation(cp string, class string) *ClassInfo {
	cfList := parseAll(cp)
	classNameMap := make(map[string]*classfile.ClassFile)
	for _, v := range cfList {
		name := v.ClassName()
		classNameMap[name] = v
	}

	classname := strings.ReplaceAll(class, ".", "/")
	return wrapClass(classname, classNameMap)
}

func wrapClass(name string, cmap map[string]*classfile.ClassFile) *ClassInfo {
	cf, ok := cmap[name]
	if !ok {
		return nil
	}
	superClassName := cf.SuperClassName()
	parent := wrapClass(superClassName, cmap)

	interfaces := make([]*ClassInfo, 0)
	inames := cf.InterfaceNames()
	if inames != nil {
		for _, name := range inames {
			interfaces = append(interfaces, wrapClass(name, cmap))
		}
	}

	info := &ClassInfo{
		ClassName:   name,
		IsInterface: cf.IsInterface(),
		Level:       0,
		ParentClass: parent,
		Interfaces:  interfaces,
	}
	return info
}