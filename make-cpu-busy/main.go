package main

import (
	crand "crypto/rand"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

func main() {
	//makeCpuBusy4()
	//makeCpuBusy3()
	go sendRandomData()
	go writeRandomData2Disk()
	makeCpuBusy3()
}

func writeRandomData(file *os.File, fileSize int64) error {
	_, err := io.CopyN(file, crand.Reader, fileSize)
	return err
}

func writeRandomData2Disk() {
	const fileSize = 15 * 1024 * 1024 * 1024 // 15GB
	for {
		file, err := os.Create("random_data.bin")
		if err != nil {
			fmt.Println("Error creating file:", err)
			file.Close()
			continue
		}

		if err := writeRandomData(file, fileSize); err != nil {
			fmt.Println("Error writing data:", err)
			file.Close()
			continue
		}
		fmt.Println("Data written successfully.")

		if err := os.Remove("random_data.bin"); err != nil {
			fmt.Println("Error removing file:", err)
			file.Close()
			continue
		}
		file.Close()
		fmt.Println("File removed successfully.")
	}

}

// 随机放松udp数据
func sendRandomData() {
	rand.NewSource(time.Now().UnixNano())
	data := make([]byte, 1024*9)
	rand.Read(data)

	for {
		randomIP := net.IPv4(byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)))
		randomPort := strconv.Itoa(rand.Intn(65535))

		addr, err := net.ResolveUDPAddr("udp", randomIP.String()+":"+randomPort)
		if err != nil {
			println("Error resolving UDP address: ", err.Error())
			continue
			//return
		}

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			println("Error dialing UDP: ", err.Error())
			continue
			//return
		}
		defer conn.Close()

		_, err = conn.Write(data)
		if err != nil {
			println("Error writing to UDP: ", err.Error())
			continue
			//return
		}

		time.Sleep(time.Millisecond * 100)
	}
}

// makeCpuBusy6
//
//	@Description:
func makeCpuBusy6() {
	rand.NewSource(time.Now().UnixNano())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				math.Sin(rand.Float64())
			}
		}()
	}

	select {}
}

// makeCpuBusy5
//
//	@Description:
func makeCpuBusy5() {
	rand.NewSource(time.Now().UnixNano())

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				math.Sqrt(float64(rand.Int63()))
			}
		}()
	}

	select {}
}

// makeCpuBusy4
//
//	@Description:
func makeCpuBusy4() {
	rand.NewSource(time.Now().UnixNano())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			var x, y int
			for i := 0; i < 100000000; i++ {
				x, y = rand.Int(), rand.Int()
				_ = x * y
			}
		}()
	}

	select {}
}

// makeCpuBusy3
//
//	@Description:
func makeCpuBusy3() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				runtime.Gosched()
			}
		}()
	}
	for {
		time.Sleep(time.Second)
	}
}

// makeCpuBusy1
//
//	@Description:
func makeCpuBusy1() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				// Spin lock
			}
		}()
	}

	select {}
}

// makeCpuBusy2
//
//	@Description:
func makeCpuBusy2() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			var x, y int
			for i := 0; i < 100000000; i++ {
				x, y = rand.Int(), rand.Int()
				_ = x * y
			}
		}()
	}
	select {}
}
