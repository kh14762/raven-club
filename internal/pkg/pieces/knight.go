package pieces

var Module = ("") // Create DI Module here, Names as Knight

type Knight struct {
	raven Raven
}

type Bishop struct {
	kite Kite
}

type Rook struct {
	heron Heron
}

type Pawn struct {
	finch Finch
}

type Raven struct {
	munitionRack []Element
	composite    Heirloom
}

type Kite struct {
	weatherRadar Element
	composite    Heirloom
}

type Heron struct {
	composite Heirloom
}

type Finch struct {
	composite Heirloom
}

type Motor struct{}

type Heirloom struct {
	title    string
	eID      string
	material string

	batteries        Element
	dcDistribution   Element
	txRadioAltimeter []Element
	markerBeacon     Element

	motor Motor
}

type Element struct{}
