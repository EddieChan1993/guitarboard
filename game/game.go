package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"guitarboard/img"
	"image/color"
)

type Game struct {
	AllWords map[WordPkId]*Words
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.ShowWord()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.HideAll()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.ShowAll()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(81, 0)
	screen.DrawImage(img.EbitenGuitarBoardImg, op)
	for _, words := range g.AllWords {
		if !words.IsShow {
			continue
		}
		ebitenutil.DrawRect(screen, words.X, words.Y, width+3, width, color.RGBA{
			R: 236,
			G: 237,
			B: 237,
			A: 255,
		})
	}
	tt, _ := opentype.Parse(fonts.PressStart2P_ttf)
	const dpi = 72
	mplusNormalFont, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	text.Draw(screen, "S Show/C Hide", mplusNormalFont, 1300, 400, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1662, 400
}

//ShowWord 显示某一个
func (g *Game) ShowWord() {
	cX, cY := ebiten.CursorPosition()
	for _, words := range g.AllWords {
		if !words.In(float64(cX), float64(cY)) {
			continue
		}
		words.Show()
	}
}

//ShowAll 展示全部
func (g *Game) ShowAll() {
	for _, words := range g.AllWords {
		words.Show()
	}
}

//HideAll 隐藏全部
func (g *Game) HideAll() {
	for _, words := range g.AllWords {
		words.Hide()
	}
}

func NewGame() *Game {
	res := &Game{
		AllWords: make(map[WordPkId]*Words),
	}
	xCd := 90
	yCd := 60
	baseX := 180
	baseY := 10
	keyIndex := 0
	var x, y float64
	for i := 0; i < 15*6; i++ {
		line := i / 15
		switch line {
		case 0:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY + 5)
			keyIndex = (i%15 + 5) % 12
		case 1:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY)
			keyIndex = (i % 15) % 12
		case 2:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY - 2)
			keyIndex = (i%15 + 8) % 12
		case 3:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY - 4)
			keyIndex = (i%15 + 3) % 12
		case 4:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY - 9)
			keyIndex = (i%15 + 10) % 12
		case 5:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY - 12)
			keyIndex = (i%15 + 5) % 12
		}
		insWord := InitWords(x, y, WordKeys[keyIndex])
		res.AllWords[insWord.PkId()] = insWord

	}
	return res
}
