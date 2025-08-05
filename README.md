# DS1302 Driver for TinyGo

Драйвер для микросхемы DS1302 Real Time Clock (RTC) для использования с TinyGo на микроконтроллерах ESP32.

## Описание

DS1302 - это недорогая микросхема часов реального времени с низким энергопотреблением. Она обеспечивает секунды, минуты, часы, дату, месяц и год. Дата автоматически корректируется для месяцев с менее чем 31 днем, включая коррекцию для високосного года.

## Подключение к ESP32

| DS1302 Pin | ESP32 GPIO | Описание |
|------------|------------|----------|
| VCC        | 3.3V       | Питание |
| GND        | GND        | Земля |
| CLK        | GPIO18     | Тактовый сигнал |
| DAT        | GPIO19     | Линия данных |
| RST        | GPIO5      | Сигнал сброса |

## Установка

```bash
go get github.com/golangworker/ds1302-driver
```

## Использование

```go
package main

import (
    "machine"
    "time"
"github.com/golangworker/ds1302-driver"
)

func main() {
    // Создаем экземпляр DS1302
    rtc := ds1302.NewDS1302(machine.GPIO18, machine.GPIO19, machine.GPIO5)
    rtc.Init()
    
    // Устанавливаем время
    rtc.SetTime(time.Now())
    
    // Читаем время
    currentTime := rtc.ReadTime()
    println("Current time:", currentTime.String())
}
```

## API

### `NewDS1302(clk, dat, rst machine.Pin) *DS1302`
Создает новый экземпляр драйвера.

### `Init()`
Инициализирует пины GPIO.

### `SetTime(t time.Time)`
Устанавливает время в RTC.

### `ReadTime() time.Time`
Читает текущее время из RTC.

## Лицензия

MIT License
