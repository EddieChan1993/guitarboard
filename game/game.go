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

const dpi = 72

type Game struct {
	AllWords map[WordPkId]*Words
	font     font.Face
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
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
	op.GeoM.Translate(150, 30)
	screen.DrawImage(img.EbitenGuitarBoardImg, op)
	g.DrawCircleFloor(screen)
	g.DrawDesc(screen)
	g.DrawWord(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1662, 400
}

//DrawDesc 描述文字
func (g *Game) DrawDesc(dst *ebiten.Image) {
	text.Draw(dst, "S Show/C Hide", g.font, 1300, 400, color.White)
}

//DrawCircleFloor 画底板圆
func (g *Game) DrawCircleFloor(dst *ebiten.Image) {
	for _, words := range g.AllWords {
		if _, had := DefHideWordKeys[words.key]; had {
			//默认隐藏的音名没有底板
			continue
		}
		ebitenutil.DrawCircle(dst, words.X, words.Y, width, color.RGBA{
			R: 236,
			G: 237,
			B: 237,
			A: 255,
		})
	}
}

//DrawWord 画音名
func (g *Game) DrawWord(dst *ebiten.Image) {
	for _, words := range g.AllWords {
		if !words.IsShow {
			continue
		}
		text.Draw(dst, words.key, g.font, int(words.X-width/2), int(words.Y+width/2+5), color.Black)
	}
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
	tt, _ := opentype.Parse(fonts.PressStart2P_ttf)
	mplusNormalFont, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	res := &Game{
		AllWords: make(map[WordPkId]*Words),
		font:     mplusNormalFont,
	}
	xCd := 85
	yCd := 50
	baseX := 220
	baseY := 35
	keyIndex := 0
	var x, y float64
	for i := 0; i < 15*6; i++ {
		line := i / 15
		switch line {
		case 0:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY + 3)
			keyIndex = (i%15 + 5) % 12
		case 1:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY + 5)
			keyIndex = (i % 15) % 12
		case 2:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY + 10)
			keyIndex = (i%15 + 8) % 12
		case 3:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY + 13)
			keyIndex = (i%15 + 3) % 12
		case 4:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY + 15)
			keyIndex = (i%15 + 10) % 12
		case 5:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY + 20)
			keyIndex = (i%15 + 5) % 12
		}
		insWord := InitWords(x, y, WordKeys[keyIndex])
		res.AllWords[insWord.PkId()] = insWord

	}
	return res
}
