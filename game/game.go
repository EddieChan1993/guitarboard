package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"guitarboard/img"
	"image/color"
	"strings"
)

const dpi = 72

var width = float64(20) //基准值

type Mode = uint8      //模式规则
type WordStyle = uint8 //音名显示类型

const (
	WordKey WordStyle = 1 //音名
	WordNum WordStyle = 2 //唱名
)

const (
	ModeSuper      Mode = 1 //超级模式,只会显示相同品格的音名
	ModeNormal     Mode = 2 //普通模式
	ModeFreedom    Mode = 3 //自由模式
	ModePentatonic Mode = 4 //五声音阶
)

type Game struct {
	AllWords        map[WordPkId]*Words //所有音名信息
	font            font.Face
	smallFont       font.Face
	mode            Mode                //当前模式
	wordStyle       WordStyle           //当前显示音名样式
	touchFret       int                 //最近所点品格
	limitFret       int                 //最多有效品格
	defRoot         string              //根音
	WordNumKeys     map[string]string   //音名和音级的关系
	DefHideWordKeys map[string]struct{} //非大调组成音
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
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.ChangeMode(ModePentatonic)
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
	return 1900, 400
}

//DrawDesc 描述文字
func (g *Game) DrawDesc(dst *ebiten.Image) {
	offset := 85
	fontColor := color.RGBA{
		R: 95,
		G: 153,
		B: 92,
		A: 255,
	}
	if g.wordStyle == WordKey {
		text.Draw(dst, "C Style-Word", g.font, 100+offset, 400, fontColor)
	}
	if g.wordStyle == WordNum {
		text.Draw(dst, "C Style-Number", g.font, 100+offset, 400, fontColor)
	}
	switch g.mode {
	case ModeNormal:
		text.Draw(dst, fmt.Sprintf("%d Mode-Normal", ModeNormal), g.font, 700+offset, 400, fontColor)
	case ModeSuper:
		text.Draw(dst, fmt.Sprintf("%d Mode-Super", ModeSuper), g.font, 700+offset, 400, fontColor)
	case ModeFreedom:
		text.Draw(dst, fmt.Sprintf("%d Mode-Free", ModeFreedom), g.font, 700+offset, 400, fontColor)
	case ModePentatonic:
		text.Draw(dst, fmt.Sprintf("%d Mode-Pentatonic", ModePentatonic), g.font, 700+offset, 400, fontColor)
	}
	text.Draw(dst, "S Show/H Hide", g.font, 1300+offset, 400, fontColor)
	text.Draw(dst, g.defRoot, g.font, 70, 80, fontColor)
}

//DrawCircleFloor 画底板圆
func (g *Game) DrawCircleFloor(dst *ebiten.Image) {
	if g.mode == ModeFreedom {
		//自由模式
		for _, words := range g.AllWords {
			//只画现实的音名
			if !words.IsShow {
				continue
			}
			if g.IsRoot(words.key) {
				//根音显示外圈
				ebitenutil.DrawCircle(dst, words.X, words.Y, width+6, color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 235,
				})
			}
			ebitenutil.DrawCircle(dst, words.X, words.Y, width, color.RGBA{
				R: 236,
				G: 237,
				B: 237,
				A: 244,
			})
		}
	} else {
		for _, words := range g.AllWords {
			//不显示半音
			if _, is := g.DefHideWordKeys[words.key]; is {
				continue
			}
			if g.mode == ModePentatonic && !g.IsPentatonic(words.key) {
				//五声音阶模式
				continue
			}
			if g.IsRoot(words.key) {
				//根音显示外圈
				ebitenutil.DrawCircle(dst, words.X, words.Y, width+6, color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 235,
				})
			}
			ebitenutil.DrawCircle(dst, words.X, words.Y, width, color.RGBA{
				R: 236,
				G: 237,
				B: 237,
				A: 244,
			})
		}
	}
}

