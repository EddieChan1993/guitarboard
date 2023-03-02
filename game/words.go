package game

import "fmt"

type WordPkId = string

var width = float64(20)

//WordKeys 所有音名
var WordKeys = []string{
	"C", "#C", "D", "#D", "E", "F", "#F", "G", "#G", "A", "#A", "B",
}

//DefHideWordKeys 默认隐藏音名
var DefHideWordKeys = map[string]struct{}{
	"#C": {}, "#D": {}, "#F": {}, "#G": {}, "#A": {},
}

type Words struct {
	X, Y   float64
	key    string
	IsShow bool //音名显示
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
	if _, had := DefHideWordKeys[this_.key]; had {
		this_.IsShow = false
		return
	}
	this_.IsShow = false
}

//Hide 显示字母
func (this_ *Words) Show() {
	if _, had := DefHideWordKeys[this_.key]; had {
		this_.IsShow = false
		return
	}
	this_.IsShow = true
}

func (this_ *Words) In(currentX, currentY float64) bool {
	wordX := this_.X + width
	wordY := this_.Y + width

	if !(currentX >= this_.X-width && currentX <= wordX) {
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
