package main

import (
	"fmt"
	"math/rand"
	"time"
)

func getId() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), 100*rand.Int())
}

type ParkingLot struct {
	ParkingSpotSvc *ParkingSpotSvc
	TicketSvc      *TicketSvc
	VehicleSvc     *VehicleSvc
}

func NewParkingLot(svc *ParkingSpotSvc, ticketSvc *TicketSvc,
	vehicleSvc *VehicleSvc) *ParkingLot {
	return &ParkingLot{
		ParkingSpotSvc: svc,
		TicketSvc:      ticketSvc,
		VehicleSvc:     vehicleSvc,
	}
}

func getSpotType(vehicleType VehicleType) SpotType {
	switch vehicleType {
	case VehicleTypeTwoWheeler:
		return SpotTypeTwoWheeler
	default:
		return SpotTypeFourWheeler
	}
}

func (p *ParkingLot) Entry(vehicleId string) (*Ticket, error) {
	vehicle, err := p.VehicleSvc.GetVehicle(vehicleId)
	if err != nil {
		return nil, err
	}
	spot, err := p.ParkingSpotSvc.GetAvailableParkingSpot(getSpotType(vehicle.VehicleType))
	if err != nil {
		return nil, err
	}
	p.ParkingSpotSvc.SetAvailability(spot.Id, false)
	ticket := p.TicketSvc.CreateTicket(vehicleId, spot.Id)
	return ticket, nil
}

func (p *ParkingLot) Exit(ticketId string) (float64, error) {
	ticket, err := p.TicketSvc.ExitVehicle(ticketId)
	if err != nil {
		return 0, err
	}
	p.ParkingSpotSvc.SetAvailability(ticket.ParkingSpotId, true)
	return ticket.Charge, nil
}

func main() {
	parkingSpotSvc := NewParkingSpotSvc()
	ticketSvc := NewTicketSvc(&DefaultPaymentStrategy{})
	vehicleSvc := NewVehicleSvc()

	pl := NewParkingLot(parkingSpotSvc, ticketSvc, vehicleSvc)

	spotTwoWheeler := NewParkingSpot(getId(), 1, SpotTypeTwoWheeler)
	spotFourWheeler := NewParkingSpot(getId(), 1, SpotTypeFourWheeler)

	_ = parkingSpotSvc.CreateParkingSpot(spotTwoWheeler)
	_ = parkingSpotSvc.CreateParkingSpot(spotFourWheeler)

	vehicleTwoWheeler := NewVehicle(getId(), "bike", "reg-1", VehicleTypeTwoWheeler)
	vehicleFourWheeler := NewVehicle(getId(), "car", "reg-2", VehicleTypeFourWheeler)

	vehicleSvc.CreateVehicle(vehicleTwoWheeler)
	vehicleSvc.CreateVehicle(vehicleFourWheeler)

	ticket, err := pl.Entry(vehicleTwoWheeler.Id)
	if err != nil {
		panic(err)
	}
	charge, err := pl.Exit(ticket.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ticket charge: %f\n", charge)
}
