// Copyright 2019 pfwu. All rights reserved.
//
// @Author: pfwu
// @Date: 2019/1/23 20:43

// 装饰者模式

package 装饰者模式

import "fmt"

//定义咖啡公用接口
type Beverage interface {
	GetDescription() string
	Cost() float32
	GetSize() string
	SetSize(str string)
}

//定义咖啡类型，以下是两个例子

//深焙咖啡
type DarkRoast struct {
	name  string
	price float32
	size  string
}

func (d *DarkRoast) GetDescription() string {
	return d.name
}

func (d *DarkRoast) Cost() float32 {
	if d.GetSize() == "tall" {
		d.price += .10
	} else if d.GetSize() == "grande" {
		d.price += .15
	} else if d.GetSize() == "venti" {
		d.price += .20
	}
	return d.price
}

func (d *DarkRoast) GetSize() string {
	return d.size
}

func (d *DarkRoast) SetSize(str string) {
	d.size = str
}

//浓咖啡
type Espresso struct {
	name  string
	price float32
	size  string
}

func (e *Espresso) GetDescription() string {
	return e.name
}

func (e *Espresso) Cost() float32 {
	if e.GetSize() == "tall" {
		e.price += .10
	} else if e.GetSize() == "grande" {
		e.price += .15
	} else if e.GetSize() == "venti" {
		e.price += .20
	}
	return e.price
}

func (e *Espresso) GetSize() string {
	return e.size
}

func (e *Espresso) SetSize(str string) {
	e.size = str
}

//装饰者
type Mocha struct {
	beverage Beverage
	name     string
	price    float32
}

func (m *Mocha) GetDescription() string {
	return m.beverage.GetDescription() + "+" + m.name
}

func (m *Mocha) Cost() float32 {
	return m.beverage.Cost() + m.price
}

func (m *Mocha) GetSize() string {
	return m.beverage.GetSize()
}

func (m *Mocha) SetSize(str string) {
	m.beverage.SetSize(str)
}

type info map[string]float32

var (
	darkRoast = &DarkRoast{
		name:  "DarkRoast",
		price: 0.99,
	}
	espresso = &Espresso{
		name:  "Espresso",
		price: 1.99,
	}
	infoTable = info{
		"Mocha": 0.20,
		"Soy":   0.30,
		"Milk":  0.55,
		"Whip":  0.88,
	}
)

func main() {
	//example 1
	fmt.Printf("only Espresson: %s, %.2f\n", espresso.GetDescription(), espresso.Cost())

	//example 2
	mocha := &Mocha{
		beverage: darkRoast,
		name:     "Mocha",
		price:    infoTable["Mocha"],
	}
	mocha2 := &Mocha{
		beverage: mocha,
		name:     "Mocha",
		price:    infoTable["Mocha"],
	}
	mocha2.SetSize("tall")

	fmt.Printf("double Mocha and whip: %s, %.2f\n", mocha2.GetDescription(), mocha2.Cost())
}
