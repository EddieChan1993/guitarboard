package game

type Scale12Typ = string            //12个基本音阶
type ScaleNaturalTyp = []Scale12Typ //自然音阶
const ScaleNatureLen = 8            //自然大调总数
const Scale12Len = 12               //12个基本音阶

var MajorNatureInterval = []int{2, 2, 1, 2, 2, 2, 1} //自然大调间隔
var MinorNatureInterval = []int{2, 1, 2, 2, 1, 2, 2} //自然小调间隔

// scale12 8度音名
var scale12 = []Scale12Typ{"C", "bD", "D", "bE", "E", "F", "bG", "G", "bA", "A", "bB", "B"}
var scaleNum12 = []Scale12Typ{"1", "b2", "2", "b3", "3", "4", "b5", "5", "b6", "6", "b7", "7"}

//==================== Api ====================//

var ScaleSys IScaleSys

type IScaleSys interface {
	IsMajorNatural(scale Scale12Typ) bool                //是否是自然大调
	GetMajorNaturalScale(key Scale12Typ) ScaleNaturalTyp //获取对应大调自然音阶
	GetMinorNaturalScale(key Scale12Typ) ScaleNaturalTyp //获取对应小调自然音阶
	ScaleNumsByRoot(root string) (res map[string]string, hide map[string]struct{})
}

func (s *ScaleSysT) GetMajorNaturalScale(key Scale12Typ) ScaleNaturalTyp {
	return s.MajorNaturalScale[key]
}

func (s *ScaleSysT) GetMinorNaturalScale(key Scale12Typ) ScaleNaturalTyp {
	return s.MinorNaturalScale[key]
}

func (s *ScaleSysT) IsMajorNatural(scale Scale12Typ) bool {
	_, had := s.MajorNaturalScale[scale]
	return had
}

//ScaleNumsByRoot 根音获取音名
func (s *ScaleSysT) ScaleNumsByRoot(root string) (res map[string]string, hide map[string]struct{}) {
	index := 0
	//大调组成音
	major := s.GetMajorNaturalScale(root)
	majorMap := make(map[string]struct{}, len(major))
	for _, typ := range major {
		majorMap[typ] = struct{}{}
	}
	//大调排列
	for i, typ := range scale12 {
		if root == typ {
			index = i
			break
		}
	}
	leftHalf := scale12[:index]
	rightHalf := scale12[index:]
	all := make([]Scale12Typ, 0, len(scale12))
	all = append(all, rightHalf...)
	all = append(all, leftHalf...)
	//获取结果
	res = make(map[string]string, len(scale12))           //大调全排列和音级关系
	hide = make(map[string]struct{}, len(all)-len(major)) //大调半音
	for i, typ := range all {
		res[typ] = scaleNum12[i]
		_, had := majorMap[typ]
		if !had {
			//大调组成音中不包含则隐藏
			hide[typ] = struct{}{}
		}
	}
	return res, hide
}

//==================== Obj ====================//

type ScaleSysT struct {
	ScaleIndex        map[Scale12Typ]int             //音阶索引
	MajorNaturalScale map[Scale12Typ]ScaleNaturalTyp //大调自然音阶
	MinorNaturalScale map[Scale12Typ]ScaleNaturalTyp //小调调自然音阶
}

func initScale() {
	ins := &ScaleSysT{
		ScaleIndex:        make(map[Scale12Typ]int, len(scale12)),
		MajorNaturalScale: make(map[Scale12Typ]ScaleNaturalTyp, len(scale12)),
		MinorNaturalScale: make(map[Scale12Typ]ScaleNaturalTyp, len(scale12)),
	}
	for index, c := range scale12 {
		ins.ScaleIndex[c] = index
	}

	//自然音阶
	for _, s := range scale12 {
		sIndex := ins.ScaleIndex[s]
		majorNList := make(ScaleNaturalTyp, 0, ScaleNatureLen)
		majorNList = append(majorNList, s)
		for i := 0; i < ScaleNatureLen-1; i++ {
			sIndex += MajorNatureInterval[i]
			index := sIndex % Scale12Len
			majorNList = append(majorNList, scale12[index])
		}
		ins.MajorNaturalScale[s] = majorNList

		sIndex2 := sIndex
		minorNList := make(ScaleNaturalTyp, 0, ScaleNatureLen)
		minorNList = append(minorNList, s)
		for i := 0; i < ScaleNatureLen-1; i++ {
			sIndex2 += MinorNatureInterval[i]
			index := sIndex2 % Scale12Len
			minorNList = append(minorNList, scale12[index])
		}
		ins.MinorNaturalScale[s] = minorNList
	}
	ScaleSys = ins
}
