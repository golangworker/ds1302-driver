//go:build tinygo

// Package ds1302 предоставляет драйвер для микросхемы DS1302 Real Time Clock (RTC)
// для использования с TinyGo на микроконтроллерах ESP32.
//
// DS1302 - это недорогая микросхема часов реального времени с низким энергопотреблением.
// Она обеспечивает секунды, минуты, часы, дату, месяц и год.
// Дата автоматически корректируется для месяцев с менее чем 31 днем,
// включая коррекцию для високосного года.
//
// Пример использования:
//
//     import "github.com/golangworker/ds1302-driver"
//
//     rtc := ds1302.NewDS1302(machine.GPIO18, machine.GPIO19, machine.GPIO5)
//     rtc.Init()
//     rtc.SetTime(time.Now())
//     currentTime := rtc.ReadTime()
//
package ds1302

import (
    "machine"
    "time"
)

// Регистры DS1302 для записи и чтения времени.
// DS1302 использует отдельные адреса для операций чтения и записи.
// Младший бит адреса определяет операцию: 0 - запись, 1 - чтение.
const (
    DS1302_SECONDS_WRITE = 0x80 // Регистр записи секунд (0-59)
    DS1302_SECONDS_READ  = 0x81 // Регистр чтения секунд (0-59)
    DS1302_MINUTES_WRITE = 0x82 // Регистр записи минут (0-59)
    DS1302_MINUTES_READ  = 0x83 // Регистр чтения минут (0-59)
    DS1302_HOURS_WRITE   = 0x84 // Регистр записи часов (0-23, 24-часовой формат)
    DS1302_HOURS_READ    = 0x85 // Регистр чтения часов (0-23, 24-часовой формат)
    DS1302_DATE_WRITE    = 0x86 // Регистр записи даты месяца (1-31)
    DS1302_DATE_READ     = 0x87 // Регистр чтения даты месяца (1-31)
    DS1302_MONTH_WRITE   = 0x88 // Регистр записи месяца (1-12)
    DS1302_MONTH_READ    = 0x89 // Регистр чтения месяца (1-12)
    DS1302_DAY_WRITE     = 0x8A // Регистр записи дня недели (1-7)
    DS1302_DAY_READ      = 0x8B // Регистр чтения дня недели (1-7)
    DS1302_YEAR_WRITE    = 0x8C // Регистр записи года (00-99, представляет 2000-2099)
    DS1302_YEAR_READ     = 0x8D // Регистр чтения года (00-99, представляет 2000-2099)
    DS1302_WP_WRITE      = 0x8E // Регистр записи защиты от записи (0x00 - разрешить, 0x80 - запретить)
    DS1302_WP_READ       = 0x8F // Регистр чтения защиты от записи
)

// DS1302 представляет драйвер для микросхемы DS1302 Real Time Clock.
// Структура содержит пины для взаимодействия с микросхемой через 3-проводной интерфейс.
//
// Подключение к ESP32:
//   - CLK (Serial Clock): Тактовый сигнал для синхронизации передачи данных
//   - DAT (Serial Data): Двунаправленная линия данных
//   - RST (Reset): Сигнал выбора микросхемы (активный высокий уровень)
//
// DS1302 использует последовательный протокол передачи данных,
// где каждый байт передается младшими битами вперед (LSB first).
type DS1302 struct {
    clk machine.Pin  // CLK (Serial Clock) - тактовый сигнал
    dat machine.Pin  // DAT (Serial Data) - линия передачи данных
    rst machine.Pin  // RST (Reset) - сигнал выбора микросхемы
}

// NewDS1302 создает новый экземпляр DS1302
func NewDS1302(clk, dat, rst machine.Pin) *DS1302 {
    return &DS1302{
        clk: clk,
        dat: dat,
        rst: rst,
    }
}

// Init инициализирует DS1302
func (d *DS1302) Init() {
    d.clk.Configure(machine.PinConfig{Mode: machine.PinOutput})
    d.dat.Configure(machine.PinConfig{Mode: machine.PinOutput})
    d.rst.Configure(machine.PinConfig{Mode: machine.PinOutput})
    
    d.clk.Low()
    d.rst.Low()
    d.dat.Low()
}

// writeByte записывает байт в DS1302
func (d *DS1302) writeByte(data uint8) {
    d.dat.Configure(machine.PinConfig{Mode: machine.PinOutput})
    
    for i := 0; i < 8; i++ {
        if data&(1<<i) != 0 {
            d.dat.High()
        } else {
            d.dat.Low()
        }
        d.clk.High()
        time.Sleep(time.Microsecond)
        d.clk.Low()
        time.Sleep(time.Microsecond)
    }
}

// readByte читает байт из DS1302
func (d *DS1302) readByte() uint8 {
    var data uint8
    d.dat.Configure(machine.PinConfig{Mode: machine.PinInput})
    
    for i := 0; i < 8; i++ {
        d.clk.High()
        time.Sleep(time.Microsecond)
        if d.dat.Get() {
            data |= (1 << i)
        }
        d.clk.Low()
        time.Sleep(time.Microsecond)
    }
    return data
}

// writeRegister записывает в регистр DS1302
func (d *DS1302) writeRegister(reg, value uint8) {
    d.rst.High()  // Начать передачу
    d.writeByte(reg)
    d.writeByte(value)
    d.rst.Low()   // Закончить передачу
}

// readRegister читает из регистра DS1302
func (d *DS1302) readRegister(reg uint8) uint8 {
    d.rst.High()  // Начать передачу
    d.writeByte(reg)
    value := d.readByte()
    d.rst.Low()   // Закончить передачу
    return value
}

// bcdToDec конвертирует BCD в десятичное
func bcdToDec(bcd uint8) uint8 {
    return ((bcd >> 4) * 10) + (bcd & 0x0F)
}

// decToBcd конвертирует десятичное в BCD
func decToBcd(dec uint8) uint8 {
    return ((dec / 10) << 4) + (dec % 10)
}

// SetTime устанавливает время в DS1302
func (d *DS1302) SetTime(t time.Time) {
    // Отключить защиту от записи
    d.writeRegister(DS1302_WP_WRITE, 0x00)
    
    // Записать время
    d.writeRegister(DS1302_SECONDS_WRITE, decToBcd(uint8(t.Second())))
    d.writeRegister(DS1302_MINUTES_WRITE, decToBcd(uint8(t.Minute())))
    d.writeRegister(DS1302_HOURS_WRITE, decToBcd(uint8(t.Hour())))
    d.writeRegister(DS1302_DATE_WRITE, decToBcd(uint8(t.Day())))
    d.writeRegister(DS1302_MONTH_WRITE, decToBcd(uint8(t.Month())))
    d.writeRegister(DS1302_YEAR_WRITE, decToBcd(uint8(t.Year()-2000)))
    
    // Включить защиту от записи
    d.writeRegister(DS1302_WP_WRITE, 0x80)
}

// ReadTime читает время из DS1302
func (d *DS1302) ReadTime() time.Time {
    seconds := bcdToDec(d.readRegister(DS1302_SECONDS_READ) & 0x7F)
    minutes := bcdToDec(d.readRegister(DS1302_MINUTES_READ))
    hours := bcdToDec(d.readRegister(DS1302_HOURS_READ))
    day := bcdToDec(d.readRegister(DS1302_DATE_READ))
    month := bcdToDec(d.readRegister(DS1302_MONTH_READ))
    year := int(2000) + int(bcdToDec(d.readRegister(DS1302_YEAR_READ)))
    
    return time.Date(int(year), time.Month(month), int(day), 
                    int(hours), int(minutes), int(seconds), 0, time.UTC)
}

