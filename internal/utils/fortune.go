package utils

import (
	"fmt"
	"math"
	"time"
)

var heavenlyStems = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
var earthlyBranches = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

var stemToElement = map[string]string{
	"甲": "木", "乙": "木",
	"丙": "火", "丁": "火",
	"戊": "土", "己": "土",
	"庚": "金", "辛": "金",
	"壬": "水", "癸": "水",
}

var branchToElement = map[string]string{
	"寅": "木", "卯": "木",
	"巳": "火", "午": "火",
	"辰": "土", "戌": "土", "丑": "土", "未": "土",
	"申": "金", "酉": "金",
	"亥": "水", "子": "水",
}

var heavenlyStemPairs = map[string]string{
	"甲": "己", "己": "甲",
	"乙": "庚", "庚": "乙",
	"丙": "辛", "辛": "丙",
	"丁": "壬", "壬": "丁",
	"戊": "癸", "癸": "戊",
}

var heavenlyStemClashes = map[string]string{
	"甲": "庚", "庚": "甲",
	"乙": "辛", "辛": "乙",
	"丙": "壬", "壬": "丙",
	"丁": "癸", "癸": "丁",
	"戊": "己", "己": "戊",
}

var earthlyBranchSixPairs = map[string]string{
	"子": "丑", "丑": "子",
	"寅": "亥", "亥": "寅",
	"卯": "戌", "戌": "卯",
	"辰": "酉", "酉": "辰",
	"巳": "申", "申": "巳",
	"午": "未", "未": "午",
}

var earthlyBranchThreePairs = map[string][]string{
	"寅": {"午", "戌"}, "午": {"寅", "戌"}, "戌": {"寅", "午"},
	"亥": {"卯", "未"}, "卯": {"亥", "未"}, "未": {"亥", "卯"},
	"巳": {"酉", "丑"}, "酉": {"巳", "丑"}, "丑": {"巳", "酉"},
	"申": {"子", "辰"}, "子": {"申", "辰"}, "辰": {"申", "子"},
}

var earthlyBranchClashes = map[string]string{
	"子": "午", "午": "子",
	"丑": "未", "未": "丑",
	"寅": "申", "申": "寅",
	"卯": "酉", "酉": "卯",
	"辰": "戌", "戌": "辰",
	"巳": "亥", "亥": "巳",
}

var earthlyBranchResentment = map[string]string{
	"子": "未", "未": "子",
	"丑": "午", "午": "丑",
	"寅": "酉", "酉": "寅",
	"卯": "申", "申": "卯",
	"辰": "亥", "亥": "辰",
	"巳": "戌", "戌": "巳",
}

func CalculateFortunePillars(year, month, day, hour int) (yearStem, yearBranch, monthStem, monthBranch, dayStem, dayBranch, hourStem, hourBranch string) {
	yearStem, yearBranch = calculateYearPillar(year)
	monthStem, monthBranch = calculateMonthPillar(year, month)
	dayStem, dayBranch = calculateDayPillar(year, month, day)
	hourStem, hourBranch = calculateHourPillar(dayStem, hour)
	return
}

func calculateYearPillar(year int) (string, string) {
	idx := (year - 4) % 60
	return heavenlyStems[idx%10], earthlyBranches[idx%12]
}

func calculateMonthPillar(year, month int) (string, string) {
	yearStem, _ := calculateYearPillar(year)
	yearStemIdx := indexOf(heavenlyStems, yearStem)
	
	monthBranchIdx := (month + 1) % 12
	monthStemIdx := (yearStemIdx*2 + monthBranchIdx) % 10
	
	return heavenlyStems[monthStemIdx], earthlyBranches[monthBranchIdx]
}

func calculateDayPillar(year, month, day int) (string, string) {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	daysSince1900 := int(t.Sub(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)).Hours() / 24)
	
	idx := (daysSince1900 + 9) % 60
	return heavenlyStems[idx%10], earthlyBranches[idx%12]
}

func calculateHourPillar(dayStem string, hour int) (string, string) {
	dayStemIdx := indexOf(heavenlyStems, dayStem)
	hourBranchIdx := (hour + 1) / 2 % 12
	hourStemIdx := (dayStemIdx*2 + hourBranchIdx) % 10
	
	return heavenlyStems[hourStemIdx], earthlyBranches[hourBranchIdx]
}

func indexOf(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return 0
}

