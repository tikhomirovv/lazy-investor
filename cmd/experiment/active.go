package main

import (
	"fmt"
	"math/rand"
	"sort"
)

type Asset string

const (
	RUB  Asset = "RUB"
	GOLD Asset = "GOLD"
	USD  Asset = "USD"
	OZON Asset = "OZON"
	SBER Asset = "SBER"
)

type Portfolio struct {
	Actives map[Asset]float32
}

type MarketState struct {
	Portfolio Portfolio
	Prices    map[Asset]float32
	History   map[Asset][]float32
}

func (m *MarketState) Act(from Asset, to Asset) bool {
	// Проверяем, есть ли актив в портфеле
	amount, exists := m.Portfolio.Actives[from]
	if !exists || amount <= 0 {
		// fmt.Printf("Не удалось перевести %s -> %s: недостаточно средств\n", from, to)
		return false // Нечего переводить
	}
	// Проверяем наличие цены актива
	fromPrice, fromExists := m.Prices[from]
	toPrice, toExists := m.Prices[to]
	if !fromExists || !toExists || toPrice == 0 {
		// fmt.Printf("Не удалось перевести %s -> %s: отсутствуют цены\n", from, to)
		return false // Невозможно рассчитать перевод
	}

	// Рассчитываем количество нового актива
	valueInRUB := amount * fromPrice  // Стоимость в RUB
	newAmount := valueInRUB / toPrice // Количество нового актива

	// Логируем сделку
	fmt.Printf("Перевод: %s -> %s | %.4f %s (%.2f RUB) -> %.4f %s\n",
		from, to, amount, from, valueInRUB, newAmount, to)

	// Обновляем портфель
	delete(m.Portfolio.Actives, from)    // Полностью убираем старый актив
	m.Portfolio.Actives[to] += newAmount // Добавляем новый актив
	return true
}

type Market struct {
	State MarketState
}

func NewMarket(initialPortfolio map[Asset]float32, initialPrices map[Asset]float32) *Market {
	portfolio := Portfolio{
		Actives: make(map[Asset]float32),
	}
	for asset, price := range initialPortfolio {
		portfolio.Actives[asset] = price
	}
	state := MarketState{
		Prices:  make(map[Asset]float32),
		History: make(map[Asset][]float32),
	}
	for asset, price := range initialPrices {
		state.Portfolio = portfolio
		state.Prices[asset] = price
		state.History[asset] = []float32{price}
	}

	return &Market{State: state}
}

func getRandomChange(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func (m *Market) SimulateNextStep() {
	changes := map[Asset][2]float32{
		RUB:  {0, 0},
		GOLD: {-500, 500},
		USD:  {-5, 5},
		OZON: {-800, 800},
		SBER: {-100, 100},
	}

	beforeAfter := map[Asset][3]float32{}

	// меняем цены и вычисляем изменения
	for asset, changeRange := range changes {
		before := m.State.Prices[asset]
		after := before + getRandomChange(changeRange[0], changeRange[1])
		m.State.Prices[asset] = after
		// fmt.Printf("%s: %.2f -> %.2f\n", asset, before, after)
		change := (after - before) / before * 100
		beforeAfter[asset] = [3]float32{before, after, change}
		m.State.History[asset] = append(m.State.History[asset], after)
	}
	fmt.Println("--------------")
	fmt.Println(beforeAfter)

	// Структура для хранения результатов
	type DiffResult struct {
		Asset1         Asset
		Asset2         Asset
		DiffPercentage float32
	}

	var diffs []DiffResult

	for asset, change := range beforeAfter {
		for asset2, change2 := range beforeAfter {
			diff := (change[2] - change2[2])
			if diff == 0 {
				continue
			}
			diffs = append(diffs, DiffResult{Asset1: asset, Asset2: asset2, DiffPercentage: diff})
			fmt.Printf("%s -> %s: %.2f, %.2f -> %.2f%%\n",
				asset, asset2, change[2], change2[2], diff)
		}
	}

	// Сортируем по убыванию разницы
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].DiffPercentage > diffs[j].DiffPercentage
	})

	// Выводим `X` наибольших значений
	fmt.Println("Изменения:")
	isAct := false
	for i := 0; i < len(diffs); i++ {

		fmt.Printf("%s -> %s: %.2f%%\n", diffs[i].Asset1, diffs[i].Asset2, diffs[i].DiffPercentage)
		// переводим актив из дорогого в дешевый?
		if !isAct {
			isAct = m.State.Act(diffs[i].Asset1, diffs[i].Asset2)
		}
	}

	fmt.Println("Портфель:")
	for asset, value := range m.State.Portfolio.Actives {
		fmt.Printf("%s: %.2f (%.2f RUB)\n", asset, value, value*m.State.Prices[asset])
	}
}

func main() {
	market := NewMarket(
		map[Asset]float32{
			RUB: 1000,
		},
		map[Asset]float32{
			RUB:  1,
			GOLD: 8357,
			USD:  89.1,
			OZON: 4004,
			SBER: 313,
		})

	for i := 0; i < 10; i++ {
		market.SimulateNextStep()
	}

	fmt.Println(market.State)
}
