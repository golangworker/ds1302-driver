//go:build !tinygo

package ds1302

import (
    "time"
)

// Заглушечная реализация для обычного Go окружения (без TinyGo).
// Она предназначена только для успешного прохождения go get / go list,
// и не взаимодействует с аппаратурой.

type DS1302 struct{}

// NewDS1302 возвращает пустой экземпляр. Параметры не используются в заглушке.
func NewDS1302(_, _, _ any) *DS1302 { return &DS1302{} }

// Init ничего не делает в заглушке.
func (d *DS1302) Init() {}

// SetTime ничего не делает в заглушке.
func (d *DS1302) SetTime(_ time.Time) {}

// ReadTime возвращает нулевое время в заглушке.
func (d *DS1302) ReadTime() time.Time { return time.Time{} }
