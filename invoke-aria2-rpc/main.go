package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const Aria2RPCEndpoint = "http://216.127.164.234:6800/jsonrpc"
const DownloadDirectory = "/root/downloads"

var aria2Token = os.Getenv("Aria2Token")

// ReadMagnetUrl
//
//	@Description: 读取magnet 文件
//	@param path2CsvFile
//	@return []string
func ReadMagnetUrl(path2CsvFile string) []string {
	csvFile, err := os.Open(path2CsvFile)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	result := []string{}

	for _, record := range records {
		s := record[0]
		mag := record[1]
		result = append(result, s+"-"+mag)
	}
	return result
}

type Aria2Request struct {
	JsonRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      string      `json:"id"`
}
type Aria2Response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}

// InvokeAria2ForMagnetUrl
//
//	@Description: aria2 添加magnet任务
//	@param magnetUrl
//	@param rpcUrl
//	@param ariaToken
//	@param downloadPath
//	@return error
func InvokeAria2ForMagnetUrl(magnetUrl string, rpcUrl string, ariaToken string, downloadPath string) (Aria2Response, error) {
	var respJson Aria2Response

	magnetURL := magnetUrl
	downloadDir := downloadPath
	rpcEndpoint := rpcUrl
	jsonReqBody := Aria2Request{
		JsonRPC: "2.0",
		Method:  "aria2.addUri",
		Params: []interface{}{
			"token:" + ariaToken,
			[]string{magnetURL},
			map[string]string{
				"dir": downloadDir,
			},
		},
		ID: "1",
	}

	reqBody, err := json.Marshal(jsonReqBody)
	if err != nil {
		return respJson, err
	}

	resp, err := http.Post(rpcEndpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return respJson, err
	}
	if resp.StatusCode != 200 {
		return respJson, errors.New("Invoke error with request: " + string(resp.StatusCode))
	}
	err = json.NewDecoder(resp.Body).Decode(&respJson)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Printf("Magnet URL %s added to Aria2\n", magnetURL)
	return respJson, nil
}

type Aria2ActiveResponse struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  []struct {
		CompletedLength string `json:"completedLength"`
		Connections     string `json:"connections"`
		DownloadSpeed   string `json:"downloadSpeed"`
		Gid             string `json:"gid"`
		NumSeeders      string `json:"numSeeders"`
		Seeder          string `json:"seeder"`
		Status          string `json:"status"`
		TotalLength     string `json:"totalLength"`
		UploadSpeed     string `json:"uploadSpeed"`
	} `json:"result"`
}

// GetAllActiveTaskInfo
//
//	@Description: 获取所有激活的任务
//	@param rpcUrl
//	@param ariaToken
func GetAllActiveTaskInfo(rpcUrl string, ariaToken string) Aria2ActiveResponse {
	rpcEndpoint := rpcUrl

	jsonReqBody := Aria2Request{
		JsonRPC: "2.0",
		Method:  "aria2.tellActive",
		Params: []interface{}{
			"token:" + ariaToken,
			[]string{
				"gid",
				"totalLength",
				"completedLength",
				"uploadSpeed",
				"downloadSpeed",
				"connections",
				"numSeeders",
				"seeder",
				"status",
				"errorCode",
				"verifiedLength",
				"verifyIntegrityPending",
			},
		},
		ID: "1",
	}
	reqBody, err := json.Marshal(jsonReqBody)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(rpcEndpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result Aria2ActiveResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic(err)
	}
	return result
}

type Aria2StatusResponse struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  Result `json:"result"`
}
type Bittorrent struct {
	AnnounceList [][]string `json:"announceList"`
}
type Files struct {
	CompletedLength string        `json:"completedLength"`
	Index           string        `json:"index"`
	Length          string        `json:"length"`
	Path            string        `json:"path"`
	Selected        string        `json:"selected"`
	Uris            []interface{} `json:"uris"`
}
type Result struct {
	Bitfield        string     `json:"bitfield"`
	Bittorrent      Bittorrent `json:"bittorrent"`
	CompletedLength string     `json:"completedLength"`
	Connections     string     `json:"connections"`
	Dir             string     `json:"dir"`
	DownloadSpeed   string     `json:"downloadSpeed"`
	ErrorCode       string     `json:"errorCode"`
	ErrorMessage    string     `json:"errorMessage"`
	Files           []Files    `json:"files"`
	FollowedBy      []string   `json:"followedBy"`
	Gid             string     `json:"gid"`
	InfoHash        string     `json:"infoHash"`
	NumPieces       string     `json:"numPieces"`
	NumSeeders      string     `json:"numSeeders"`
	PieceLength     string     `json:"pieceLength"`
	Status          string     `json:"status"`
	TotalLength     string     `json:"totalLength"`
	UploadLength    string     `json:"uploadLength"`
	UploadSpeed     string     `json:"uploadSpeed"`
}

