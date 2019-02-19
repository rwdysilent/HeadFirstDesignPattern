// Copyright 2019 pfwu. All rights reserved.
//
// @Author: pfwu
// @Date: 2019/1/30 17:41

//工厂模式

package main

import (
	"fmt"
	"log"
)

type P interface {
	Prepare()
	Bake()
	Cut()
	Box()
}

type Pizza struct {
	Name string
	Type string
}

func (p *Pizza) Prepare() {
	fmt.Println("prepare Done")
}

func (p *Pizza) Bake() {
	fmt.Println("bake Done")
}

func (p *Pizza) Cut() {
	fmt.Println("cut Done")
}

func (p *Pizza) Box() {
	fmt.Println("box Done")
}

// check struct implements the interface
//var _ = P(&Pizza{})
var _ P = &Pizza{}

// pizza 制作工厂函数
type PizzaFactory func() *Pizza

// NewYork factory
func NYPizzaFactory() *Pizza {
	var pizza = new(Pizza)
	pizza.Name = "NewYork Pizza"
	return pizza
}

// chicago factory
func ChicagoPizzaFactory() *Pizza {
	var pizza = new(Pizza)
	pizza.Name = "Chicago Pizza"
	return pizza
}

var pizzaFactories = make(map[string]PizzaFactory)

func RegisterF(name string, factory PizzaFactory) {
	if factory == nil {
		log.Panicf("pizza factory %s dose not exist. ", name)
	}
	_, register := pizzaFactories[name]
	if register {
		log.Printf("pizza factory %s already registered. Ignored.", name)
	}

	pizzaFactories[name] = factory
}

func init() {
	RegisterF("NewYork", NYPizzaFactory)
	RegisterF("Chicago", ChicagoPizzaFactory)
}

// 原料工厂
type IngredientFactory interface {
	CreateDough()
	CreateSauce()
	CreateCheese()
}

type NYIngredientFactory struct {
	Dough  string
	Sauce  string
	Cheese string
}

func (NIn *NYIngredientFactory) CreateDough() {
	NIn.Dough = "NewYork factory Dough "
}

func (NIn *NYIngredientFactory) CreateSauce() {
	NIn.Sauce = "NewYork factory Sauce "
}

func (NIn *NYIngredientFactory) CreateCheese() {
	NIn.Cheese = "NewYork factory Cheese "
}

type ChicagoIngredientFactory struct {
	Dough  string
	Sauce  string
	Cheese string
}

func (CHIn *ChicagoIngredientFactory) CreateDough() {
	CHIn.Dough = "Chicago Dough "
}
func (CHIn *ChicagoIngredientFactory) CreateSauce() {
	CHIn.Sauce = "Chicago Sauce "
}
func (CHIn *ChicagoIngredientFactory) CreateCheese() {
	CHIn.Cheese = "Chicago Cheese "
}

// 商店工厂
type PizzaStore interface {
	CreatePizza(pType string) Pizza
	OrderPizza(pType string) Pizza
}

type NYPizzaStore struct {
	pizza Pizza
	city  string
}

func (NY *NYPizzaStore) CreatePizza(t string) Pizza {
	pizzaFactory, ok := pizzaFactories[NY.city]
	if !ok {
		log.Panicf("no factory %s found", NY.city)
	}

	NY.pizza = *pizzaFactory()
	NY.pizza = NY.Prepare(NY.pizza)

	if t == "cheese" {
		NY.pizza.Name += " NYCheese"
	} else if t == "clam" {
		NY.pizza.Name += " NYClam"
	} else {
		NY.pizza.Name = " not found"
	}
	return NY.pizza
}

func (NY *NYPizzaStore) Cut() {
	fmt.Println("NewYork Cut...Done")
}

func (NY *NYPizzaStore) Prepare(pizza Pizza) Pizza {
	ingredient := &NYIngredientFactory{}
	ingredient.CreateDough()
	ingredient.CreateSauce()
	ingredient.CreateCheese()

	pizza.Type += ingredient.Cheese + ingredient.Sauce + ingredient.Dough
	return pizza
}

func (NY *NYPizzaStore) OrderPizza(t string) Pizza {
	pizza := NY.CreatePizza(t)
	pizza.Bake()
	//pizza.Cut()
	// NewYork Cut func
	NY.Cut()
	pizza.Box()
	return pizza
}

var _ PizzaStore = &NYPizzaStore{}

type ChicagoPizzaStore struct {
	pizza Pizza
	city  string
}

func (CH *ChicagoPizzaStore) Prepare(pizza Pizza) Pizza {
	ingredient := &ChicagoIngredientFactory{}
	ingredient.CreateDough()
	ingredient.CreateSauce()
	ingredient.CreateCheese()

	pizza.Type += ingredient.Dough + ingredient.Sauce + ingredient.Cheese
	return pizza
}

func (CH *ChicagoPizzaStore) CreatePizza(t string) Pizza {
	pizzaFactory, ok := pizzaFactories[CH.city]
	if !ok {
		log.Panicf("no factory %s found", CH.city)
	}

	CH.pizza = *pizzaFactory()
	CH.pizza = CH.Prepare(CH.pizza)

	if t == "cheese" {
		CH.pizza.Name += " ChicagoCheese"
	} else if t == "clam" {
		CH.pizza.Name += " ChicagoClam"
	} else {
		CH.pizza.Name = " not found"
	}
	return CH.pizza
}

func (CH *ChicagoPizzaStore) OrderPizza(t string) Pizza {
	pizza := CH.CreatePizza(t)
	pizza.Prepare()
	pizza.Bake()
	pizza.Cut()
	pizza.Box()
	return pizza
}

var _ PizzaStore = &ChicagoPizzaStore{}

var pizzaStores = make(map[string]PizzaStore)

func RegisterStore(city string, storeFactory PizzaStore) {
	if storeFactory == nil {
		log.Panicf("Store Factory %s not found.", storeFactory)
	}

	_, registered := pizzaStores[city]
	if registered {
		log.Printf("Store Factory %s already registered. Ignoring.", city)
	}

	pizzaStores[city] = storeFactory
}

func init() {
	RegisterStore("NewYork", &NYPizzaStore{city: "NewYork"})
	RegisterStore("Chicago", &ChicagoPizzaStore{city: "Chicago"})
}

func main() {
	// example1
	store := pizzaStores["NewYork"]
	pizza := store.OrderPizza("cheese")
	fmt.Println(pizza)

	// example 2
	store = pizzaStores["Chicago"]
	pizza = store.OrderPizza("clam")
	fmt.Println(pizza)
}
