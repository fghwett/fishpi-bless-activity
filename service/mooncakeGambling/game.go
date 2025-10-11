package mooncakeGambling

import (
	"math/rand/v2"
	"sort"
)

// PrizeLevel 奖励等级
type PrizeLevel int

const (
	PrizeLevelNone        PrizeLevel = 0  // 无奖
	PrizeLevelYiXiu       PrizeLevel = 1  // 一秀（1个4点）
	PrizeLevelErJu        PrizeLevel = 2  // 二举（2个4点）
	PrizeLevelSiJin       PrizeLevel = 3  // 四进（4个相同点数，非4点）
	PrizeLevelSanHong     PrizeLevel = 4  // 三红（3个4点）
	PrizeLevelDuiTang     PrizeLevel = 5  // 对堂（1~6各1个）
	PrizeLevelZSiDianHong PrizeLevel = 6  // 状元四点红（4个4点）
	PrizeLevelZYWuZi      PrizeLevel = 7  // 状元五子登科（5个相同点数，非4点）
	PrizeLevelZYWuHong    PrizeLevel = 8  // 状元五红（5个4点）
	PrizeLevelZYHeiLiuBo  PrizeLevel = 9  // 状元黑六勃（6个相同点数，非1,4点）
	PrizeLevelZBianDiJin  PrizeLevel = 10 // 状元遍地锦（6个1点）
	PrizeLevelZYJinHua    PrizeLevel = 11 // 状元插金花（4个4点+2个1点）
	PrizeLevelZYLiuBo4    PrizeLevel = 12 // 状元六勃红（6个4点）
)

func (level PrizeLevel) IsTop() bool {
	return level >= PrizeLevelZSiDianHong
}

// PrizeLevelName 奖励等级名称映射
var PrizeLevelName = map[PrizeLevel]string{
	PrizeLevelNone:        "无奖",
	PrizeLevelYiXiu:       "一秀",
	PrizeLevelErJu:        "二举",
	PrizeLevelSiJin:       "四进",
	PrizeLevelSanHong:     "三红",
	PrizeLevelDuiTang:     "对堂",
	PrizeLevelZSiDianHong: "状元四点红",
	PrizeLevelZYWuZi:      "状元五子登科",
	PrizeLevelZYWuHong:    "状元五红",
	PrizeLevelZYHeiLiuBo:  "状元黑六勃",
	PrizeLevelZBianDiJin:  "状元遍地锦",
	PrizeLevelZYJinHua:    "状元插金花",
	PrizeLevelZYLiuBo4:    "状元六勃红",
}

// GameResult 游戏结果
type GameResult struct {
	Dices      [6]int     // 6个骰子的点数
	PrizeLevel PrizeLevel // 奖励等级
	PrizeName  string     // 奖励名称
}

// MooncakeGame 博饼游戏
type MooncakeGame struct {
}

// NewMooncakeGame 创建新的博饼游戏实例
func NewMooncakeGame() *MooncakeGame {
	return &MooncakeGame{}
}

// RollDices 掷骰子
func (g *MooncakeGame) RollDices() [6]int {
	var dices [6]int
	for i := 0; i < 6; i++ {
		dices[i] = rand.IntN(6) + 1 // 1-6点
	}
	return dices
}

// Play 进行一次博饼游戏
func (g *MooncakeGame) Play() GameResult {
	dices := g.RollDices()
	prizeLevel := g.CalculatePrize(dices)

	return GameResult{
		Dices:      dices,
		PrizeLevel: prizeLevel,
		PrizeName:  PrizeLevelName[prizeLevel],
	}
}

// CalculatePrize 计算奖励等级
func (g *MooncakeGame) CalculatePrize(dices [6]int) PrizeLevel {
	// 统计每个点数出现的次数
	counts := make(map[int]int)
	for _, dice := range dices {
		counts[dice]++
	}

	count4 := counts[4] // 4点出现的次数

	// 按优先级从高到低判断

	// 状元六勃红：6个4点
	if count4 == 6 {
		return PrizeLevelZYLiuBo4
	}

	// 状元插金花：4个4点+2个1点
	if count4 == 4 && counts[1] == 2 {
		return PrizeLevelZYJinHua
	}

	// 状元遍地锦：6个1点
	if counts[1] == 6 {
		return PrizeLevelZBianDiJin
	}

	// 状元黑六勃：6个相同点数（非1,4点）
	for dice, count := range counts {
		if count == 6 && dice != 1 && dice != 4 {
			return PrizeLevelZYHeiLiuBo
		}
	}

	// 状元五红：5个4点
	if count4 == 5 {
		return PrizeLevelZYWuHong
	}

	// 状元五子登科：5个相同点数（非4点）
	for dice, count := range counts {
		if count == 5 && dice != 4 {
			return PrizeLevelZYWuZi
		}
	}

	// 状元四点红：4个4点
	if count4 == 4 {
		return PrizeLevelZSiDianHong
	}

	// 对堂：1~6各1个
	if len(counts) == 6 {
		hasAll := true
		for i := 1; i <= 6; i++ {
			if counts[i] != 1 {
				hasAll = false
				break
			}
		}
		if hasAll {
			return PrizeLevelDuiTang
		}
	}

	// 三红：3个4点
	if count4 == 3 {
		return PrizeLevelSanHong
	}

	// 四进：4个相同点数（非4点）
	for dice, count := range counts {
		if count == 4 && dice != 4 {
			return PrizeLevelSiJin
		}
	}

	// 二举：2个4点
	if count4 == 2 {
		return PrizeLevelErJu
	}

	// 一秀：1个4点
	if count4 == 1 {
		return PrizeLevelYiXiu
	}

	// 无奖
	return PrizeLevelNone
}

