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

type Mode = uint8      //模式规则
type WordStyle = uint8 //音名显示类型

const (
	WordKey WordStyle = 1 //音名
	WordNum WordStyle = 2 //唱名
)

const (
	ModeSuper   Mode = 1 //超级模式,只会显示相同品格的音名
	ModeNormal  Mode = 2 //普通模式
	ModeFreedom Mode = 3 //自由模式
)

type Game struct {
	AllWords  map[WordPkId]*Words //所有音名信息
	font      font.Face
	smallFont font.Face
	mode      Mode      //当前模式
	wordStyle WordStyle //当前显示音名样式
	touchFret int       //最近所点品格
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.touchEventThink()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.HideAll()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.ShowAll()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.ChangeWordStyle()
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.ChangeMode(ModeNormal)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.ChangeMode(ModeSuper)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.ChangeMode(ModeFreedom)
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
	if g.wordStyle == WordKey {
		text.Draw(dst, "C Style-Word", g.font, 100, 400, color.White)
	}
	if g.wordStyle == WordNum {
		text.Draw(dst, "C Style-Number", g.font, 100, 400, color.White)
	}
	switch g.mode {
	case ModeNormal:
		text.Draw(dst, "1 Mode-Normal", g.font, 700, 400, color.White)
	case ModeSuper:
		text.Draw(dst, "2 Mode-Super", g.font, 700, 400, color.White)
	case ModeFreedom:
		text.Draw(dst, "3 Mode-Free", g.font, 700, 400, color.White)
	}
	text.Draw(dst, "S Show/H Hide", g.font, 1300, 400, color.White)
}

//DrawCircleFloor 画底板圆
func (g *Game) DrawCircleFloor(dst *ebiten.Image) {
	for _, words := range g.AllWords {
		if g.mode != ModeFreedom {
			if _, is := DefHideWordKeys[words.key]; is {
				//非自由模式不显示
				continue
			}
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
		//非自由模式不显示
		if _, is := DefHideWordKeys[words.key]; is {
			if g.wordStyle == WordKey {
				text.Draw(dst, words.key, g.smallFont, int(words.X-width/2-7), int(words.Y+width/2+2), color.Black)
			}
			if g.wordStyle == WordNum {
				num := WordNumKeys[words.key]
				text.Draw(dst, num, g.smallFont, int(words.X-width/2-7), int(words.Y+width/2+2), color.Black)
			}
		} else {
			if g.wordStyle == WordKey {
				text.Draw(dst, words.key, g.font, int(words.X-width/2), int(words.Y+width/2+5), color.Black)
			}
			if g.wordStyle == WordNum {
				num := WordNumKeys[words.key]
				text.Draw(dst, num, g.font, int(words.X-width/2), int(words.Y+width/2+5), color.Black)
			}
		}
	}
}

//touchEventThink 点击事件
func (g *Game) touchEventThink() {
	cX, cY := ebiten.CursorPosition()
	for _, words := range g.AllWords {
		if !words.In(float64(cX), float64(cY)) {
			continue
		}
		if g.mode != ModeFreedom {
			//非自由模式不显示
			if _, is := DefHideWordKeys[words.key]; is {
				continue
			}
		}
		if g.mode == ModeSuper {
			//超级模式
			if g.touchFret != words.Fret {
				g.HideAll()
			}
		}
		words.Show()
		g.touchFret = words.Fret
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

//ChangeWordStyle 显示切换
func (g *Game) ChangeWordStyle() {
	if g.wordStyle == WordKey {
		g.wordStyle = WordNum
	} else {
		g.wordStyle = WordKey
	}
}

func (g *Game) ChangeMode(mode Mode) {
	g.mode = mode
	g.HideAll()
}

//initXYPos 音名坐标初始化
func (g *Game) initXYPos() {
	xCd := 85     //间隔
	yCd := 50     //间隔
	baseX := 220  //坐标（0，0）距离
	baseY := 35   //坐标（0，0）距离
	keyIndex := 0 //取key的索引
	var x, y float64
	for i := 0; i < 15*6; i++ {
		line := i / 15   //第几行
		fret := i%15 + 1 //品格
		switch line {
		case 0:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY + 3)
			keyIndex = (i%15 + 5) % 12
		case 1:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY + 2)
			keyIndex = (i % 15) % 12
		case 2:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY + 5)
			keyIndex = (i%15 + 8) % 12
		case 3:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY + 9)
			keyIndex = (i%15 + 3) % 12
		case 4:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY + 12)
			keyIndex = (i%15 + 10) % 12
		case 5:
			x = float64(i%15*xCd + baseX)
			y = float64(i/15*yCd + baseY + 18)
			keyIndex = (i%15 + 5) % 12
		}
		//if _, had := DefHideWordKeys[WordKeys[keyIndex]]; had {
		//	continue
		//}
		insWord := InitWords(x, y, WordKeys[keyIndex], fret)
		g.AllWords[insWord.PkId()] = insWord
	}
}

//initFont 字体初始化
func (g *Game) initFont() {
	tt, _ := opentype.Parse(fonts.PressStart2P_ttf)
	mplusNormalFont, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	smailFont, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    18,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	g.font = mplusNormalFont
	g.smallFont = smailFont
}

func NewGame() *Game {
	res := &Game{
		AllWords:  make(map[WordPkId]*Words),
		mode:      ModeSuper,
		wordStyle: WordNum,
	}
	res.initXYPos()
	res.initFont()
	return res
}
