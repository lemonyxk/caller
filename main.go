/**
* @program: caller
*
* @description:
*
* @author: lemon
*
* @create: 2021-07-01 22:46
**/

package caller

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
)

type Info struct {
	Line   int
	File   string
	Func   string
	Module string
}

func (i Info) String() string {
	return fmt.Sprintf("%s:%d %s", i.File, i.Line, i.Func)
}

var pwd = ""
var mux sync.RWMutex
var cache = make(map[string]Info)

func init() {
	_, codePath, _, _ := runtime.Caller(0)
	pwd = filepath.Dir(codePath)
}

func Deep(deep int) Info {

	var file, line = "", 0
	pc, codePath, codeLine, ok := runtime.Caller(deep)
	if !ok {
		return Info{}
	}

	file, line = codePath, codeLine
	var key = file + ":" + strconv.Itoa(line)

	mux.RLock()
	if info, ok := cache[key]; ok {
		mux.RUnlock()
		return info
	}
	mux.RUnlock()

	var m, f, l = clipFileAndLine(file, line)

	var info = Info{
		Line:   l,
		File:   f,
		Func:   filepath.Base(runtime.FuncForPC(pc).Name()),
		Module: m,
	}

	mux.Lock()
	cache[key] = info
	mux.Unlock()

	return info
}

func Deeps(deep int) []Info {
	var res []Info
	for skip := deep; true; skip++ {
		var info = Deep(skip)
		if info.File == "" {
			break
		}
		res = append(res, info)
	}
	return res
}

func Stack(deep int) (string, string, int) {
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

func clipFileAndLine(file string, line int) (string, string, int) {
	if file == "" || line == 0 {
		return "", "", 0
	}

	switch runtime.GOOS {
	case "windows":
		pwd = strings.Replace(pwd, "\\", "/", -1)
	}

	var i = getCommonStr(file, pwd)

	var m = "main"

	if i <= 1 {
		return m, file, line
	}

	file = file[i:]
	if file[0] == '/' {
		file = file[1:]
	}

	var arr = strings.Split(file, "/")
	if len(arr) == 1 {
		return m, file, line
	}

	if strings.Contains(arr[0], "@v") {
		arr = arr[1:]
		m = strings.Split(arr[0], "@")[0]
	} else {
		m = arr[0]
	}

	if len(arr) == 1 {
		return m, arr[0], line
	}

	if len(arr) == 2 {
		return m, strings.Join(arr, "/"), line
	}

	return arr[1], strings.Join(arr[2:], "/"), line
}

func getCommonStr(str string, str1 string) int {
	var i = 0
	for ; i < len(str) && i < len(str1); i++ {
		if str[i] != str1[i] {
			break
		}
	}
	return i
}
