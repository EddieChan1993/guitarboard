package game

import "fmt"

type WordPkId = string

var width = float64(27)

//WordKeys 所有音名
var WordKeys = []string{
	"C", "bC", "D", "bD", "E", "F", "bF", "G", "bG", "A", "bA", "B",
}

//defShowWordKeys 默认显示音名
var defShowWordKeys = map[string]struct{}{
	"bC": {}, "bD": {}, "bF": {}, "bG": {}, "bA": {},
}

type Words struct {
	X, Y   float64
	key    string
	IsShow bool //开启遮罩
}

func InitWords(x, y float64, key string) *Words {
	works := &Words{
		X:   x,
		Y:   y,
		key: key,
	}
	works.Hide()
	return works
}

func (this_ *Words) Trigger() {
	if this_.IsShow {
		this_.IsShow = false
	} else {
		this_.IsShow = true
	}
}

//Hide 隐藏字母
func (this_ *Words) Hide() {
	if _, had := defShowWordKeys[this_.key]; had {
		this_.IsShow = false
		return
	}
	this_.IsShow = true
}

//Hide 显示字母
func (this_ *Words) Show() {
	if _, had := defShowWordKeys[this_.key]; had {
		this_.IsShow = false
		return
	}
	this_.IsShow = false
}

func (this_ *Words) In(currentX, currentY float64) bool {
	wordX := this_.X + width
	wordY := this_.Y + width

	if !(currentX >= this_.X && currentX <= wordX) {
		return false
	}
	if !(currentY >= this_.Y && currentY <= wordY) {
		return false
	}
	return true
}

func (this_ *Words) PkId() WordPkId {
	return fmt.Sprintf("%f_%f", this_.X, this_.Y)
}
