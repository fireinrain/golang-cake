package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	//RemoveEmptyDir()
	PrintEmptyDir()
}

func PrintEmptyDir() {
	var diskPath = "/Volumes/移动数据盘2/v2ph-images"
	fmt.Println("当前移动盘挂载路径: " + diskPath)
	separator := string(os.PathSeparator)

	dir, err := os.ReadDir(diskPath)
	if err != nil {
		fmt.Println("读取文件夹失败: " + err.Error())
		return
	}
	emptyDirCounter := 0
	for _, entry := range dir {
		if entry.IsDir() {
			//fmt.Println("找到目录: " + entry.Name())
			//判断是否是空目录
			currentDir := diskPath + separator + entry.Name()
			readDir, _ := os.ReadDir(currentDir)
			if len(readDir) <= 0 {
				fmt.Println("找到空目录: " + currentDir)
				emptyDirCounter += 1
				fromSlash := filepath.FromSlash(currentDir)
				fmt.Println(fromSlash)
				err := os.RemoveAll(fromSlash)
				if err != nil {
					fmt.Println("删除目录: " + currentDir + "失败")
					continue
				}
			}
		} else {
			fmt.Println("找到文件: " + entry.Name())
		}
	}
	fmt.Printf("一共找到：%d\n", emptyDirCounter)

}

func RemoveEmptyDir() {
	var diskPath = "/Volumes/移动数据盘2/文件"
	fmt.Println("当前移动盘挂载路径: " + diskPath)
	separator := string(os.PathSeparator)

	dir, err := os.ReadDir(diskPath)
	if err != nil {
		fmt.Println("读取文件夹失败: " + err.Error())
		return
	}
	emptyDirCounter := 0
	for _, entry := range dir {
		if entry.IsDir() {
			fmt.Println("找到目录: " + entry.Name())
			//判断是否是空目录
			currentDir := diskPath + separator + entry.Name()
			readDir, _ := os.ReadDir(currentDir)
			if len(readDir) <= 0 {
				fmt.Println("找到空目录: " + currentDir)
				emptyDirCounter += 1
				err := os.RemoveAll(currentDir)
				if err != nil {
					fmt.Println("删除目录: " + currentDir + "失败")
					return
				}
			}
		} else {
			fmt.Println("找到文件: " + entry.Name())
		}
	}
	fmt.Printf("成功删除: %d 个空目录\n", emptyDirCounter)

	//files, err := GetAllFiles(diskPath)
	//if err != nil {
	//	fmt.fmt.Println("读取文件夹失败: " + err.Error())
	//}
	//for _, file := range files {
	//	fmt.Println(file)
	//}
}

// GetFilesAndDirs 获取指定目录下的所有文件和目录
func GetFilesAndDirs(dirPth string) (files []string, dirs []string, err error) {
	dir, err := os.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetFilesAndDirs(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			//ok := strings.HasSuffix(fi.Name(), ".go")
			//if ok {
			files = append(files, dirPth+PthSep+fi.Name())
			//}
		}
	}

	return files, dirs, nil
}

// GetAllFiles 获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := os.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			//ok := strings.HasSuffix(fi.Name(), ".go")
			//if ok {
			files = append(files, dirPth+PthSep+fi.Name())
			//}
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}