//DrawWord 画音名
func (g *Game) DrawWord(dst *ebiten.Image) {
	fontColor := color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}
	for _, words := range g.AllWords {
		if !words.IsShow {
			continue
		}
		//非自由模式不显示
		showKey := words.key
		if g.wordStyle == WordNum {
			showKey = g.WordNumKeys[words.key]
		}
		if g.mode != ModeFreedom {
			//非自由模式，不显示非组成音
			if _, is := g.DefHideWordKeys[words.key]; is {
				continue
			}
		}
		if g.mode == ModePentatonic && !g.IsPentatonic(words.key) {
			//五声音阶模式
			continue
		}
		if strings.Contains(showKey, "b") {
			//如果有降b，就缩小显示
			text.Draw(dst, showKey, g.smallFont, int(words.X-width/2-7), int(words.Y+width/2+2), fontColor)
		} else {
			text.Draw(dst, showKey, g.font, int(words.X-width/2), int(words.Y+width/2+5), fontColor)
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
		//主音切换
		if g.isPressMoreKey([]ebiten.Key{ebiten.KeyControlLeft}) {
			g.setRoot(words.key)
			return
		}
		if g.mode != ModeFreedom {
			//非自由模式，无法点击非组成音
			if _, is := g.DefHideWordKeys[words.key]; is {
				continue
			}
		}
		if g.mode == ModeSuper {
			//超级模式
			//只保留一个品格显示
			if g.touchFret != words.Fret {
				g.HideAll()
			}
			g.touchFret = words.Fret
		}
		words.Trigger()
	}
}

//setRoot 设置主音
func (g *Game) setRoot(root string) {
	g.WordNumKeys, g.DefHideWordKeys = ScaleSys.ScaleNumsByRoot(root)
	g.defRoot = root
	g.ShowAll() //全部显示
}

//isPressMoreKey 是否按了该键
func (g *Game) isPressMoreKey(needKeys []ebiten.Key) bool {
	nows := inpututil.PressedKeys()
	if len(needKeys) == 0 || len(nows) == 0 {
		return false
	}
	for _, now := range nows {
		for _, need := range needKeys {
			if now == need {
				return true
			}
		}
	}
	return false
}

//ShowAll 展示全部
func (g *Game) ShowAll() {
	for _, words := range g.AllWords {
		if g.mode != ModeFreedom {
			//非自由模式不显示
			if _, is := g.DefHideWordKeys[words.key]; is {
				continue
			}
		}
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
	g.ShowAll()
}

//initXYPos 音名坐标初始化
func (g *Game) initXYPos() {
	xCd := 85     //间隔
	yCd := 50     //间隔
	baseX := 220  //坐标（0，0）距离
	baseY := 35   //坐标（0，0）距离
	keyIndex := 0 //取key的索引
	maxFret := g.limitFret
	var x, y float64
	for i := 0; i < maxFret*6; i++ {
		line := i / maxFret   //第几行
		fret := i%maxFret + 1 //品格
		switch line {
		case 0:
			x = float64(i%maxFret*xCd + baseX)
			y = float64(i/maxFret*yCd + baseY + 3)
			keyIndex = (i%maxFret + 5) % 12
		case 1:
			x = float64(i%maxFret*xCd + baseX)
			y = float64(i/maxFret*yCd + baseY + 2)
			keyIndex = (i % maxFret) % 12
		case 2:
			x = float64(i%maxFret*xCd + baseX)
			y = float64(i/maxFret*yCd + baseY + 5)
			keyIndex = (i%maxFret + 8) % 12
		case 3:
			x = float64(i%maxFret*xCd + baseX)
			y = float64(i/maxFret*yCd + baseY + 9)
			keyIndex = (i%maxFret + 3) % 12
		case 4:
			x = float64(i%maxFret*xCd + baseX)
			y = float64(i/maxFret*yCd + baseY + 12)
			keyIndex = (i%maxFret + 10) % 12
		case 5:
			x = float64(i%maxFret*xCd + baseX)
			y = float64(i/maxFret*yCd + baseY + 18)
			keyIndex = (i%maxFret + 5) % 12
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

//IsRoot 是否是根音
func (g *Game) IsRoot(key string) bool {
	return g.defRoot == key
}

//IsPentatonic 是否是五声音阶
func (g *Game) IsPentatonic(word string) bool {
	nums := g.WordNumKeys[word]
	return nums != "4" && nums != "7"
}

func NewGame() *Game {
	initScale()
	wordNumKeys, defHideWordKeys := ScaleSys.ScaleNumsByRoot(DefRootKey)
	res := &Game{
		AllWords:        make(map[WordPkId]*Words),
		mode:            ModeFreedom,
		wordStyle:       WordKey,
		limitFret:       18,
		defRoot:         DefRootKey,
		WordNumKeys:     wordNumKeys,
		DefHideWordKeys: defHideWordKeys,
	}
	res.initXYPos()
	res.initFont()
	res.ChangeMode(ModeFreedom)
	return res
}
