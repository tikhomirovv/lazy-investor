package services

import (
	"fmt"
	"sort"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/dto"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

type StrategyService struct {
	logger  logging.Logger
	tinkoff *TinkoffService
}

func NewStrategyService(logger logging.Logger, tinkoff *TinkoffService) *StrategyService {
	return &StrategyService{
		logger:  logger,
		tinkoff: tinkoff,
	}
}

// ---------

type Portfolio struct {
	Name    string
	Amount  float64 // RUB
	Actives map[dto.Isin]float64
}

type Market struct {
	Instruments map[dto.Isin]*dto.Instrument
	State       MarketState
}
type MarketState struct {
	Portfolio   *Portfolio
	Alternative map[dto.Isin]*Portfolio
	Prices      map[dto.Isin]float64
	History     map[dto.Isin][]float64
}

type TimePrices struct {
	Time   time.Time
	Prices map[dto.Isin]dto.Candle
}

const GoldIsin = "RU000A101NZ2"

func (m *Market) PrintPortfolio(p *Portfolio, onlyValue bool) {
	if !onlyValue {
		fmt.Printf("💼 %s: RUB - %.2f\n", p.Name, p.Amount)
	}
	for isin, value := range p.Actives {
		name := m.Instruments[isin].Name
		rub := value * m.State.Prices[isin]
		gold := value * m.State.Prices[isin] / m.State.Prices[GoldIsin]
		if onlyValue {
			fmt.Printf("\t %.2f RUB,\t%.2f GOLD |  %s\n", rub, gold, name)
		} else {
			fmt.Printf("\t - %s: %.2f (%.2f RUB, %.2f GOLD)\n", name, value, rub, gold)
			fmt.Println("---")
		}
	}
}
func (p *Portfolio) Buy(isin dto.Isin, price float64) bool {
	if p.Amount <= 0 {
		// fmt.Println("No money in portfolio")
		return false
	}
	value := p.Amount / price
	p.Amount = 0
	p.Actives[isin] = value
	fmt.Printf("💰: %.4f %s by %.4f (%.4f RUB)\n",
		value, isin, price, p.Amount)
	return true
}

func (p *Portfolio) Sale(isin dto.Isin, price float64) bool {
	value, exists := p.Actives[isin]
	if !exists || value <= 0 {
		// fmt.Printf("No active in portfolio `%s`\n", isin)
		return false
	}
	delete(p.Actives, isin)
	p.Amount = value * price
	fmt.Printf("🔻: %.4f %s by %.4f (%.4f RUB)\n",
		value, isin, price, p.Amount)
	return true
}

func (m *Market) GetAmountInRub(p *Portfolio) float64 {
	amount := p.Amount
	for isin, active := range p.Actives {
		amount += active * m.State.Prices[isin]
	}

	return amount
}

func NewPortfolio(name string) *Portfolio {
	return &Portfolio{Name: name, Amount: 1000, Actives: make(map[dto.Isin]float64)}
}

func NewMarket() *Market {
	state := MarketState{
		Portfolio:   NewPortfolio("main"),
		Alternative: make(map[dto.Isin]*Portfolio),
		Prices:      make(map[dto.Isin]float64),
		History:     make(map[dto.Isin][]float64),
	}
	instruments := make(map[dto.Isin]*dto.Instrument)
	return &Market{State: state, Instruments: instruments}
}

func (m *Market) SimulateNextStep(tp TimePrices, prev TimePrices) {
	// fmt.Printf("Time: %s\n", tp.Time)
	// for isin, candle := range tp.Prices {
	// 	fmt.Printf("\t%s -> Open: %.2f,\tClose: %.2f\t| %s\n", isin, candle.Open, candle.Close, m.Instruments[isin].Name)
	// }

	fmt.Println("--------------")
	beforeAfter := map[dto.Isin][3]float64{}
	for isin := range tp.Prices {
		before := prev.Prices[isin].Close
		if before <= 0 {
			continue
		}
		after := tp.Prices[isin].Close
		change := (after - before) / before * 100
		beforeAfter[isin] = [3]float64{before, after, change}
		m.State.Prices[isin] = after
		fmt.Printf("%s: %.2f -> %.2f (%.2f%%) | %s\n", isin, before, after, change, m.Instruments[isin].Name)
	}
	fmt.Println("--------------")

	// Структура для хранения результатов
	type DiffResult struct {
		Active1        dto.Isin
		Active2        dto.Isin
		DiffPercentage float64
	}
	type Step struct {
		diffs []DiffResult
	}
	var step Step
	for isin, change := range beforeAfter {
		for isin2, change2 := range beforeAfter {
			diff := (change[2] - change2[2])
			if diff == 0 || isin == GoldIsin || isin2 == GoldIsin {
				continue
			}
			step.diffs = append(step.diffs, DiffResult{Active1: isin, Active2: isin2, DiffPercentage: diff})
			fmt.Printf("%s / %s:  %.2f%%\n",
				m.Instruments[isin].Name, m.Instruments[isin2].Name, diff)
		}
	}

	// Сортируем по убыванию разницы
	sort.Slice(step.diffs, func(i, j int) bool {
		// return step.diffs[i].DiffPercentage > step.diffs[j].DiffPercentage
		return step.diffs[i].DiffPercentage < step.diffs[j].DiffPercentage
	})

	// 1 Strategy
	isAct := false
	fmt.Println("Изменения:")
	for i := 0; i < len(step.diffs); i++ {
		fmt.Printf("%s / %s: %.2f%%\n", m.Instruments[step.diffs[i].Active1].Name, m.Instruments[step.diffs[i].Active2].Name, step.diffs[i].DiffPercentage)
		// переводим актив из дорогого в дешевый?
		if !isAct && step.diffs[i].DiffPercentage < 0 {
			m.State.Portfolio.Sale(step.diffs[i].Active2, tp.Prices[step.diffs[i].Active2].Close)
			// if isSale {
			isBuy := m.State.Portfolio.Buy(step.diffs[i].Active1, tp.Prices[step.diffs[i].Active1].Close)
			if isBuy {
				isAct = true
				continue
			}
		}
	}

	// 2 Alt
	for i := 0; i < len(step.diffs); i++ {
		isin := step.diffs[i].Active1
		if m.State.Alternative[isin] == nil {
			m.State.Alternative[isin] = NewPortfolio(fmt.Sprintf("Alternative #%d", i))
			pp := m.State.Alternative[isin]
			pp.Amount = m.GetAmountInRub(m.State.Portfolio)
			// // pp.Sale(diffs[i].Active1, tp.Prices[diffs[i].Active1].Close)
			pp.Buy(step.diffs[i].Active1, tp.Prices[step.diffs[i].Active1].Close)
			m.State.Alternative[isin] = pp
		}
	}

	m.PrintPortfolio(m.State.Portfolio, false)
	for _, alt := range m.State.Alternative {
		m.PrintPortfolio(alt, true)
	}
}

func GroupCandlesByTime(candlesByIsin map[dto.Isin][]dto.Candle) []TimePrices {
	timeMap := make(map[time.Time]map[dto.Isin]dto.Candle)
	for isin, candles := range candlesByIsin {
		for _, candle := range candles {
			if _, exists := timeMap[candle.Time]; !exists {
				timeMap[candle.Time] = make(map[dto.Isin]dto.Candle)
			}
			timeMap[candle.Time][isin] = candle
		}
	}
	// Преобразуем map в слайс с сортировкой по времени
	var timePrices []TimePrices
	for t, prices := range timeMap {
		timePrices = append(timePrices, TimePrices{Time: t, Prices: prices})
	}
	sort.Slice(timePrices, func(i, j int) bool {
		return timePrices[i].Time.Before(timePrices[j].Time)
	})
	return timePrices
}

func (ss *StrategyService) Test(instruments []*dto.Instrument) {
	ss.logger.Debug("Instruments", "is", instruments)

	instrs := make(map[dto.Isin]*dto.Instrument)
	candlesByIsin := make(map[dto.Isin][]dto.Candle)
	for _, inst := range instruments {
		instrs[inst.Isin] = inst
		candles, err := ss.tinkoff.GetCandles(inst)
		if err != nil {
			ss.logger.Error("StrategyService.Test: %w", err)
		}
		candlesByIsin[inst.Isin] = candles
	}
	market := NewMarket()
	market.Instruments = instrs
	// market.State.Portfolio.Amount = 1000
	grouped := GroupCandlesByTime(candlesByIsin)
	for i := range grouped {
		if i <= 0 { // skip first
			continue
		}
		market.SimulateNextStep(grouped[i], grouped[i-1])
		fmt.Scanln()
	}
	fmt.Println(market.State.Portfolio)
}
