package reflect

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"testing"
)

func TestReflect(t *testing.T) {
	var name = "xiaoqian"
	ShowInfo(name)

}

func ShowInfo(args interface{}) {
	fmt.Println("args Type: ", reflect.TypeOf(args))
	fmt.Println("args Value: ", reflect.ValueOf(args))

}

func TestReflect2(t *testing.T) {
	people := People{
		Name:    "xiyang",
		Age:     18,
		Address: "nanjing",
	}

	PrintInfo(people)

}
func PrintInfo(args any) {
	typeOf := reflect.TypeOf(args)
	fmt.Println("Type is: ", typeOf)
	valueOf := reflect.ValueOf(args)
	fmt.Println("Value is: ", valueOf)

	for index := range typeOf.NumField() {
		field := typeOf.Field(index)
		value := valueOf.Field(index)

		fmt.Println("index: ", index, ",field name: ", field.Name, ",value is: ", value)
	}

	invokeArgs := []reflect.Value{
		reflect.ValueOf("你好"),
	}
	for i := range typeOf.NumMethod() {
		method := typeOf.Method(i)
		fmt.Println("method name: ", method.Name)

		values := method.Func.Call(append([]reflect.Value{valueOf}, invokeArgs...))
		for _, va := range values {
			fmt.Println(va.Interface())

		}
	}

	//解析字段tag
	for i := range typeOf.NumField() {
		field := typeOf.Field(i)
		tagValue := field.Tag.Get("json")
		fmt.Println(tagValue)

		tagValue2 := field.Tag.Get("cityNum")
		if tagValue2 != "" {
			fmt.Println(tagValue2)

		}

	}

}

type People struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address" cityNum:"12"`
}

func (p People) ShowInfo(msg string) string {
	fmt.Printf("Name is: %s Age is: %d Address is: %s", p.Name, p.Age, p.Address)
	fmt.Println(msg)
	return "Nike" + msg
}

func TestChannel(t *testing.T) {
	var ch = make(chan int)
	go func() {
		fmt.Println("I send an int value ...")
		ch <- 100
		close(ch)
	}()

	value := <-ch
	fmt.Println("I receive a value: ", value)

}

func TestSlices(t *testing.T) {
	collect := slices.Collect(slices.Values([]int{1, 2, 3, 4}))
	fmt.Println(collect)

}

func TestEmptyInterface(t *testing.T) {
	i := NewEmpty()
	if i == nil {
		fmt.Println("E is nil")
	} else {
		fmt.Println("E is not nil")
	}
}

type Empty struct {
}

func NewEmpty() interface{} {
	var a *Empty = nil
	return a
}

func TestExtends(t *testing.T) {
	xiaoMin := XiaoMin{&Peopler{
		Name: "xiaoqian",
	}}
	xiaoMin.GetName()

	SayHello(xiaoMin)
}

type Peopler struct {
	Name string
}

func (p *Peopler) GetName() string {
	fmt.Println(p.Name)
	return p.Name
}

type IPeopler interface {
	GetName() string
}

func SayHello(p IPeopler) {
	fmt.Println("Hello ", p.GetName())
}

type XiaoMin struct {
	*Peopler
}

var (
	ErrNotFound    = errors.New("not found error")
	ErrArgsInvalid = errors.New("args is not valid")
	ErrServiceE    = errors.New("service invoke error")
)

func TestErrHandler(t *testing.T) {
	HandlerError()
}
func Handler1(age int) (int, error) {
	if age < 18 {
		//return -1, errors.New("handle1:check age -> cant smaller than 18")
		return -1, fmt.Errorf("handle1:check age -> %w %w", errors.New("cant smaller than 18"), ErrArgsInvalid)
	}
	return 0, nil
}

func Handler2() (int, error) {
	handler1, err := Handler1(17)
	if err != nil {
		return -1, fmt.Errorf("handler2:invoken h1 -> %w %w", err, ErrServiceE)
	}
	return handler1, nil
}

func Handle3() (int, error) {
	handler2, err := Handler2()
	if err != nil {
		return -1, fmt.Errorf("handler3:invoke h2 -> %w %w", err, ErrServiceE)
	}
	return handler2, nil
}

func HandlerError() {
	handle3, err := Handle3()
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, ErrServiceE) {
			fmt.Println("Error of service")
		}
	}
	fmt.Println(handle3)
}
