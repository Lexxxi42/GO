package main

import (
	"fmt"
	"lab01/src"
)

func main() {
	// ===== Задание 2. JSON: элементы =====
	src.Chapter(2)
	els := src.ReadJSON[src.Elements]("data/elements.json")
	fmt.Printf("Всего элементов %d\n\n", len(els.Elements))

	// все элементы таблицей
	src.PrintAllElements(els, "Все элементы:")

	// все газы
	gases := src.Gases(els)
	src.PrintRows(gases, src.ElementsRow("Газы:"))

	// металлы с T плавления выше 2000
	metals := src.MetalsAbove(els, 2000)
	src.PrintRows(metals, src.ElementsRow("Металлы с T>2000°C:"))

	// самый поздний открытый неметалл
	if nm, ok := src.LatestNonMetal(els); ok {
		src.PrintRows([]src.Element{nm}, src.ElementsRow("Самый поздний неметалл:"))
	}

	// подсчёт по фазам
	fmt.Println("Кол-во по агрегатным состояниям:")
	counts := src.CountByPhase(els)
	fmt.Printf("  solid: %d, gas: %d, liq: %d\n",
		counts["solid"], counts["gas"], counts["liq"])

	// ===== Задание 3. CSV: игры =====
	src.Chapter(3)
	games := src.ReadCSV[src.Vgsale]("data/vgsales.csv")

	// 3 самых продаваемых
	src.PrintRows(src.TopSellers(games, 3), src.VgRow("Топ-3 по продажам:"))

	// жанр последней игры Nintendo
	fmt.Printf("Жанр последней игры Nintendo: %s\n\n",
		src.LastNintendoGenre(games))

	// где выше доходы на PlayStation: Европа или Япония
	switch src.PlayStationRegionCompare(games) {
	case 1:
		fmt.Println("На PlayStation доходы выше в Европе")
	case -1:
		fmt.Println("На PlayStation доходы выше в Японии")
	case 0:
		fmt.Println("На PlayStation доходы примерно равны/нет данных")
	}

	src.Chapter(4)
	health := src.ReadCSV[src.Insurance]("data/Health_Insurance.csv")
	src.PrintRows(health, src.InsuranceRow("Страхование здоровья: "))
	fmt.Println("Сколько страховок по регионам:")
	regions := src.CountByRegion(health)
	for r, n := range regions {
		fmt.Printf("  %-12s: %d\n", r, n)
	}
	fmt.Printf("Максимальная стоимость страховки: %.2f\n", src.MaxCharges(health))
	fmt.Printf("Средний возраст застрахованного: %.1f\n", src.AvgAge(health))

	src.Chapter(5)
	proverbs := src.ReadJSON[src.Proverbs]("data/proverbs.json")
	src.PrintProverbs(proverbs)
	topic, n := src.MaxTopicProverbs(proverbs)
	fmt.Printf("Тема с наибольшим числом пословиц: %s (%d шт.)\n", topic, n)
	fmt.Printf("Самая длинная пословица: %s\n", src.LongestProverb(proverbs))
	fmt.Printf("Средняя длина пословицы: %.1f символов\n", src.AvgProverbLen(proverbs))
}
