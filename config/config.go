package config

// todo 暂时没有使用
type config struct {
	configInfo map[string]string
}

const classPath string = "classPath"

var conf config = config{configInfo:make(map[string]string)}

func SetClassPath(cp string) {
	conf.configInfo[classPath] = cp
}

func GetClassPath() (string, error) {
	cp, ok := conf.configInfo[classPath]
	if ok {
		return cp, nil
	}
	return ".", nil
}
