package main

import "errors"

type Vehicle struct {
	Id          string
	Name        string
	VehicleType VehicleType
	RegNumber   string
}

func NewVehicle(id string, name, regNo string, vt VehicleType) *Vehicle {
	return &Vehicle{
		Id:          id,
		Name:        name,
		VehicleType: vt,
		RegNumber:   regNo,
	}
}

type VehicleSvc struct {
	Vehicles map[string]*Vehicle
}

func NewVehicleSvc() *VehicleSvc {
	return &VehicleSvc{
		Vehicles: make(map[string]*Vehicle),
	}
}

func (v *VehicleSvc) GetVehicle(id string) (*Vehicle, error) {
	vehicle, ok := v.Vehicles[id]
	if !ok {
		return nil, errors.New("vehicle not found")
	}
	return vehicle, nil
}

func (v *VehicleSvc) CreateVehicle(vehicle *Vehicle) {
	v.Vehicles[vehicle.Id] = vehicle
}
