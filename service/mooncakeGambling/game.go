package mooncakeGambling

import (
	"math/rand"
	"sort"
	"time"
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
	rng *rand.Rand
}

// NewMooncakeGame 创建新的博饼游戏实例
func NewMooncakeGame() *MooncakeGame {
	return &MooncakeGame{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// RollDices 掷骰子
func (g *MooncakeGame) RollDices() [6]int {
	var dices [6]int
	for i := 0; i < 6; i++ {
		dices[i] = g.rng.Intn(6) + 1 // 1-6点
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

// ComparePrize 比较两个奖励等级的大小
// 返回值：1表示a>b，-1表示a<b，0表示a==b
func ComparePrize(a, b PrizeLevel) int {
	if a > b {
		return 1
	} else if a < b {
		return -1
	}
	return 0
}

// SortResults 对游戏结果按奖励等级排序（从高到低）
func SortResults(results []GameResult) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].PrizeLevel > results[j].PrizeLevel
	})
}
