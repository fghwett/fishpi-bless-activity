package mooncakeGambling

import (
	"fmt"
	"testing"
)

func TestMooncakeGame_Play(t *testing.T) {
	game := NewMooncakeGame()

	// 测试10次游戏
	for i := 0; i < 10; i++ {
		result := game.Play()
		t.Logf("第%d次: 骰子=%v, 奖励=%s (等级=%d)",
			i+1, result.Dices, result.PrizeName, result.PrizeLevel)
	}
}

func TestMooncakeGame_CalculatePrize(t *testing.T) {
	game := NewMooncakeGame()

	tests := []struct {
		name   string
		dices  [6]int
		expect PrizeLevel
	}{
		{"状元六勃红", [6]int{4, 4, 4, 4, 4, 4}, PrizeLevelZYLiuBo4},
		{"状元插金花", [6]int{4, 4, 4, 4, 1, 1}, PrizeLevelZYJinHua},
		{"状元遍地锦", [6]int{1, 1, 1, 1, 1, 1}, PrizeLevelZBianDiJin},
		{"状元黑六勃", [6]int{2, 2, 2, 2, 2, 2}, PrizeLevelZYHeiLiuBo},
		{"状元五红", [6]int{4, 4, 4, 4, 4, 1}, PrizeLevelZYWuHong},
		{"状元五子登科", [6]int{3, 3, 3, 3, 3, 1}, PrizeLevelZYWuZi},
		{"状元四点红", [6]int{4, 4, 4, 4, 1, 2}, PrizeLevelZSiDianHong},
		{"对堂", [6]int{1, 2, 3, 4, 5, 6}, PrizeLevelDuiTang},
		{"三红", [6]int{4, 4, 4, 1, 2, 3}, PrizeLevelSanHong},
		{"四进", [6]int{2, 2, 2, 2, 1, 3}, PrizeLevelSiJin},
		{"二举", [6]int{4, 4, 1, 2, 3, 5}, PrizeLevelErJu},
		{"一秀", [6]int{4, 1, 2, 3, 3, 5}, PrizeLevelYiXiu},
		{"无奖", [6]int{1, 2, 3, 5, 6, 6}, PrizeLevelNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := game.PlayWithDices(tt.dices)
			if result.PrizeLevel != tt.expect {
				t.Errorf("期望奖励等级 %d (%s), 但得到 %d (%s)",
					tt.expect, PrizeLevelName[tt.expect],
					result.PrizeLevel, result.PrizeName)
			} else {
				t.Logf("✓ %s: 骰子=%v, 奖励=%s", tt.name, tt.dices, result.PrizeName)
			}
		})
	}
}

func TestComparePrize(t *testing.T) {
	tests := []struct {
		a      PrizeLevel
		b      PrizeLevel
		expect int
	}{
		{PrizeLevelZYLiuBo4, PrizeLevelYiXiu, 1},
		{PrizeLevelYiXiu, PrizeLevelZYLiuBo4, -1},
		{PrizeLevelSanHong, PrizeLevelSanHong, 0},
	}

	for _, tt := range tests {
		result := ComparePrize(tt.a, tt.b)
		if result != tt.expect {
			t.Errorf("ComparePrize(%d, %d) = %d, 期望 %d", tt.a, tt.b, result, tt.expect)
		}
	}
}

func TestSortResults(t *testing.T) {
	results := []GameResult{
		{PrizeLevel: PrizeLevelYiXiu, PrizeName: "一秀"},
		{PrizeLevel: PrizeLevelZYLiuBo4, PrizeName: "状元六勃红"},
		{PrizeLevel: PrizeLevelSanHong, PrizeName: "三红"},
		{PrizeLevel: PrizeLevelNone, PrizeName: "无奖"},
	}

	SortResults(results)

	// 验证排序结果（从高到低）
	expected := []PrizeLevel{
		PrizeLevelZYLiuBo4,
		PrizeLevelSanHong,
		PrizeLevelYiXiu,
		PrizeLevelNone,
	}

	for i, result := range results {
		if result.PrizeLevel != expected[i] {
			t.Errorf("排序后索引%d的奖励等级 = %d, 期望 %d", i, result.PrizeLevel, expected[i])
		}
		t.Logf("第%d名: %s (等级=%d)", i+1, result.PrizeName, result.PrizeLevel)
	}
}

