package main

import (
	"archive/zip"
	"better-cf-ip/cf"
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/carlmjohnson/requests"
	"io"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var cftestPath = "/Users/sunrise/BackSoftWares/CloudflareST_darwin_amd64/CloudflareST"
var gitRepo = "https://ghproxy.com/https://github.com/hello-earth/cloudflare-better-ip/archive/refs/heads/main.zip"
var speedTestUrl = "https://cloudflarest.gssmc.tk/100mb.zip"
var ipZipFile = "cf-ip.zip"

func main() {
	email := cf.CloudflareConfigValue.Email
	fmt.Println(email)
	//result.cvs
	_, failed := DownloadZipFile()
	if failed {
		return
	}
	fmt.Println("--------------------------------------------------------------------------")
	failed = UnzipIpFile()
	if failed {
		return
	}
	fmt.Println("--------------------------------------------------------------------------")

	resultIPText, notSuccess := ExtractAndCombineIp()
	if notSuccess {
		return
	}
	fmt.Println("--------------------------------------------------------------------------")
	CloudflareSpeedTest(resultIPText)
}

func CloudflareSpeedTest(resultIPText string) {
	fmt.Println("正在运行优选ip程序,请稍后...")
	//获取ip.txt的绝对路径
	absPath, err := filepath.Abs(resultIPText)
	if err != nil {
		_ = fmt.Errorf(err.Error())
	}
	//RunCloudflareST(cftestPath, absPath)
	cmdParams := []string{
		cftestPath,
		"-dn 20",
		"-p 20",
		"-url " + speedTestUrl,
		"-f " + absPath,
	}
	cmd := strings.Join(cmdParams, " ")
	fmt.Println("运行: " + cmd)
	RunWithCancelCommand(cmd)
	//cmd2 := "ping baidu.com"
	//RunWithCancelCommand(cmd2)
}

func ExtractAndCombineIp() (string, bool) {
	var resultIPText = "ip.txt"
	if FileOrDirExists(resultIPText) {
		_ = os.Remove(resultIPText)
	}
	var lineCounter = 0
	resultFile, err := os.OpenFile(resultIPText, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	defer resultFile.Close()

	if err != nil {
		fmt.Printf("文件创建失败: %s\n", err.Error())
		return "", true
	}

	var ipResults = []string{}

	//解析获得ip
	filepath.WalkDir("./ip/cloudflare-better-ip-main/cloudflare", func(path string, d fs.DirEntry, err error) error {
		if IsDir(path) {
			return nil
		}
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			fmt.Printf("文件读取失败: %s", err)
			return err
		}

		scanner := bufio.NewScanner(file)

		// Read line by line until the end of the file
		for scanner.Scan() {
			line := scanner.Text()
			split := strings.Split(line, "|")
			realIp := split[0]
			realIp = strings.TrimSpace(realIp)
			ipStr := strings.Split(realIp, ":")
			ip := ipStr[0]
			ipResults = append(ipResults, ip)
			//fmt.Println(line)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
			return err
		}

		return err
	})
	//remove duplicate lines
	cfIps := RemoveDuplicates(ipResults)
	//write to file
	for _, ip := range cfIps {
		lineCounter += 1
		resultFile.WriteString(string(ip) + "\n")
	}

	fmt.Printf("读取解析写入完成: %s,共有: %d个cloudflare ip \n", resultIPText, lineCounter)
	return resultIPText, false
}

func UnzipIpFile() bool {
	//解压
	reader, err := zip.OpenReader(ipZipFile)
	//自动释放
	defer reader.Close()
	if err != nil {
		fmt.Printf("文件: %s 解压缩失败: %s", ipZipFile, err.Error())
		return true
	}

	var dstDir = "ip"
	if !FileOrDirExists(dstDir) {
		os.Mkdir(dstDir, os.ModeDir|os.ModePerm)
	}
	err2 := Unzip(ipZipFile, "./"+dstDir)
	if err2 != nil {
		fmt.Println("Error unzipping:", err2)
		return true
	}
	fmt.Println("Unzipped ", ipZipFile, "successfully.")
	return false
}

func DownloadZipFile() (error, bool) {
	var ipPremium = "result.csv"
	exists := FileOrDirExists(ipPremium)
	if exists {
		currentTime := time.Now()
		format := currentTime.Format("2006-01-02 15:04:05")
		fmt.Println("Current date and time:", format)
		err := os.Rename(ipPremium, format+"-"+ipPremium)
		if err != nil {
			fmt.Printf("文件重命名失败: " + err.Error())
			return nil, true
		}
	}

	ctx := context.Background()
	err := requests.
		URL(gitRepo).
		ToFile(ipZipFile).
		Fetch(ctx)
	fmt.Printf("正在下载文件: %s,请稍后...\n", ipZipFile)
	if err != nil {
		fmt.Printf("文件下载失败: %s\n", err.Error())
		return nil, true
	}
	fmt.Printf("ip文件下载成功: %s\n", ipZipFile)
	return err, false
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

func Unzip(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(path, file.Mode())
			if err != nil {
				return err
			}
			continue
		}

		writer, err := os.Create(path)
		if err != nil {
			return err
		}
		defer writer.Close()

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		_, err = io.Copy(writer, fileReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func RemoveDuplicates(arr []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueArr []string

	for _, num := range arr {
		if !uniqueMap[num] {
			uniqueMap[num] = true
			uniqueArr = append(uniqueArr, num)
		}
	}
	return uniqueArr
}
