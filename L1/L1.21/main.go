/*
Реализовать паттерн проектирования «Адаптер» на любом примере.
Описание: паттерн Adapter позволяет сконвертировать интерфейс одного класса в интерфейс другого, который ожидает клиент.
Продемонстрируйте на простом примере в Go: у вас есть существующий интерфейс (или структура) и другой,
несовместимый по интерфейсу потребитель — напишите адаптер, который реализует нужный интерфейс и делегирует вызовы к встроенному объекту.

Поясните применимость паттерна, его плюсы и минусы, а также приведите реальные примеры использования.
*/
package main

import "fmt"

func main() {

	plugRus := &PlugRUS{}
	plugEU := &PlugEU{}

	plugAdapterRUStoEU := &PlugRUStoEUAdapter{
		rusPlug: plugRus,
	}

	plugAdapterEUtoRUS := &PlugEUtoRUSAdapter{
		euPlug: plugEU,
	}

	outletEU := &PowerOutletEU{}
	outletEU.InsertPlug(plugEU)
	outletEU.InsertPlug(plugAdapterRUStoEU)

	fmt.Println()

	outletRUS := &PowerOutletRUS{}
	outletRUS.InsertPlug(plugRus)
	outletRUS.InsertPlug(plugAdapterEUtoRUS)

}

type Plug interface {
	InsertIntoOutlet()
}

type PowerOutletEU struct {
}

func (ow *PowerOutletEU) InsertPlug(plug Plug) {
	fmt.Println("Inserting plug into the EU outlet")
	plug.InsertIntoOutlet()
}

type PowerOutletRUS struct {
}

func (ow *PowerOutletRUS) InsertPlug(plug Plug) {
	fmt.Println("Inserting plug into the RUS outlet")
	plug.InsertIntoOutlet()
}

type PlugEU struct {
}

func (p *PlugEU) InsertIntoOutlet() {
	fmt.Println("Inserted EU Plug into the outlet")
}

type PlugRUS struct {
}

func (p *PlugRUS) InsertIntoOutlet() {
	fmt.Println("Inserted RUS Plug into the outlet")
}

type PlugRUStoEUAdapter struct {
	rusPlug *PlugRUS
}

func (pa *PlugRUStoEUAdapter) InsertIntoOutlet() {
	fmt.Println("Adapter converts PlugRUS into PlugEU")
	pa.rusPlug.InsertIntoOutlet()
}

type PlugEUtoRUSAdapter struct {
	euPlug *PlugEU
}

func (pa *PlugEUtoRUSAdapter) InsertIntoOutlet() {
	fmt.Println("Adapter converts PlugEU into PlugRUS")
	pa.euPlug.InsertIntoOutlet()
}
