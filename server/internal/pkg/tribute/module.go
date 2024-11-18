package tribute

import "go.uber.org/fx"

var Module = fx.Module("TributeService", fx.Provide(fx.Annotate(initTributeList, fx.As(new(TributeService)))))

type TributeService interface {
	List() []Tribute
	Get(id uint16) (Tribute, bool)
	Add(Tribute)
}

type TributeList struct {
	id       uint16
	tributes []Tribute
}

type Tribute struct {
	Id     uint16 `json:"id"`
	Name   string `json:"name"`
	Amount int    `json:"amount"` //TODO: Convert this into balance
}

func (tl *TributeList) List() []Tribute {
	return tl.tributes
}

func (tl *TributeList) Get(id uint16) (Tribute, bool) {
	for _, tribute := range tl.tributes {
		if tribute.Id == id {
			return tribute, true
		}
	}
	return Tribute{}, true
}

func (tl *TributeList) Add(tribute Tribute) {
	tl.id++
	tribute.Id = tl.id
	tl.tributes = append(tl.tributes, tribute)
}

func initTributeList() *TributeList {
	return &TributeList{
		id:       77,
		tributes: make([]Tribute, 0, 7), // fill with Cato, Pompey, Scipio, Hector, Nero, Phillis, Juno, Cyrus,
	}
}
