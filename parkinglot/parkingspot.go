package main

import "errors"

type ParkingSpot struct {
	Id        string
	Level     int
	Available bool
	SpotType  SpotType
}

func (p *ParkingSpot) IsAvailable() bool {
	return p.Available
}

func (p *ParkingSpot) SetAvailable(available bool) {
	p.Available = available
}

func NewParkingSpot(id string, level int, spotType SpotType) *ParkingSpot {
	return &ParkingSpot{
		Id:        id,
		Level:     level,
		Available: true,
		SpotType:  spotType,
	}
}

type ParkingSpotSvc struct {
	ParkingSpots map[int]map[string]*ParkingSpot
}

func NewParkingSpotSvc() *ParkingSpotSvc {
	return &ParkingSpotSvc{
		ParkingSpots: make(map[int]map[string]*ParkingSpot),
	}
}

func (p *ParkingSpotSvc) GetAvailableParkingSpot(spotType SpotType) (*ParkingSpot, error) {
	for _, level := range p.ParkingSpots {
		for _, parkingSpot := range level {
			if parkingSpot.IsAvailable() && parkingSpot.SpotType == spotType {
				return parkingSpot, nil
			}
		}
	}
	return nil, errors.New("parking spot not found")
}

func (p *ParkingSpotSvc) SetAvailability(id string, available bool) {
	for _, level := range p.ParkingSpots {
		for _, parkingSpot := range level {
			if parkingSpot.Id == id {
				parkingSpot.SetAvailable(available)
				return
			}
		}
	}
}

func (p *ParkingSpotSvc) CreateParkingSpot(parkingSpot *ParkingSpot) error {
	levels, ok := p.ParkingSpots[parkingSpot.Level]
	if !ok {
		levels = make(map[string]*ParkingSpot)
	}
	levels[parkingSpot.Id] = parkingSpot
	p.ParkingSpots[parkingSpot.Level] = levels
	return nil
}
