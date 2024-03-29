package game

import (
	"fmt"
)

type WordPkId = string

const DefRootKey = "C" //默认根音

//WordKeys 所有音名
var WordKeys = scale12

type Words struct {
	X, Y   float64
	key    string
	Fret   int  //品
	IsShow bool //音名显示
}

func InitWords(x, y float64, key string, fret int) *Words {
	works := &Words{
		X:      x,
		Y:      y,
		key:    key,
		Fret:   fret,
		IsShow: false,
	}
	works.Hide()
	return works
}

//Trigger 双击
func (this_ *Words) Trigger() {
	if this_.IsShow {
		this_.IsShow = false
	} else {
		this_.IsShow = true
	}
}

//Hide 隐藏音名
func (this_ *Words) Hide() {
	this_.IsShow = false
}

//Show 显示音名
func (this_ *Words) Show() {
	this_.IsShow = true
}

//In 是否点到了音名图标
func (this_ *Words) In(currentX, currentY float64) bool {
	wordX := this_.X + width
	wordY := this_.Y + width

	if !(currentX >= this_.X-width-22 && currentX <= wordX) {
		return false
	}
	if !(currentY >= this_.Y-width && currentY <= wordY) {
		return false
	}
	return true
}

func (this_ *Words) PkId() WordPkId {
	return fmt.Sprintf("%f_%f", this_.X, this_.Y)
}
