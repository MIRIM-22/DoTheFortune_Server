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
//일진 분석 엔진(용신 개념, 십성 판별
// 십성(十星) 매핑
var tenStarsMap = map[string]map[string]string{
	"木": {"木": "比肩", "火": "食神", "土": "偏財", "金": "七殺", "水": "偏印"},
	"火": {"火": "比肩", "土": "食神", "金": "偏財", "水": "七殺", "木": "偏印"},
	"土": {"土": "比肩", "金": "食神", "水": "偏財", "木": "七殺", "火": "偏印"},
	"金": {"金": "比肩", "水": "食神", "木": "偏財", "火": "七殺", "土": "偏印"},
	"水": {"水": "比肩", "木": "食神", "火": "偏財", "土": "七殺", "金": "偏印"},
}

// 용신(用神) 판단
var godOfUseMap = map[string]string{
	"木": "火", "火": "木", "土": "木", "金": "水", "水": "金",
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

//형(刑)’ / ‘역마살’ / ‘공망’ 없음 수정함
// 지지형(刑)
var earthlyBranchPunishment = map[string][]string{
	"子": {"卯"}, "卯": {"子"},
	"丑": {"戌", "未"}, "戌": {"丑"}, "未": {"丑"},
	"寅": {"巳", "亥"}, "巳": {"寅"}, "亥": {"寅"},
	"申": {"巳", "寅"}, // 신형(申刑)
	"午": {"午"}, // 자형(自刑)
	"酉": {"酉"},
}

// 역마살(馬)
var flyingHorseBranches = map[string][]string{
	"寅": {"申"}, "申": {"寅"},
	"亥": {"巳"}, "巳": {"亥"},
	"巳": {"亥"}, "亥": {"巳"},
	"申": {"寅"}, "寅": {"申"},
}

// 공망(空亡)
var emptyTrunkBranches = map[string][]string{
	"甲乙": {"戌", "亥"},
	"丙丁": {"子", "丑"},
	"戊己": {"寅", "卯"},
	"庚辛": {"辰", "巳"},
	"壬癸": {"午", "未"},
}

// 천을귀인(天乙貴人)
var nobleStemMap = map[string][]string{
	"甲": {"寅", "午"},
	"乙": {"卯", "未"},
	"丙": {"巳", "午"},
	"丁": {"午", "未"},
	"戊": {"巳", "午"},
	"己": {"午", "未"},
	"庚": {"申", "酉"},
	"辛": {"申", "酉"},
	"壬": {"亥", "子"},
	"癸": {"亥", "子"},
}

// 월지 가충지(喜衝地)
var monthNourishingBranches = map[int]string{
	1: "寅", 2: "卯", 3: "巳", 4: "午",
	5: "未", 6: "未", 7: "申", 8: "酉",
	9: "戌", 10: "戌", 11: "子", 12: "丑",
}

//2번 피드백
//궁합 상세 결과에 오행 분포 데이터(목, 화, 토, 금, 수 총 5개)와 4대 카테고리(대화, 감정 등) 추가 필요합니다.
type FortuneResult struct {
	YearStem   string
	YearBranch string
	MonthStem  string
	MonthBranch string
	DayStem    string
	DayBranch  string
	HourStem   string
	HourBranch string
}

type CompatibilityDetail struct {
	Score              float64                    `json:"score"`
	ElementDistribution map[string]int            `json:"element_distribution"`
	Categories         map[string]CategoryScore  `json:"categories"`
	Details            string                    `json:"details"`
}

type CategoryScore struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

type FortunePrediction struct {
	Score       float64               `json:"score"`
	Keywords    map[string]string     `json:"keywords"` // 재물, 애정, 건강, 총운
	LuckyItems  LuckyItems            `json:"lucky_items"`
	Analysis    string                `json:"analysis"`
}

type LuckyItems struct {
	Element string `json:"element"`
	Color   string `json:"color"`
	Numbers []int  `json:"numbers"`
}

type DailyAnalysis struct {
	GodOfUse      string
	TenStar       string
	StemRelation  string
	BranchRelation string
	HasNobleInfluence bool
	HasFlyingHorse bool
	HasEmptyTrunk bool
}

type SimilarityResultItem struct {
	Score  float64
	Rank   int
	UserID string
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

//삼주 가중치 계산 수정
func CalculateSaJuWeightedScore(dayScore, monthScore, yearScore, hourScore float64) float64 {
	return dayScore*0.4 + monthScore*0.3 + yearScore*0.2 + hourScore*0.1
}

//일(50%) + 월(30%) + 연(20%) = 100%
func CalculateSamJuWeightedScore(dayScore, monthScore, yearScore float64) float64 {
	return dayScore*0.5 + monthScore*0.3 + yearScore*0.2
}


func CalculateCompatibilityScore(fortune1, fortune2 map[string]string) CompatibilityDetail {
	detail := CompatibilityDetail{
		Score:              0.0,
		ElementDistribution: make(map[string]int),
		Categories:        make(map[string]CategoryScore),
	}

	// 일주 점수
	dayScore := calculatePillarCompatibility(
		fortune1["day_stem"], fortune1["day_branch"],
		fortune2["day_stem"], fortune2["day_branch"],
	) * 0.4

	// 월주 점수
	monthScore := calculatePillarCompatibility(
		fortune1["month_stem"], fortune1["month_branch"],
		fortune2["month_stem"], fortune2["month_branch"],
	) * 0.3

	// 연주 점수
	yearScore := calculatePillarCompatibility(
		fortune1["year_stem"], fortune1["year_branch"],
		fortune2["year_stem"], fortune2["year_branch"],
	) * 0.2

	//시주
	hourScore := calculatePillarCompatibility(
		fortune1["hour_stem"], fortune1["hour_branch"],
		fortune2["hour_stem"], fortune2["hour_branch"],
	) * 0.1

	detail.Score = dayScore + monthScore + yearScore

	// 오행 분포
	elem1 := GetFiveElements(fortune1)
	detail.ElementDistribution = elem1

	// 4대 카테고리
	detail.Categories = CalculateCategories(fortune1, fortune2, elem1)

	return detail
}

func calculatePillarCompatibility(stem1, branch1, stem2, branch2 string) float64 {
	score := 50.0

	// 천간합: +20
	if IsHeavenlyStemPair(stem1, stem2) {
		score += 20
	}

	// 지지합: +20
	if IsEarthlyBranchSixPair(branch1, branch2) {
		score += 20
	} else if IsEarthlyBranchThreePair(branch1, branch2) {
		score += 20
	}
	//조후 보완 +15
	elem1 := GetElement(stem1)
	elem2 := GetElement(stem2)
	if isElementGenerating(elem1, elem2) || isElementGenerating(elem2, elem1) {
		score += 15
	}
	// 부정 요소
	if IsHeavenlyStemClash(stem1, stem2) {
		score -= 10
	}
	if IsEarthlyBranchClash(branch1, branch2) {
		score -= 15
	}
	if IsEarthlyBranchPunishment(branch1, branch2) {
		score -= 15
	}
	if IsEarthlyBranchResentment(branch1, branch2) {
		score -= 10
	}

	return math.Min(100, math.Max(0, score))
}

//카테고리 맵핑 로직
func CalculateCategories(fortune1, fortune2 map[string]string, elem1 map[string]int) map[string]CategoryScore {
	categories := make(map[string]CategoryScore)

	// 대화(의사소통)
	comScore := 50.0
	if IsHeavenlyStemPair(fortune1["day_stem"], fortune2["day_stem"]) {
		comScore += 20
	}
	categories["대화"] = CategoryScore{Name: "대화", Score: math.Min(100, comScore)}

	// 감정
	emotScore := 50.0
	if IsEarthlyBranchSixPair(fortune1["month_branch"], fortune2["month_branch"]) {
		emotScore += 20
	}
	categories["감정"] = CategoryScore{Name: "감정", Score: math.Min(100, emotScore)}

	// 재물
	wealthScore := 50.0
	if elem1["木"] > 0 {
		wealthScore += 15
	}
	categories["재물"] = CategoryScore{Name: "재물", Score: math.Min(100, wealthScore)}

	// 건강
	healthScore := 50.0
	if IsEarthlyBranchThreePair(fortune1["day_branch"], fortune2["day_branch"]) {
		healthScore += 15
	}
	categories["건강"] = CategoryScore{Name: "건강", Score: math.Min(100, healthScore)}

	return categories
}

//십성 계산
func CalculateTenStar(userStem, todayStem string) string {
	userElement := GetElement(userStem)
	todayElement := GetElement(todayStem)
	
	if starMap, ok := tenStarsMap[userElement]; ok {
		if star, ok := starMap[todayElement]; ok {
			return star
		}
	}
	return "기타"
}

//일진 분석 엔딩
func AnalyzeDailyPillar(fortune map[string]string, todayStem, todayBranch string) DailyAnalysis {
	analysis := DailyAnalysis{}

	userDayStem := fortune["day_stem"]
	userDayBranch := fortune["day_branch"]

	// 용신
	userElement := GetElement(userDayStem)
	analysis.GodOfUse = godOfUseMap[userElement]

	// 십성: User 일간 vs Today 천간
	analysis.TenStar = CalculateTenStar(userDayStem, todayStem)

	// User 천간 vs Today 천간
	if IsHeavenlyStemPair(userDayStem, todayStem) {
		analysis.StemRelation = "합"
	} else if IsHeavenlyStemClash(userDayStem, todayStem) {
		analysis.StemRelation = "충"
	} else {
		analysis.StemRelation = "중립"
	}

	// Today 지지 vs User 일지
	if IsEarthlyBranchSixPair(userDayBranch, todayBranch) {
		analysis.BranchRelation = "육합"
	} else if IsEarthlyBranchThreePair(userDayBranch, todayBranch) {
		analysis.BranchRelation = "삼합"
	} else if IsEarthlyBranchClash(userDayBranch, todayBranch) {
		analysis.BranchRelation = "충"
	} else if IsEarthlyBranchPunishment(userDayBranch, todayBranch) {
		analysis.BranchRelation = "형"
	} else {
		analysis.BranchRelation = "중립"
	}

	// 천을귀인
	if IsNobleInfluence(userDayStem, todayBranch) {
		analysis.HasNobleInfluence = true
	}

	// 역마살
	if HasFlyingHorse(userDayBranch, todayBranch) {
		analysis.HasFlyingHorse = true
	}

	// 공망
	if HasEmptyTrunk(userDayBranch, todayBranch) {
		analysis.HasEmptyTrunk = true
	}

	return analysis
}

func IsNobleInfluence(stem, branch string) bool {
	if nobles, ok := nobleStemMap[stem]; ok {
		for _, noble := range nobles {
			if noble == branch {
				return true
			}
		}
	}
	return false
}

func HasFlyingHorse(userBranch, todayBranch string) bool {
	if clash, ok := flyingHorseBranches[userBranch]; ok {
		return clash == todayBranch
	}
	return false
}

func HasEmptyTrunk(userBranch, todayBranch string) bool {
	for _, emptyBranches := range emptyTrunkBranches {
		for _, empty := range emptyBranches {
			if empty == todayBranch {
				return true
			}
		}
	}
	return false
}

func CalculateTodayFortune(fortune map[string]string) FortunePrediction {
	now := time.Now()
	todayStem, todayBranch := calculateDayPillar(now.Year(), int(now.Month()), now.Day())

	prediction := FortunePrediction{
		Score:    70.0,
		Keywords: make(map[string]string),
	}

	userDayBranch := fortune["day_branch"]
	userDayStem := fortune["day_stem"]

	// 지지합: +20
	if IsEarthlyBranchSixPair(userDayBranch, todayBranch) {
		prediction.Score += 20
	}

	// 용신운: +15
	godOfUse := godOfUseMap[GetElement(userDayStem)]
	if GetElement(todayBranch) == godOfUse {
		prediction.Score += 15
	}

	// 천을귀인: +10
	if IsNobleInfluence(userDayStem, todayBranch) {
		prediction.Score += 10
	}

	// 지지충: -20
	if IsEarthlyBranchClash(userDayBranch, todayBranch) {
		prediction.Score -= 20
	}

	// 기신운: -15
	if IsEarthlyBranchPunishment(userDayBranch, todayBranch) {
		prediction.Score -= 15
	}

	prediction.Score = math.Min(100, math.Max(0, prediction.Score))

	// 4대 운세 키워드
	prediction.Keywords["재물"] = getWealthFortune(GetElement(todayStem))
	prediction.Keywords["애정"] = getEmotionFortune(todayBranch)
	prediction.Keywords["건강"] = getHealthFortune(userDayStem)
	prediction.Keywords["총운"] = fmt.Sprintf("%.0f", prediction.Score)

	return prediction
}

func getWealthFortune(element string) string {
	fortuneMap := map[string]string{
		"木": "진행 중인 프로젝트에서 성과가 기대됩니다.",
		"火": "창의적인 아이디어가 수익으로 이어질 수 있습니다.",
		"土": "안정적인 재정 관리 시기입니다.",
		"金": "투자나 거래에서 신중함이 필요합니다.",
		"水": "유동성 있는 기회가 찾아올 수 있습니다.",
	}
	if fortune, ok := fortuneMap[element]; ok {
		return fortune
	}
	return "재물운이 평온합니다."
}

func getEmotionFortune(branch string) string {
	fortuneMap := map[string]string{
		"子": "차분한 감정 상태입니다.",
		"丑": "인내심이 필요한 시기입니다.",
		"寅": "활기찬 감정이 넘칩니다.",
		"卯": "섬세한 감정 표현이 좋습니다.",
		"辰": "신중한 판단이 필요합니다.",
		"巳": "열정적인 기운이 강합니다.",
		"午": "감정 표현이 적극적입니다.",
		"未": "부드러운 에너지가 흐릅니다.",
		"申": "이성적인 판단이 우선입니다.",
		"酉": "말조심이 필요한 시기입니다.",
		"戌": "충실한 기운이 흐릅니다.",
		"亥": "휴식과 성찰이 필요합니다.",
	}
	if fortune, ok := fortuneMap[branch]; ok {
		return fortune
	}
	return "감정이 평온합니다."
}

func getHealthFortune(stem string) string {
	fortuneMap := map[string]string{
		"甲": "신체 활동이 좋은 시기입니다.",
		"乙": "휴식을 충분히 취하세요.",
		"丙": "에너지가 넘치는 시기입니다.",
		"丁": "스트레스 관리가 중요합니다.",
		"戊": "소화기 건강을 챙기세요.",
		"己": "면역력 강화가 필요합니다.",
		"庚": "호흡기 건강에 주의하세요.",
		"辛": "피부 관리를 신경 쓰세요.",
		"壬": "수분 섭취에 신경 쓰세요.",
		"癸": "신장 건강을 돌보세요.",
	}
	if fortune, ok := fortuneMap[stem]; ok {
		return fortune
	}
	return "건강한 하루입니다."
}
//행운 아이템 로직
func CalculateLuckyItem(fortune map[string]string) LuckyItems {
	elements := GetFiveElements(fortune)
	
	// 억부: 가장 부족한 오행
	var luckyElement string
	minCount := 999
	
	for element, count := range elements {
		if count < minCount {
			minCount = count
			luckyElement = element
		}
	}

	// 월지 가충지 우선 로직
	monthBranch := fortune["month_branch"]
	nowMonth := time.Now().Month()
	nourishing := monthNourishingBranches[int(nowMonth)]
	
	if nourishing != "" {
		if nourishingElement := GetElement(nourishing); nourishingElement != "" {
			if elements[nourishingElement] < minCount {
				luckyElement = nourishingElement
			}
		}
	}

	//조후 우선 로직
	dayStem := fortune["day_stem"]
	dayElement := GetElement(dayStem)
	
	// 조후 판단: 일간과 상생 관계인 오행 중 가장 부족한 것을 선택
	if luckyConditionElement := findLuckyConditionElement(dayElement, elements); luckyConditionElement != "" {
		luckyElement = luckyConditionElement
	}

	// 통관 판단
	if minCount == 0 {
		luckyElement = findTransitionElement(fortune)
	}

	colorName, colorHex := GetLuckyColor(luckyElement)
	luckyNumbers := GetLuckyNumbers(luckyElement)

	return LuckyItems{
		Element: luckyElement,
		Color:   colorName + "(" + colorHex + ")",
		Numbers: luckyNumbers,
	}
}

// 조후 우선 로직: 일간의 조후(약한 오행) 찾기
// 일간과 상생 관계인 오행 중 가장 부족한 오행 선택
func findLuckyConditionElement(dayElement string, elements map[string]int) string {
	// 일간이 생하는 오행(자식): 일간 -> 조후
	// 예: 목일간 -> 화를 생함 -> 화가 조후
	generatingMap := map[string]string{
		"木": "火", "火": "土", "土": "金", "金": "水", "水": "木",
	}

	if condition, ok := generatingMap[dayElement]; ok {
		// 조후가 가장 부족하면 그것을 선택
		minCount := 999
		var luckyCondition string
		
		if elements[condition] < minCount {
			minCount = elements[condition]
			luckyCondition = condition
		}
		
		if luckyCondition != "" {
			return luckyCondition
		}
	}
	
	return ""
}

func findTransitionElement(fortune map[string]string) string {
	dayStem := fortune["day_stem"]
	stemElement := GetElement(dayStem)
	
	// 오행 상생: 木→火→土→金→水→木
	transitionMap := map[string]string{
		"木": "火", "火": "土", "土": "金", "金": "水", "水": "木",
	}
	
	if element, ok := transitionMap[stemElement]; ok {
		return element
	}
	return "土"
}


//유사도 동점자 처리
func HandleSimilarityTie(results []SimilarityResultItem) []SimilarityResultItem {
	// 점수 내림차순 정렬
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		// 동점이면 UserID 알파벳순
		return results[i].UserID < results[j].UserID
	})

	// 순위 지정
	rank := 1
	for i := 0; i < len(results); i++ {
		if i > 0 && results[i].Score != results[i-1].Score {
			rank = i + 1
		}
		results[i].Rank = rank
	}

	return results
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

func IsEarthlyBranchPunishment(branch1, branch2 string) bool {
	if punishments, ok := earthlyBranchPunishment[branch1]; ok {
		for _, punishment := range punishments {
			if punishment == branch2 {
				return true
			}
		}
	}
	return false
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

//오늘의 운세 프롬프트 기준
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