func ExampleMooncakeGame_Play() {
	game := NewMooncakeGame()
	result := game.Play()
	fmt.Printf("骰子点数: %v\n", result.Dices)
	fmt.Printf("获得奖励: %s\n", result.PrizeName)
}

func ExampleMooncakeGame_PlayWithDices() {
	game := NewMooncakeGame()
	// 使用指定的骰子点数
	dices := [6]int{4, 4, 4, 4, 1, 1}
	result := game.PlayWithDices(dices)
	fmt.Printf("骰子点数: %v\n", result.Dices)
	fmt.Printf("获得奖励: %s\n", result.PrizeName)
	// Output will vary based on implementation
}

// TestMooncakeGame_ProbabilityStatistics 测试10000次博饼后各种奖励的概率分布
func TestMooncakeGame_ProbabilityStatistics(t *testing.T) {
	game := NewMooncakeGame()
	rounds := 10000

	// 统计各个奖励等级出现的次数
	prizeCount := make(map[PrizeLevel]int)

	// 进行10000次博饼
	for i := 0; i < rounds; i++ {
		result := game.Play()
		prizeCount[result.PrizeLevel]++
	}

	t.Logf("\n========== 博饼概率统计 (总次数: %d) ==========", rounds)

	// 按奖励等级从高到低输出统计结果
	prizeLevels := []PrizeLevel{
		PrizeLevelZYLiuBo4,
		PrizeLevelZYJinHua,
		PrizeLevelZBianDiJin,
		PrizeLevelZYHeiLiuBo,
		PrizeLevelZYWuHong,
		PrizeLevelZYWuZi,
		PrizeLevelZSiDianHong,
		PrizeLevelDuiTang,
		PrizeLevelSanHong,
		PrizeLevelSiJin,
		PrizeLevelErJu,
		PrizeLevelYiXiu,
		PrizeLevelNone,
	}

	for _, level := range prizeLevels {
		count := prizeCount[level]
		probability := float64(count) / float64(rounds) * 100
		name := PrizeLevelName[level]

		// 根据概率显示不同长度的进度条
		barLength := int(probability * 2) // 每1%显示2个字符
		if barLength > 100 {
			barLength = 100
		}
		bar := ""
		for i := 0; i < barLength; i++ {
			bar += "█"
		}

		t.Logf("%-12s (等级%2d): %5d次 | %6.2f%% | %s",
			name, level, count, probability, bar)
	}

	t.Logf("=================================================\n")

	// 验证总次数是否正确
	totalCount := 0
	for _, count := range prizeCount {
		totalCount += count
	}
	if totalCount != rounds {
		t.Errorf("统计总次数 %d 不等于实际次数 %d", totalCount, rounds)
	}

	// 验证一些基本的概率规律（粗略验证）
	// 无奖的概率应该是最高的
	if prizeCount[PrizeLevelNone] < rounds/10 {
		t.Logf("警告: 无奖概率似乎偏低 (%.2f%%)", float64(prizeCount[PrizeLevelNone])/float64(rounds)*100)
	}

	// 一秀的概率应该比较高（第二高）
	if prizeCount[PrizeLevelYiXiu] < rounds/20 {
		t.Logf("警告: 一秀概率似乎偏低 (%.2f%%)", float64(prizeCount[PrizeLevelYiXiu])/float64(rounds)*100)
	}

	// 状元级别的奖励应该很少
	champCount := prizeCount[PrizeLevelZYLiuBo4] +
		prizeCount[PrizeLevelZYJinHua] +
		prizeCount[PrizeLevelZBianDiJin] +
		prizeCount[PrizeLevelZYHeiLiuBo] +
		prizeCount[PrizeLevelZYWuHong] +
		prizeCount[PrizeLevelZYWuZi] +
		prizeCount[PrizeLevelZSiDianHong]

	champProbability := float64(champCount) / float64(rounds) * 100
	t.Logf("\n所有状元级别奖励总概率: %.4f%%", champProbability)

	if champProbability > 5 {
		t.Logf("警告: 状元级别奖励概率似乎偏高 (%.4f%%)", champProbability)
	}
}
