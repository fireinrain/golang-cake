package variable

import (
	"fmt"
	_ "reflect"
	"testing"
)

func TestName(t *testing.T) {
	var age int
	fmt.Println(age)

	var age2 int = 100
	fmt.Println(age2)

	var age3 = 120
	fmt.Println(age3)
	fmt.Printf("%T", age2)
	fmt.Printf("%T", age3)

}

func TestIOAT(t *testing.T) {
	const PI = 3.1415

	const (
		Spring = iota + 1
		Summer
		Autumn
		Winter
	)

	fmt.Println(Spring)
}

func TestFunc(t *testing.T) {
	ab, cd := MaxMe(1, 2)
	fmt.Println(ab, cd)

}

func MaxMe(a int, b int) (ab int, cd int) {
	ab = 12
	cd = 18

	as := []int{1, 2, 3, 4}

	s := make([]int, 3)
	copy(s, as[1:])

	fmt.Printf("s value:%v", s)
	return 7, 0

}
func TestSlice2(t *testing.T) {

	as := []int{1, 2, 3, 4}

	s := make([]int, 3)
	copy(s, as[1:])

	m := make(map[int]string, 10)
	m[1] = "name"

	fmt.Printf("s value:%v", s)
	fmt.Printf("map is : %v", m)

}

type People interface {
	Name() string
}

type Man interface {
	Gender() string
}

type Xiaomin struct {
}

func (x Xiaomin) Name() string {
	fmt.Println("Name")
	return ""
}

func (x Xiaomin) Gender() string {
	fmt.Println("Gender")
	return ""
}

func TestInterface2(t *testing.T) {
	xiao := Xiaomin{}
	var p People
	p = xiao
	p.Name()
	var m Man
	m, ok := p.(Man)
	m.Gender()
	fmt.Println(ok)
	for i2 := range 15 {
		fmt.Println(i2)
	}

}

func TestExtends(t *testing.T) {
	//var bigDog BigDog = BigDog{
	//	&Dog{},
	//	&MidDog{},
	//}
	//
	//bigDog.Dog.Show()

	//必须要BigDog也实现Animal的所有方法
	var dog Animal
	dog = &BigDog{
		&Dog{},
		&MidDog{},
	}
	println(dog.Name())
}

type BigDog struct {
	*Dog
	*MidDog
}

func (m BigDog) Name() string {
	return "BigDog"
}

func (m BigDog) Show() {
	fmt.Println("big dog")
}

type MidDog struct {
}

func (m MidDog) Name() string {
	return "MidDog"
}

func (m MidDog) Show() {
	fmt.Println("mid dog")
}

type Dog struct {
}

func (d Dog) Name() string {
	return "Dog"
}

func (d Dog) Show() {
	fmt.Println("Dog")
}

type Animal interface {
	Name() string
	Show()
}
