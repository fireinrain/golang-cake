package main

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"unsafe"
)

func main() {
	var a float64 = 12.3434342342
	b := float32(a)
	fmt.Println(b)
	var name string = "你好"
	size := unsafe.Sizeof(name)
	fmt.Printf("name occupy %d bytes\n", size)

	var arr0 [5]int = [5]int{1, 2, 3, 4, 5}
	fmt.Printf("%v\n", arr0)
	var arr1 = [5]string{"n", "a", "b", "c", "k"}
	fmt.Printf("%v", arr1)
	var arr2 = [...]int{1, 3, 5, 7, 9}
	fmt.Printf("%v", arr2)

	d := [2]struct {
		name string
		age  int
	}{
		{"nihao", 12},
		{"wohao", 16},
	}

	fmt.Printf("%v\n", d)

	//多维数组
	var multi_array = [3][2]int{{1, 2}, {2, 3}, {3, 4}}
	fmt.Printf("%v", multi_array)

	var nu [2]int = [2]int{1, 2}
	fmt.Printf("%v\n", nu)
	test(nu)
	fmt.Printf("%v\n", nu)

	var nu2 = &nu
	test2(nu2)
	fmt.Printf("%v\n", nu2)

	var sum = [3]int{3, 4, 5}
	array := sunArray(sum)
	fmt.Printf("Array sum is: %v", array)

	s := []int{1, 2, 3}
	testSlice(s)
	fmt.Printf("%v\n", s)

	//test copy
	ints := []int{1, 2, 3}
	newInts := make([]int, 3, 3)
	fmt.Printf("before: %v\n", newInts)
	copy(newInts, ints)
	fmt.Printf("after: %v\n", newInts)

	var p *string
	fmt.Println(p)
	fmt.Printf("%s", p)
	if p == nil {
		fmt.Println("空指针")
	} else {
		fmt.Println("非空指针")
	}

	var aa *int
	aa = new(int)
	*aa = 100
	fmt.Println(*aa)
	fmt.Println("%p", aa)

	//var bb map[string]int
	//bb["测试"] = 100
	//fmt.Println(bb)

	scoreMap := make(map[string]int, 8)
	scoreMap["张三"] = 90
	scoreMap["小明"] = 100
	fmt.Println(scoreMap)
	fmt.Println(len(scoreMap))

	var numm *int
	fmt.Printf("%p\n", numm)

	var hashMap = make(map[string]int, 3)
	hashMap["xiaoqian"] = 12
	hashMap["baibai"] = 13

	if _, ok := hashMap["xiaoqian"]; ok {
		fmt.Println("xiaoqian in map")
	}
	delete(hashMap, "xiaoqian")
	fmt.Printf("%v\n", hashMap)

	type MyInt = int
	var age MyInt = 12
	var age2 int = 12
	if age == age2 {
		fmt.Println("type equal")
	}

	type MyInt2 int
	var peek MyInt2 = 12
	var size2 int = 12
	if int(peek) == size2 {
		fmt.Println("type not equal")
	}

	type student2 struct {
		name string
		age  int
	}
	m := make(map[string]*student2)
	stus := []student2{
		{name: "pprof.cn", age: 18},
		{name: "测试", age: 23},
		{name: "博客", age: 28},
	}

	for _, stu := range stus {
		m[stu.name] = &stu
	}
	for k, v := range m {
		fmt.Println(k, "=>", *v)
	}

	dog := Dog{
		Head: 12,
		Age:  8,
		Animal: &Animal{
			name: "乐乐",
		},
	}
	dog.move()

	var ce []student //定义一个切片类型的结构体
	ce = []student{
		student{1, "xiaoming", 22},
		student{2, "xiaozhang", 33},
	}
	fmt.Println(ce)
	demo(ce)
	fmt.Println(ce)

	var c1, c2, c3 chan int
	var i1, i2 int
	select {
	case i1 = <-c1:
		fmt.Printf("received ", i1, " from c1\n")
	case c2 <- i2:
		fmt.Printf("sent ", i2, " to c2\n")
	case i3, ok := (<-c3): // same as: i3, ok := <-c3
		if ok {
			fmt.Printf("received ", i3, " from c3\n")
		} else {
			fmt.Printf("c3 is closed\n")
		}
	default:
		fmt.Printf("no communication\n")
	}

	//for {
	//	fmt.Println("nonono")
	//}

	var i int = 0

	var upper = func(name string) string {
		return strings.ToUpper(name)
	}

	s2 := upper("xiaoqian")
	fmt.Println(s2)

	//var whatever [5]struct{}
	//for i := range whatever {
	//	var ii = i
	//	defer func(index int) { fmt.Println(ii) }(ii)
	//}

	var whatever [5]struct{}
	for i := range whatever {
		defer func(index int) { fmt.Println(i) }(i)
	}

start: // 这是一个标签
	fmt.Println(i)
	i++
	if i < 5 {
		goto start // 跳转回 start 标签
	}
	fmt.Println("End")

	fmt.Println("--------")
	stru := NewCustomStruct(5, WithAddValue(-1), WithAddValue(2), WithMinusValue(2))
	fmt.Println(stru.A)

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

func Test5(t *testing.T) {
	people := People{
		name: "xiaoqa",
	}
	defer people.print()
	people.name = "gaogao"
}

type People struct {
	name string
}

func (p People) print() {
	fmt.Println(p.name)
}

type T struct {
}

func (t T) a() {

}
func (t T) b() {

}

type A interface {
	a()
	b()
}

type B interface {
	a()
	b()
	c()
}

type CustomStruct struct {
	A int
}

type Options func(*CustomStruct)

func WithAddValue(value int) Options {
	return func(c *CustomStruct) {
		c.A += value
	}
}
func WithMinusValue(value int) Options {
	return func(c *CustomStruct) {
		c.A -= value
	}
}
func NewCustomStruct(initValue int, opts ...Options) *CustomStruct {
	res := &CustomStruct{
		A: initValue,
	}

	for _, opt := range opts {
		opt(res)
	}
	return res
}

type student struct {
	id   int
	name string
	age  int
}

func demo(ce []student) {
	//切片是引用传递，是可以改变值的
	ce[1].age = 999
	// ce = append(ce, student{3, "xiaowang", 56})
	// return ce
}

type Animal struct {
	name string
}

func (a *Animal) move() {
	fmt.Println("%s move", a.name)
}

type Dog struct {
	Head int
	Age  int
	*Animal
}

func testSlice(a []int) {
	fmt.Printf("%v\n", a)
	a[0] = 12
}

func test(num [2]int) {
	fmt.Printf("%v", num)
	num[0] = 2
}

func test2(nump *[2]int) {
	fmt.Printf("%v", nump)
	nump[0] = 2
}

func sunArray(a [3]int) int {
	result := 0
	for _, i2 := range a {
		result += i2
	}
	return result
}

func Bigger(a any, b any) {

}
func CreatUser(name, address string, age int) {

}
func CreateUser2(age int, name, adress string) {

}

func init() {
	fmt.Println("init1")
}

func init() {
	fmt.Println("init2")
}

func init() {
	fmt.Println("init3")
}