func CalculateCompatibilityScore(fortune1, fortune2 map[string]string) float64 {
	score := 50.0

	dayStem1 := fortune1["day_stem"]
	dayStem2 := fortune2["day_stem"]
	dayBranch1 := fortune1["day_branch"]
	dayBranch2 := fortune2["day_branch"]

	if IsHeavenlyStemPair(dayStem1, dayStem2) && IsEarthlyBranchSixPair(dayBranch1, dayBranch2) {
		return 100.0
	}

	if IsHeavenlyStemPair(dayStem1, dayStem2) {
		score += 30
	}

	if IsHeavenlyStemClash(dayStem1, dayStem2) {
		score -= 10
	}

	if IsEarthlyBranchSixPair(dayBranch1, dayBranch2) {
		score += 25
	} else if IsEarthlyBranchThreePair(dayBranch1, dayBranch2) {
		score += 20
	}

	if IsEarthlyBranchClash(dayBranch1, dayBranch2) {
		score -= 15
	}

	user1Elements := GetFiveElements(fortune1)
	user2Elements := GetFiveElements(fortune2)
	complementCount := CountComplementaryElements(user1Elements, user2Elements)
	if complementCount >= 2 {
		score += 20
	}

	if HasElementBias(user1Elements, user2Elements) {
		score -= 10
	}

	if IsEarthlyBranchResentment(dayBranch1, dayBranch2) {
		score -= 10
	}

	return math.Min(100, math.Max(0, score))
}

func CalculateSimilarityScore(fortune1, fortune2 map[string]string) float64 {
	dayScore := calculatePillarSimilarity(
		fortune1["day_stem"], fortune1["day_branch"],
		fortune2["day_stem"], fortune2["day_branch"],
	)

	monthScore := calculatePillarSimilarity(
		fortune1["month_stem"], fortune1["month_branch"],
		fortune2["month_stem"], fortune2["month_branch"],
	)

	yearScore := calculatePillarSimilarity(
		fortune1["year_stem"], fortune1["year_branch"],
		fortune2["year_stem"], fortune2["year_branch"],
	)

	return dayScore*0.5 + monthScore*0.3 + yearScore*0.2
}

func calculatePillarSimilarity(stem1, branch1, stem2, branch2 string) float64 {
	score := 0.0

	if stem1 == stem2 {
		score += 50
	} else if GetElement(stem1) == GetElement(stem2) {
		score += 25
	}

	if branch1 == branch2 {
		score += 50
	} else if GetElement(branch1) == GetElement(branch2) {
		score += 25
	}

	return score
}

func CalculateConflictScore(fortune1, fortune2 map[string]string) float64 {
	score := 50.0

	dayBranch1 := fortune1["day_branch"]
	dayBranch2 := fortune2["day_branch"]

	if IsEarthlyBranchClash(dayBranch1, dayBranch2) {
		score -= 30
	}

	if IsEarthlyBranchResentment(dayBranch1, dayBranch2) {
		score -= 25
	}

	user1Elements := GetFiveElements(fortune1)
	user2Elements := GetFiveElements(fortune2)
	if HasElementBias(user1Elements, user2Elements) {
		score -= 15
	}

	return math.Min(100, math.Max(0, score))
}

func GetTodayFortune(fortuneInfo map[string]string) string {
	dayStem := fortuneInfo["day_stem"]
	dayBranch := fortuneInfo["day_branch"]
	
	fortuneMap := map[string]string{
		"甲子": "오늘은 새로운 시작에 좋은 날입니다. 자신감을 가지고 도전해보세요.",
		"乙丑": "인내심이 필요한 하루입니다. 서두르지 말고 차근차근 진행하세요.",
		"丙寅": "활동적인 하루가 예상됩니다. 에너지를 잘 활용하세요.",
		"丁卯": "창의적인 아이디어가 떠오를 수 있는 날입니다.",
		"戊辰": "안정적인 하루입니다. 기존 일을 마무리하는 데 좋습니다.",
		"己巳": "변화를 준비하는 날입니다. 새로운 기회를 주시하세요.",
		"庚午": "의사소통이 중요한 하루입니다. 타인과의 협력이 도움이 됩니다.",
		"辛未": "세심한 주의가 필요한 날입니다. 작은 실수를 조심하세요.",
		"壬申": "유연성이 필요한 하루입니다. 상황에 맞게 대응하세요.",
		"癸酉": "깊이 있는 사고가 필요한 날입니다. 중요한 결정은 신중하게 하세요.",
	}
	
	key := dayStem + dayBranch
	if fortune, ok := fortuneMap[key]; ok {
		return fortune
	}
	
	return "오늘은 평범한 하루입니다. 긍정적인 마음가짐으로 하루를 보내세요."
}