// CheckAriaTaskStatus
//
//	@Description: 根据GID 获取任务的完成状态
//	@param id
//	@param rpcUrl
//	@param ariaToken
//	@return Aria2StatusResponse
//	@return bool
func CheckAriaTaskStatus(id string, rpcUrl string, ariaToken string) (Aria2StatusResponse, bool) {
	gid := id
	rpcEndpoint := rpcUrl

	jsonReqBody := Aria2Request{
		JsonRPC: "2.0",
		Method:  "aria2.tellStatus",
		Params: []interface{}{
			"token:" + ariaToken,
			gid,
		},
		ID: "1",
	}
	reqBody, err := json.Marshal(jsonReqBody)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(rpcEndpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result Aria2StatusResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Status: %s\n", result.Result.Status)
	//fmt.Printf("Total length: %s bytes\n", result.Result.TotalLength)
	//fmt.Printf("Completed length: %s bytes\n", result.Result.CompletedLength)
	if result.Result.Status == "complete" {
		fmt.Println("Task is complete")
		return result, true
	} else {
		fmt.Printf("Task is not yet complete. Completed percent: %s bytes\n", result.Result.CompletedLength)
		return result, false
	}
	return result, false
}

// CreateFolderInOnedrive
//
//	@Description: 在onedirve上创建目录
func CreateFolderInOnedrive() {

}

// UploadDir2OneDrive
//
//	@Description: 将目录上传到onedirve
func UploadDir2OneDrive() {

}

// GetAllDirFiles
//
//	@Description: 获取该目录所有文件
//	@param dirPath
//	@return []string
func GetAllDirFiles(dirPath string) []string {
	return []string{}
}

// IsMetaData
//
//	@Description: 判断下载的是否为磁力元数据
//	@param fileLength
//	@return bool
func IsMetaData(fileLength string) bool {
	metaDataSize := 50 * 1024 * 1024
	atoi, err := strconv.Atoi(fileLength)
	if err != nil {
		fmt.Println("格式化数字错误: ", err.Error())
		panic(err)
	}
	if metaDataSize > atoi {
		return true
	}
	return false
}

func main() {

	magnetUrl := ReadMagnetUrl("./麻豆1.26.csv")
	for _, value := range magnetUrl {
		fmt.Println("正在处理链接: ", value)
		infoSlice := strings.Split(value, "-")
		downloadDir := DownloadDirectory + string(os.PathSeparator) + infoSlice[0]
		resp, err := InvokeAria2ForMagnetUrl(infoSlice[1], Aria2RPCEndpoint, aria2Token, downloadDir)
		if err != nil {
			fmt.Errorf("invoke aria2 error: %s", err.Error())
		}
		fmt.Printf("id: %s \n", resp.Id)
		fmt.Printf("result: %s\n", resp.Result)
		fmt.Println("---------------------------------------")

	checkLoop:
		for {
			//获取所有活动的任务
			taskInfo := GetAllActiveTaskInfo(Aria2RPCEndpoint, aria2Token)
			tasks := taskInfo.Result
			//没有活动的任务就跳出检查
			if len(tasks) <= 0 {
				break checkLoop
			}
			for _, task := range tasks {
				length := task.TotalLength
				//check length
				isMetaData := IsMetaData(length)
				if isMetaData {
					continue
				}
				//获取GID
				gid := task.Gid
				fmt.Println("下载任务GID: ", gid)
				result, status := CheckAriaTaskStatus(gid, Aria2RPCEndpoint, aria2Token)
				if !status {
					time.Sleep(5 * time.Second)
				} else {
					//下载完成
					fmt.Println("任务下载完成,正在上传onedrive...")

					fmt.Println(result.Result.Dir)

					break checkLoop
				}

			}

		}

		//上传onedrive
		fmt.Println("上传onedrive完成...")

		//删除源文件
		fmt.Println("删除源文件成功...")

		break
	}

	//for {
	//	id := "a3845f0338ed9c7a"
	//	_, status := CheckAriaTaskStatus(id, Aria2RPCEndpoint, aria2Token)
	//	if status {
	//		break
	//	}
	//	time.Sleep(1 * time.Second)
	//
	//}
	//
	//for {
	//	result := GetAllActiveTaskInfo(Aria2RPCEndpoint, aria2Token)
	//	fmt.Printf("result: %s", result)
	//	time.Sleep(1 * time.Second)
	//}

}
