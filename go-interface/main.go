package main

func main() {
	var p Phone
	p = &XiaoMiPhone{}
	//p是小米指针类型 如果要类型转换的话 必须要赚化为同为指针类型
	phone := p.(*XiaoMiPhone)
	phone.MIUI()
}

type Phone interface {
	call()
	text()
	wifi()
}

type XiaoMiPhone struct {
}

func (x XiaoMiPhone) call() {
	//TODO implement me
	panic("implement me")
}

func (x XiaoMiPhone) text() {
	//TODO implement me
	panic("implement me")
}

func (x XiaoMiPhone) wifi() {
	//TODO implement me
	panic("implement me")
}

func (x XiaoMiPhone) MIUI() {
	println("xiaomu UI")
}

type Readme interface {
	Read()
	ReadAll()
}

type WriteMe interface {
	Write()
	WriteAll()
}

type IOBound interface {
	Readme
	WriteMe
	Close()
}

type USB struct {
}

func (U USB) Read() {
	//TODO implement me
	panic("implement me")
}

func (U USB) ReadAll() {
	//TODO implement me
	panic("implement me")
}

func (U USB) Write() {
	//TODO implement me
	panic("implement me")
}

func (U USB) WriteAll() {
	//TODO implement me
	panic("implement me")
}

func (U USB) Close() {
	//TODO implement me
	panic("implement me")
}
