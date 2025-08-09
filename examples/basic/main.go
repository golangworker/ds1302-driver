//go:build tinygo

package main

import (
	"machine"
	"time"
"github.com/golangworker/ds1302-driver"
)

func main() {
	// Инициализируем встроенный светодиод
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Создаем экземпляр DS1302
	// CLK -> GPIO18, DAT -> GPIO19, RST -> GPIO5
	rtc := ds1302.NewDS1302(machine.GPIO18, machine.GPIO19, machine.GPIO5)
	rtc.Init()

	// Устанавливаем начальное время (только один раз)
	// В реальном проекте это можно делать через веб-интерфейс или другой способ
	initialTime := time.Date(2024, 8, 5, 21, 0, 0, 0, time.UTC)
	rtc.SetTime(initialTime)
	
	println("DS1302 RTC Example Started!")
	println("Initial time set to:", initialTime.Format("2006-01-02 15:04:05"))

	for {
		// Читаем время из RTC
		currentTime := rtc.ReadTime()
		
		// Выводим время в Serial
		println("RTC Time:", currentTime.Format("2006-01-02 15:04:05"))
		
		// Мигаем светодиодом каждую секунду
		led.High()
		time.Sleep(100 * time.Millisecond)
		led.Low()
		
		// Ждем секунду
		time.Sleep(900 * time.Millisecond)
	}
}
