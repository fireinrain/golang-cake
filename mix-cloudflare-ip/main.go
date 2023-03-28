package main

import (
	"archive/zip"
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/carlmjohnson/requests"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
	//result.cvs
	var ipPremium = "result.csv"
	exists := FileOrDirExists(ipPremium)
	if exists {
		currentTime := time.Now()
		format := currentTime.Format("2006-01-02 15:04:05")
		fmt.Println("Current date and time:", format)
		err := os.Rename(ipPremium, format+"-"+ipPremium)
		if err != nil {
			fmt.Printf("文件重命名失败: " + err.Error())
			return
		}
	}

	var cftestPath = "/Users/sunrise/BackSoftWares/CloudflareST_darwin_amd64/CloudflareST"
	var gitRepo = "https://ghproxy.com/https://github.com/ip-scanner/cloudflare/archive/refs/heads/daily.zip"
	var ipZipFile = "ip.zip"
	ctx := context.Background()
	err := requests.
		URL(gitRepo).
		ToFile(ipZipFile).
		Fetch(ctx)
	fmt.Printf("正在下载文件: %s,请稍后...\n", ipZipFile)
	if err != nil {
		fmt.Printf("文件下载失败: %s\n", err.Error())
		return
	}
	fmt.Printf("ip文件下载成功: %s\n", ipZipFile)
	fmt.Println("--------------------------------------------------------------------------")
	//解压
	reader, err := zip.OpenReader(ipZipFile)
	//自动释放
	defer reader.Close()
	if err != nil {
		fmt.Printf("文件: %s 解压缩失败: %s", ipZipFile, err.Error())
		return
	}

	var dstDir = "ip"
	if !FileOrDirExists(dstDir) {
		os.Mkdir(dstDir, os.ModeDir|os.ModePerm)
	}

	for _, file := range reader.File {
		//println(file.Name)
		if file.Name == "cloudflare-daily/" {
			continue
		}
		replaceStr := strings.Replace(file.Name, "cloudflare-daily", "", 1)

		filePath := filepath.Join(dstDir, replaceStr)
		fmt.Println("unzipping file ", filePath)
		//源文件
		fileInArchive, err := file.Open()
		if err != nil {
			panic(err)
		}
		//目标文件
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModeAppend|os.ModePerm)
		if err != nil {
			panic(err)
		}
		//源文件写入目标文件
		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}
		dstFile.Close()
		fileInArchive.Close()

		fmt.Printf("文件: %s 成功解压到: %s\n", file.Name, dstDir)
	}

	fmt.Println("--------------------------------------------------------------------------")

	//读取每一个文件的每一行 并判断是否是一个合格的ipv4格式
	//最后写入到ip.txt文件中

	var resultIPText = "ip.txt"
	if FileOrDirExists(resultIPText) {
		_ = os.Remove(resultIPText)
	}
	var lineCounter = 0
	resultFile, err := os.OpenFile(resultIPText, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	defer resultFile.Close()

	if err != nil {
		fmt.Printf("文件创建失败: %s\n", err.Error())
		return
	}

	dirEntries, err := os.ReadDir(dstDir)
	if err != nil {
		fmt.Printf("无法读取目录: %s,:%s", dstDir, err.Error())
		return
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		// is file type
		file, err := os.Open(filepath.Join(dstDir, entry.Name()))
		if err != nil {
			fmt.Printf("文件读取失败: %s", err)
			return
		}

		br := bufio.NewReader(file)
		for {
			line, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			//fmt.Println(string(line))
			//check if is ipv4
			isIpV4, err := checkIsIpV4(string(line))
			if err != nil {
				fmt.Errorf("ip格式错误跳过处理: %s", line)
				continue
			}
			if !isIpV4 {
				continue
			}
			//写入目标文件
			resultFile.WriteString(string(line) + "\n")
			lineCounter += 1
		}
		file.Close()
	}
	fmt.Printf("读取解析写入完成: %s,共有: %d个cloudflare ip", resultIPText, lineCounter)
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("正在运行优选ip程序,请稍后...")
	//获取ip.txt的绝对路径
	absPath, err := filepath.Abs(resultIPText)
	if err != nil {
		_ = fmt.Errorf(err.Error())
		return
	}
	//RunCloudflareST(cftestPath, absPath)
	cmd := cftestPath + " -f " + absPath
	fmt.Println("运行: " + cmd)
	RunWithCancelCommand(cmd)
	//cmd2 := "ping baidu.com"
	//RunWithCancelCommand(cmd2)

}

// RunCloudflareST run cloudflareSt soft to check ip that given by
func RunCloudflareST(cloudFStPath string, ipTextPath string) {
	cmd := exec.Command(cloudFStPath, "-f", ipTextPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() {
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				break
			}
			fmt.Print(readString)
		}
	}()
	_ = cmd.Run()
}

func RunWithCancelCommand(cmd string) {
	ctx, cancel := context.WithCancel(context.Background())
	go func(cancelFunc context.CancelFunc) {
		//超时时间
		time.Sleep(5 * 60 * time.Second)
		cancelFunc()
	}(cancel)
	c := exec.CommandContext(ctx, "bash", "-c", cmd)
	stdout, err := c.StdoutPipe()
	if err != nil {
		fmt.Printf("运行出错: %s", err.Error())
		return
	}
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			// 其实这段去掉程序也会正常运行，只是我们就不知道到底什么时候Command被停止了，而且如果我们需要实时给web端展示输出的话，这里可以作为依据 取消展示
			select {
			// 检测到ctx.Done()之后停止读取
			case <-ctx.Done():
				if ctx.Err() != nil {
					fmt.Printf("程序出现错误: %q", ctx.Err())
				} else {
					fmt.Println("程序被终止")
				}
				return
			default:
				readString, err := reader.ReadString('\n')
				if err != nil || err == io.EOF {
					break
				}
				fmt.Print(readString)
			}
		}
	}()
	_ = c.Run()
}

// check is a string is a valid ipv4 format
func checkIsIpV4(ipStr string) (bool, error) {
	ipv4Regex := `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`

	reg := regexp.MustCompile(ipv4Regex)
	//解释失败，返回nil
	if reg == nil {
		fmt.Println("regexp err")
		return false, errors.New("regexp err")
	}
	matchString := reg.MatchString(ipStr)
	return matchString, nil

}

func checkIfAIpv4(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.To4() != nil
}

// FileOrDirExists  判断所给路径文件/文件夹是否存在
func FileOrDirExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
