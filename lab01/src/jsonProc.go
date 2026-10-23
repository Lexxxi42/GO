package src

import (
	"fmt"
	"strings"
)

type Elements struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	AtomicNumber             int64             `json:"atomic_number"`
	Name                     string            `json:"name"`
	Symbol                   string            `json:"symbol"`
	AtomicWeight             int64             `json:"atomic_weight"`
	Period                   int64             `json:"period"`
	Group                    int64             `json:"group"`
	Phase                    Phase             `json:"phase"`
	MostStableCrystal        MostStableCrystal `json:"most_stable_crystal"`
	Type                     string            `json:"type"`
	IonicRadius              *float64          `json:"ionic_radius"`
	AtomicRadius             *float64          `json:"atomic_radius"`
	Electronegativity        *float64          `json:"electronegativity"`
	FirstIonizationPotential *float64          `json:"first_ionization_potential"`
	Density                  *float64          `json:"density"`
	MeltingPoint             *float64          `json:"melting_point"`
	BoilingPoint             *float64          `json:"boiling_point"`
	Isotopes                 *int64            `json:"isotopes"`
	Discoverer               string            `json:"discoverer"`
	YearOfDiscovery          string            `json:"year_of_discovery"`
	SpecificHeatCapacity     *float64          `json:"specific_heat_capacity"`
	ElectronConfiguration    string            `json:"electron_configuration"`
	DisplayRow               int64             `json:"display_row"`
	DisplayColumn            int64             `json:"display_column"`
}

type MostStableCrystal string

const (
	Bcc                  MostStableCrystal = "bcc"
	Cub                  MostStableCrystal = "cub"
	Fcc                  MostStableCrystal = "fcc"
	Hex                  MostStableCrystal = "hex"
	Mno                  MostStableCrystal = "mno"
	MostStableCrystalGas MostStableCrystal = "gas"
	Orth                 MostStableCrystal = "orth"
	Rho                  MostStableCrystal = "rho"
	Tet                  MostStableCrystal = "tet"
)

type Phase string

const (
	Artificial Phase = "artificial"
	Liq        Phase = "liq"
	PhaseGas   Phase = "gas"
	Solid      Phase = "solid"
)

type Proverbs struct {
	Description string                `json:"description"`
	Proverbs    []map[string][]string `json:"proverbs"`
}

// ProverbsText возвращает все пословицы как список строк "Тема: пословица".
func ProverbsText(p Proverbs) []string {
	var out []string
	for _, topic := range p.Proverbs { // []map[string][]string
		for theme, list := range topic { // map[string][]string
			for _, text := range list { // []string
				out = append(out, fmt.Sprintf("%s: %s", theme, text))
			}
		}
	}
	return out
}

func printSlice(lines []string) {
	for _, l := range lines {
		fmt.Println(l)
	}
}

// All как отдельная задача: вывести все элементы (таблица).
func PrintProverbs(data Proverbs) {
	printSlice(ProverbsText(data))
}

// ================= Задание 5: пословицы =================

// MaxTopicProverbs ищет тему, в которой больше всего пословиц.
func MaxTopicProverbs(p Proverbs) (string, int) {
	bestTopic := ""
	bestCount := 0
	for _, topic := range p.Proverbs { // каждый topic — карта тема -> список
		for theme, list := range topic {
			if len(list) > bestCount {
				bestTopic = theme
				bestCount = len(list)
			}
		}
	}
	return bestTopic, bestCount
}

// LongestProverb ищет самую длинную пословицу.
func LongestProverb(p Proverbs) string {
	best := ""
	for _, line := range ProverbsText(p) {
		// в строке уже "Тема: пословица", но для длины это ок
		if len(line) > len(best) {
			best = line
		}
	}
	return best
}

// AvgProverbLen считает среднюю длину пословицы (в символах).
func AvgProverbLen(p Proverbs) float64 {
	lines := ProverbsText(p)
	if len(lines) == 0 {
		return 0
	}
	sum := 0
	for _, l := range lines {
		sum += len(l)
	}
	return float64(sum) / float64(len(lines))
}

// XXX
// ElementCell превращает элемент в ячейки строки (переиспользуется в печати).
func ElementCell(el Element) []string {
	r := "Unknown"
	if el.IonicRadius != nil {
		r = fmt.Sprintf("%.4f", *el.IonicRadius)
	}
	return []string{fmt.Sprintf("%d", el.AtomicNumber), el.Name, r}
}

// ElementsRow описывает таблицу элементов для PrintRows.
func ElementsRow(title string) RowDesc[Element] {
	return table(title, []string{"№", "Name", "Radius"}, []int{4, 20, 10}, ElementCell)
}

// All как отдельная задача: вывести все элементы (таблица).
func PrintAllElements(data Elements, title string) {
	PrintRows(data.Elements, ElementsRow(title))
}

// Gases возвращает газообразные элементы.
func Gases(data Elements) []Element {
	return Filter(data.Elements, func(el Element) bool { return el.Phase == PhaseGas })
}

// isMetal верно, если в типе есть слово Metal и это не Metalloid.
func isMetal(el Element) bool {
	return strings.Contains(el.Type, "Metal") && el.Type != "Metalloid"
}

// MetalsAbove возвращает металлы с температурой плавления выше t.
func MetalsAbove(data Elements, t float64) []Element {
	return Filter(data.Elements, func(el Element) bool {
		return isMetal(el) && el.MeltingPoint != nil && *el.MeltingPoint > t
	})
}

// isNonMetal верно, если тип элемента относится к неметаллам.
func isNonMetal(el Element) bool {
	switch el.Type {
	case "Nonmetal", "Halogen", "Noble Gas":
		return true
	}
	return false
}

// LatestNonMetal возвращает неметалл с самой поздней датой открытия.
func LatestNonMetal(data Elements) (Element, bool) {
	var best Element
	var bestYear int
	found := false
	for _, el := range data.Elements {
		if !isNonMetal(el) {
			continue
		}
		y := yearOf(el.YearOfDiscovery)
		if !found || y > bestYear {
			best, bestYear, found = el, y, true
		}
	}
	return best, found
}

// CountByPhase считает элементы по агрегатному состоянию.
func CountByPhase(data Elements) map[Phase]int {
	return CountBy(data.Elements, func(el Element) Phase { return el.Phase })
}