func CalculateTodayPillar() (string, string) {
	now := time.Now()
	return calculateDayPillar(now.Year(), int(now.Month()), now.Day())
}


func GetElement(char string) string {
	if element, ok := stemToElement[char]; ok {
		return element
	}
	if element, ok := branchToElement[char]; ok {
		return element
	}
	return ""
}

func IsHeavenlyStemPair(stem1, stem2 string) bool {
	return heavenlyStemPairs[stem1] == stem2
}

func IsHeavenlyStemClash(stem1, stem2 string) bool {
	return heavenlyStemClashes[stem1] == stem2
}

func IsEarthlyBranchSixPair(branch1, branch2 string) bool {
	return earthlyBranchSixPairs[branch1] == branch2
}

func IsEarthlyBranchThreePair(branch1, branch2 string) bool {
	if pairs, ok := earthlyBranchThreePairs[branch1]; ok {
		for _, pair := range pairs {
			if pair == branch2 {
				return true
			}
		}
	}
	return false
}

func IsEarthlyBranchClash(branch1, branch2 string) bool {
	return earthlyBranchClashes[branch1] == branch2
}

func IsEarthlyBranchResentment(branch1, branch2 string) bool {
	return earthlyBranchResentment[branch1] == branch2
}

func GetFiveElements(fortune map[string]string) map[string]int {
	elements := map[string]int{
		"木": 0, "火": 0, "土": 0, "金": 0, "水": 0,
	}

	allChars := []string{
		fortune["year_stem"], fortune["year_branch"],
		fortune["month_stem"], fortune["month_branch"],
		fortune["day_stem"], fortune["day_branch"],
		fortune["hour_stem"], fortune["hour_branch"],
	}

	for _, char := range allChars {
		if element := GetElement(char); element != "" {
			elements[element]++
		}
	}

	return elements
}

func CountComplementaryElements(user1Elements, user2Elements map[string]int) int {
	count := 0
	for element, count1 := range user1Elements {
		if count1 == 0 && user2Elements[element] >= 2 {
			count++
		}
	}
	return count
}

func HasElementBias(user1Elements, user2Elements map[string]int) bool {
	for element := range user1Elements {
		if user1Elements[element] >= 3 && user2Elements[element] >= 3 {
			return true
		}
	}
	return false
}

func CalculateLuckyElement(fortuneInfo map[string]string, todayStem, todayBranch string) string {
	allElements := map[string]int{
		"木": 0, "火": 0, "土": 0, "金": 0, "水": 0,
	}

	allChars := []string{
		fortuneInfo["year_stem"], fortuneInfo["year_branch"],
		fortuneInfo["month_stem"], fortuneInfo["month_branch"],
		fortuneInfo["day_stem"], fortuneInfo["day_branch"],
		fortuneInfo["hour_stem"], fortuneInfo["hour_branch"],
	}

	allChars = append(allChars, todayStem, todayBranch)

	for _, char := range allChars {
		if element := GetElement(char); element != "" {
			allElements[element]++
		}
	}

	minCount := 999
	luckyElement := "土"

	for element, count := range allElements {
		if count < minCount {
			minCount = count
			luckyElement = element
		}
	}

	return luckyElement
}

func GetLuckyColor(element string) (string, string) {
	colorMap := map[string]struct {
		name string
		hex  string
	}{
		"木": {"초록", "#4CAF50"},
		"火": {"레드", "#F44336"},
		"土": {"옐로우", "#FFC107"},
		"金": {"화이트", "#FFFFFF"},
		"水": {"블루", "#2196F3"},
	}

	if color, ok := colorMap[element]; ok {
		return color.name, color.hex
	}
	return "베이지", "#795548"
}

func GetLuckyNumbers(element string) []int {
	numberMap := map[string][]int{
		"木": {3, 8},
		"火": {2, 7},
		"土": {0, 5},
		"金": {4, 9},
		"水": {1, 6},
	}

	if numbers, ok := numberMap[element]; ok {
		return numbers
	}
	return []int{0, 5}
}

func IntSliceToJSON(nums []int) string {
	if len(nums) == 0 {
		return "[]"
	}
	result := "["
	for i, num := range nums {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d", num)
	}
	result += "]"
	return result
}