// PlayWithDices 使用指定的骰子点数进行判定
func (g *MooncakeGame) PlayWithDices(dices [6]int) GameResult {
	prizeLevel := g.CalculatePrize(dices)

	return GameResult{
		Dices:      dices,
		PrizeLevel: prizeLevel,
		PrizeName:  PrizeLevelName[prizeLevel],
	}
}

// ComparePrize 保留旧的按等级比较接口（兼容）
func ComparePrize(a, b PrizeLevel) int {
	if a > b {
		return 1
	} else if a < b {
		return -1
	}
	return 0
}

// extrasForLevel 返回在同级比较时应当参与比较的“多余骰子”列表（未排序）
func extrasForLevel(dices [6]int, level PrizeLevel) []int {
	counts := make(map[int]int)
	for _, d := range dices {
		counts[d]++
	}

	mainVals := make(map[int]bool)
	switch level {
	case PrizeLevelZYLiuBo4:
		mainVals[4] = true
	case PrizeLevelZYJinHua:
		mainVals[4] = true // 4为主，1为次
	case PrizeLevelZBianDiJin:
		mainVals[1] = true
	case PrizeLevelZYHeiLiuBo:
		for v, c := range counts {
			if c == 6 && v != 1 && v != 4 {
				mainVals[v] = true
				break
			}
		}
	case PrizeLevelZYWuHong:
		mainVals[4] = true
	case PrizeLevelZYWuZi:
		for v, c := range counts {
			if c == 5 && v != 4 {
				mainVals[v] = true
				break
			}
		}
	case PrizeLevelZSiDianHong:
		mainVals[4] = true
	case PrizeLevelDuiTang:
		// 对堂没有多余骰子
	case PrizeLevelSanHong:
		mainVals[4] = true
	case PrizeLevelSiJin:
		for v, c := range counts {
			if c == 4 && v != 4 {
				mainVals[v] = true
				break
			}
		}
	case PrizeLevelErJu:
		mainVals[4] = true
	case PrizeLevelYiXiu:
		mainVals[4] = true
	case PrizeLevelNone:
		// 无奖：所有骰子都为多余骰子
	default:
		// 默认把所有当作多余
	}

	// 收集非主值的骰子（或当主值为空时收集全部）
	extras := make([]int, 0)
	if len(mainVals) == 0 {
		for _, d := range dices {
			extras = append(extras, d)
		}
		return extras
	}
	for _, d := range dices {
		if !mainVals[d] {
			extras = append(extras, d)
		}
	}
	return extras
}

// sumInts 计算整数切片总和
func sumInts(arr []int) int {
	s := 0
	for _, v := range arr {
		s += v
	}
	return s
}

// CompareGameResult 按用户要求比较两个 GameResult：先比较等级，再比较多余骰子之和，若相同则对多余骰子降序逐位比较
// 返回 1 表示 a>b，-1 表示 a<b，0 表示相等
func CompareGameResult(a, b GameResult) int {
	// 先比较等级
	if a.PrizeLevel > b.PrizeLevel {
		return 1
	} else if a.PrizeLevel < b.PrizeLevel {
		return -1
	}

	// 相同等级时比较多余骰子
	exA := extrasForLevel(a.Dices, a.PrizeLevel)
	exB := extrasForLevel(b.Dices, b.PrizeLevel)

	sumA := sumInts(exA)
	sumB := sumInts(exB)
	if sumA > sumB {
		return 1
	} else if sumA < sumB {
		return -1
	}

	// 若和相等，则将多余骰子从大到小排序后逐位比较
	sort.Slice(exA, func(i, j int) bool { return exA[i] > exA[j] })
	sort.Slice(exB, func(i, j int) bool { return exB[i] > exB[j] })
	// 比较长度时，长度较长且前缀相同的视为更大（不过长度应相同）
	maxLen := len(exA)
	if len(exB) > maxLen {
		maxLen = len(exB)
	}
	for i := 0; i < maxLen; i++ {
		var va, vb int
		if i < len(exA) {
			va = exA[i]
		} else {
			va = 0
		}
		if i < len(exB) {
			vb = exB[i]
		} else {
			vb = 0
		}
		if va > vb {
			return 1
		} else if va < vb {
			return -1
		}
	}

	return 0
}

// SortResults 对游戏结果按奖励等级 + 多余骰子规则排序（从高到低）
func SortResults(results []GameResult) {
	sort.Slice(results, func(i, j int) bool {
		return CompareGameResult(results[i], results[j]) > 0
	})
}
