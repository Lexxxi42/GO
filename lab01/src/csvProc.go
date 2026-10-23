package src

import (
	"strconv"
	"strings"
)

type Vgsales []Vgsale

type Vgsale struct {
	Rank        string  `csv:"Rank"`
	Name        *string `csv:"Name,omitempty"`
	Platform    *string `csv:"Platform,omitempty"`
	Year        *string `csv:"Year,omitempty"`
	Genre       *string `csv:"Genre,omitempty"`
	Publisher   *string `csv:"Publisher,omitempty"`
	NASales     *string `csv:"NA_Sales,omitempty"`
	EUSales     *string `csv:"EU_Sales,omitempty"`
	JPSales     *string `csv:"JP_Sales,omitempty"`
	OtherSales  *string `csv:"Other_Sales,omitempty"`
	GlobalSales *string `csv:"Global_Sales,omitempty"`
}

type Insurances []Insurance

type Insurance struct {
	Age      int     `csv:"age"`
	Sex      *string `csv:"sex"`
	Bmi      float64 `csv:"bmi"`
	Children int     `csv:"children"`
	Smoker   *string `csv:"smoker"`
	Region   *string `csv:"region"`
	Charges  float64 `csv:"charges"`
}

// InsuranceRow описывает таблицу страховок для PrintRows.
func InsuranceRow(title string) RowDesc[Insurance] {
	return table(title,
		[]string{"Age", "Sex", "Bmi", "Childern", "Smoker", "Region", "Charges"},
		[]int{4, 6, 8, 3, 4, 15, 15}, InsuranceSaleCell)
}

// VgSaleCell превращает игру в ячейки строки (Rank, Name, Year).
func InsuranceSaleCell(g Insurance) []string {
	return []string{
		strconv.Itoa(g.Age),
		str(g.Sex),
		strconv.FormatFloat(g.Bmi, 'f', -1, 64),
		strconv.Itoa(g.Children),
		str(g.Smoker),
		str(g.Region),
		strconv.FormatFloat(g.Charges, 'f', -1, 64)}
}

// str возвращает значение указателя на строку, либо default при nil/"-".
func str(p *string) string {
	if p == nil || *p == "" {
		return "-"
	}
	return *p
}

// VgSaleCell превращает игру в ячейки строки (Rank, Name, Year).
func VgSaleCell(g Vgsale) []string {
	return []string{g.Rank, str(g.Name), str(g.Year)}
}

// VgRow описывает таблицу игр для PrintRows.
func VgRow(title string) RowDesc[Vgsale] {
	return table(title, []string{"Rank", "Name", "Year"}, []int{6, 35, 6}, VgSaleCell)
}

// TopSeller возвращает n игр с наибольшими продажами (GlobalSales).
func TopSellers(games []Vgsale, n int) []Vgsale {
	return TopBy(games, n, func(a, b Vgsale) bool {
		return sale(a.GlobalSales) > sale(b.GlobalSales)
	})
}

// LastNintendoGenre возвращает жанр последней выпущенной игры издателя Nintendo.
func LastNintendoGenre(games []Vgsale) string {
	best, bestYear := "", 0
	for _, g := range games {
		if !eq(g.Publisher, "Nintendo") {
			continue
		}
		if y := yearOf(str(g.Year)); y > bestYear {
			best, bestYear = str(g.Genre), y
		}
	}
	return best
}

// isPlayStation верно, если платформа начинается с PS (PS, PS2, PS3, PS4, PSP, PSV).
func isPlayStation(platform string) bool {
	return strings.HasPrefix(platform, "PS")
}

// PlayStationRegionCompare возвращает 1, если в Европе продажи выше,
// -1 если в Японии выше, 0 если равны/нет данных.
func PlayStationRegionCompare(games []Vgsale) int {
	var eu, jp float64
	for _, g := range games {
		if !isPlayStation(str(g.Platform)) {
			continue
		}
		eu += saleF(g.EUSales)
		jp += saleF(g.JPSales)
	}
	switch {
	case eu > jp:
		return 1
	case jp > eu:
		return -1
	default:
		return 0
	}
}

// sale парсит GlobalSales в float (0 при ошибке/nil).
func sale(s *string) float64 {
	if s == nil || *s == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(*s, 64)
	return v
}

// saleF — псевдоним sale, для читаемости.
func saleF(s *string) float64 { return sale(s) }

// eq сравнивает строку-указатель с переданной, игнорируя nil.
func eq(p *string, s string) bool {
	return p != nil && *p == s
}

// ================= Задание 4: страхование =================

// CountByRegion считает, сколько страховок в каждом регионе.
func CountByRegion(health []Insurance) map[string]int {
	m := make(map[string]int)
	for _, h := range health {
		m[str(h.Region)]++
	}
	return m
}

// MaxCharges ищет самую дорогую страховку (charges).
func MaxCharges(health []Insurance) float64 {
	// пусть 0, если список пустой — для лабы сойдёт
	max := 0.0
	for _, h := range health {
		if h.Charges > max {
			max = h.Charges
		}
	}
	return max
}

// AvgAge считает средний возраст застрахованных.
func AvgAge(health []Insurance) float64 {
	if len(health) == 0 {
		return 0
	}
	sum := 0
	for _, h := range health {
		sum += h.Age
	}
	// делим на float, чтобы не потерять дробную часть
	return float64(sum) / float64(len(health))
}
