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
	"os"
	"runtime"
	"strings"
)

var rootPath, _ = os.Getwd()

var packageName = "github.com/lemoyxk/caller"

func Auto() (string, int) {

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

func Deep(deep int) (string, int) {

	var file, line = "", 0
	_, codePath, codeLine, ok := runtime.Caller(deep)
	if !ok {
		return file, line
	}

	file, line = codePath, codeLine

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
