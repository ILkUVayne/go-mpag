package common

import "os"

type Option func(st interface{})

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//IsNotExist来判断，是不是不存在的错误
	if os.IsNotExist(err) {
		//如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	//如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
	return false, err
}
