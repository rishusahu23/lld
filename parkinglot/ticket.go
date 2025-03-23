package main

import (
	"errors"
	"time"
)

type Ticket struct {
	Id            string
	ParkingSpotId string
	VehicleId     string
	EntryTime     time.Time
	ExitTime      time.Time
	Charge        float64
}

func NewTicket(id, spotId string, vehicleId string) *Ticket {
	return &Ticket{
		Id:            id,
		ParkingSpotId: spotId,
		VehicleId:     vehicleId,
		EntryTime:     time.Now(),
	}
}

type TicketSvc struct {
	Tickets         map[string]*Ticket
	PaymentStrategy PaymentStrategy
}

func NewTicketSvc(strategy PaymentStrategy) *TicketSvc {
	return &TicketSvc{
		Tickets:         make(map[string]*Ticket),
		PaymentStrategy: strategy,
	}
}

func (t *TicketSvc) CreateTicket(vehicleId, spotId string) *Ticket {
	id := getId()
	ticket := NewTicket(id, spotId, spotId)
	t.Tickets[id] = ticket
	return ticket
}

func (t *TicketSvc) ExitVehicle(id string) (*Ticket, error) {
	ticket, ok := t.Tickets[id]
	if !ok {
		return nil, errors.New("vehicle not found")
	}
	ticket.ExitTime = time.Now()
	charge := t.PaymentStrategy.Calculate(ticket)
	ticket.Charge = charge
	return ticket, nil
}

type PaymentStrategy interface {
	Calculate(ticket *Ticket) float64
}

type DefaultPaymentStrategy struct{}

func (d *DefaultPaymentStrategy) Calculate(ticket *Ticket) float64 {
	dur := ticket.ExitTime.Sub(ticket.EntryTime)
	return dur.Seconds() * 5
}
