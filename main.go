/**
* @program: caller
*
* @description:
*
* @author: lemo
*
* @create: 2021-07-01 22:46
**/

package caller

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

type Info struct {
	Line int
	File string
	Func string
}

func (i Info) String() string {
	return fmt.Sprintf("%s:%d %s", i.File, i.Line, i.Func)
}

var rootPath, _ = os.Getwd()

func Auto(packageName string) (string, int) {

	var file, line = "", 0

	for skip := 1; true; skip++ {
		pc, codePath, codeLine, ok := runtime.Caller(skip)
		if !ok {
			break
		}

		prevFunc := runtime.FuncForPC(pc).Name()

		if !strings.Contains(prevFunc, packageName) {
			file, line = codePath, codeLine
			break
		}
	}

	return clipFileAndLine(file, line)
}

func Deep(deep int) Info {

	var file, line = "", 0
	pc, codePath, codeLine, ok := runtime.Caller(deep)
	if !ok {
		return Info{}
	}

	file, line = codePath, codeLine

	var f, l = clipFileAndLine(file, line)

	var info = Info{
		Line: l,
		File: f,
		Func: filepath.Base(runtime.FuncForPC(pc).Name()),
	}

	return info
}

func Stack(deep int) (string, int) {
	var list = strings.Split(string(debug.Stack()), "\n")
	var info = strings.TrimSpace(list[deep])
	var flInfo = strings.Split(strings.Split(info, " ")[0], ":")
	var file, l = flInfo[0], flInfo[1]
	var line, _ = strconv.Atoi(l)
	return clipFileAndLine(file, line)
}

func GetFuncName(fn interface{}) string {
	t := reflect.ValueOf(fn).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	}
	return t.String()
}

func FuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func clipFileAndLine(file string, line int) (string, int) {
	if file == "" || line == 0 {
		return "", 0
	}

	if runtime.GOOS == "windows" {
		rootPath = strings.Replace(rootPath, "\\", "/", -1)
	}

	if rootPath == "/" {
		return file, line
	}

	if strings.HasPrefix(file, rootPath) {
		file = file[len(rootPath)+1:]
	}

	return file, line
}
