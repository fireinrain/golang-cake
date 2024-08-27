package main

import (
	"bytes"
	"fmt"
	"github.com/andygrunwald/go-trending"
	"io"
	"log"
	"log/slog"
	"net"
	"os/exec"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRuntime(t *testing.T) {

	go func(s string) {
		for i := 0; i < 2; i++ {
			fmt.Println(s)
		}
	}("world")
	// 主协程
	for i := 0; i < 2; i++ {
		// 切一下，再次分配任务
		runtime.Gosched()
		fmt.Println("hello")
	}
}

func TestRuntime2(t *testing.T) {
	go func() {
		defer fmt.Println("A.defer")
		func() {
			defer fmt.Println("B.defer")
			// 结束协程
			runtime.Goexit()
			defer fmt.Println("C.defer")
			fmt.Println("B")
		}()
		fmt.Println("A")
	}()
	for {
	}
}

func TestName(t *testing.T) {
	runtime.GOMAXPROCS(2)
	go a()
	go b()
	time.Sleep(10 * time.Second)

}

func TestChannel(t *testing.T) {
	var sChannel chan int = make(chan int)
	//var blockCh chan interface{} = make(chan interface{})
	go func() {
		for data := range sChannel {
			fmt.Println("I recieve a int: ", data)
		}

	}()

	go func() {
		for i := 0; i < 10; i++ {
			sChannel <- i
			fmt.Println("I have send a new value: ", i)
		}

		close(sChannel)
	}()

	//blockCh<-true

	select {}

}

func TestC(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- i
		}

		close(ch1)
	}()

	go func() {
		for {
			if value, ok := <-ch1; ok {
				ch2 <- value * 2
			} else {
				break
			}
		}
		close(ch2)
	}()

	for i2 := range ch2 {
		fmt.Println("ch2 recieve value: ", i2)
	}
}

func Test2(t *testing.T) {

	c := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			c <- i
		}
		close(c)
	}()
	for {
		if data, ok := <-c; ok {
			fmt.Println(data)
		} else {
			break
		}
	}
	fmt.Println("main结束")
}

func TestTimer(t *testing.T) {
	// 创建一个定时器，3秒后触发
	timer := time.NewTimer(3 * time.Second)

	go func() {
		time.Sleep(2 * time.Second)
		timer.Stop()
	}()

	select {
	case t := <-timer.C:
		fmt.Println("Timer fired!", t)
	case <-time.After(5 * time.Second):
		fmt.Println("Timer timeout")
		break

	}
	// 阻塞直到定时器触发
	//<-timer.C
	//fmt.Println("Timer fired!")

	<-time.After(10 * time.Second)
	fmt.Println("end for 10s")

}

func TestTimer3(t *testing.T) {
	timer := time.NewTimer(3 * time.Second)
	t2 := <-timer.C
	fmt.Printf("t2:%v\n", t2)
}

func TestTimer4(t *testing.T) {
	// 1.获取ticker对象
	ticker := time.NewTicker(1 * time.Second)
	i := 0
	// 子协程
	go func() {
		for {
			//<-ticker.C
			i++
			fmt.Println(<-ticker.C)
			if i == 5 {
				//停止
				ticker.Stop()
			}
		}
	}()
	for {
	}

}

func TestSync(t *testing.T) {
	var wg sync.WaitGroup
	var lock sync.Mutex
	var counter int = 10
	wg.Add(2)

	go func() {
		defer wg.Done()
		lock.Lock()
		counter += 1
		lock.Unlock()
	}()

	go func() {
		defer wg.Done()
		lock.Lock()
		counter += 1
		lock.Unlock()
	}()

	wg.Wait()
	fmt.Println("Result of counter: ", counter)

}

func TestAtomic(t *testing.T) {

	var counter int64 = 0
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()
	fmt.Println("Final counter: ", counter)

}

func TestTrending(t *testing.T) {
	trend := trending.NewTrending()
	projects, err := trend.GetProjects(trending.TimeToday, "go")
	if err != nil {
		panic(err)
	}
	for i, project := range projects {
		fmt.Println("index ", i, " ", project)
	}
}

func a() {
	for i := 1; i < 10; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("A:", i)
	}
}

func b() {
	for i := 1; i < 10; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("B:", i)
	}
}

func TestSync2(t *testing.T) {
	pool := sync.Pool{}
	pool.New = func() any {
		return "Hello"
	}

	pool.Put("nice")
	pool.Put("nice2")
	pool.Put("nice3")
	pool.Put("nice4")
	get := pool.Get()
	fmt.Println(get)
	get = pool.Get()
	fmt.Println(get)
	get = pool.Get()
	fmt.Println(get)
	get = pool.Get()
	fmt.Println(get)
}

func TestSyncMap(t *testing.T) {
	var syncMap sync.Map
	syncMap.Store("name", "xiaoqian")
	syncMap.Store("age", 12)
	if value, ok := syncMap.Load("name"); ok {
		fmt.Println(value)
	}

	syncMap.Range(func(key, value any) bool {
		fmt.Printf("key is: %s,value: %s", key, value)
		return true
	})

}

func TestTypeAssert(t *testing.T) {
	// 创建一个空接口切片，用于存储不同类型的切片
	var slices []interface{}

	// 添加整数切片
	intSlice := []int{1, 2, 3}
	slices = append(slices, intSlice)

	// 添加字符串切片
	strSlice := []string{"a", "b", "c"}
	slices = append(slices, strSlice)

	// 遍历并打印每个切片的内容
	for _, v := range slices {
		switch v := v.(type) {
		case []int:
			fmt.Println("Integer Slice:", v)
		case []string:
			fmt.Println("String Slice:", v)
		default:
			fmt.Println("Unknown Slice Type")
		}
	}

}

func TestReflect(t *testing.T) {
	// 定义一个 int 变量
	x := 10

	// 获取变量 x 的反射值
	v := reflect.ValueOf(&x) // 注意这里是传递的指针
	fmt.Println("Before:", x)

	// 通过反射修改值之前，必须确保它是可以修改的
	if v.Kind() == reflect.Ptr && !v.Elem().CanSet() {
		fmt.Println("Value cannot be modified")
		return
	}

	// 获取指针指向的值
	v = v.Elem()

	// 修改值
	v.SetInt(20)

	// 打印修改后的值
	fmt.Println("After:", x)

}

func TestT(t *testing.T) {
	// 需要反查的 IP 地址
	ip := "162.159.152.45"

	// 反查 DNS
	names, err := net.LookupAddr(ip)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 输出结果
	//for _, name := range names {
	//	fmt.Println("Hostname:", name)
	//}
	fmt.Printf("%v", names)
}

func PrintInfo(value any) {
	// 获取反射类型对象
	typ := reflect.TypeOf(value)
	// 获取反射值对象
	val := reflect.ValueOf(&value)

	fmt.Printf("Type: %v\n", typ)
	fmt.Printf("Value: %v\n", val)

	// 检查并展示值是否可写
	fmt.Printf("Value is settable: %v\n", val.CanSet())

	// 如果值是可写的，并且是int类型，尝试修改其值
	//if val.Kind() == reflect.Int && val.CanSet() {
	//	val.SetInt(val.Int() + 1)
	//	fmt.Printf("Modified Value: %v\n", val)
	//}

	if val.Kind() == reflect.Ptr && !val.Elem().CanSet() {
		fmt.Println("Value cannot be modified")
		return
	}

	val = val.Elem()
	val.SetInt(100)

	fmt.Println(val)

}

func TestNice(t *testing.T) {
	people := People{
		name: "xiaoqa",
	}
	defer people.print()
	people.name = "gaogao"
}

func TestOS(t *testing.T) {
	cmd1 := exec.Command("echo", "hello world")
	cmd2 := exec.Command("wc")

	reader, writer := io.Pipe()

	cmd1.Stdout = writer
	cmd2.Stdin = reader

	go func() {
		defer writer.Close()
		err := cmd1.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var buff bytes.Buffer
	cmd2.Stdout = &buff
	if err := cmd2.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("output of wc: %s", buff.String())
}

func TestSlog(t *testing.T) {
	slog.Info("Hello world")
	slog.Debug("debug is here")

	dog := Dog{}
	dog.Walk()
	println(dog.Name())

	var ages = []int{1, 2, 3, 4}

	ints := max(1, ages...)
	fmt.Println(ints)

	for i2 := range 10 {
		fmt.Println(i2)
	}

}

func TestMax(t *testing.T) {

	var ages = []int{1, 2, 3, 4}
	//golang的max内置函数定义了 不支持切片解构传参

	maxValue := max(1, ages...)
	fmt.Println(maxValue)
}

//func max(ints []int) int {
//	if len(ints) <= 0 {
//		panic("slice cant be empty")
//	}
//	temp := ints[0]
//	for _, i := range ints {
//		if i > temp {
//			temp = i
//		}
//	}
//
//	return temp
//}

type Dog struct {
}

func (d Dog) Name() string {
	return "Dog"
}

func (d Dog) Walk() {
	fmt.Println("Dog is walk")
}

type Animals interface {
	Name() string
	Walk()
}

type People struct {
	name string
}

func (p People) print() {
	fmt.Println(p.name)
}
