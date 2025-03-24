package main

import (
	"fmt"
	"sync"
	"time"
)

type SignalState interface {
	ChangeSignal(manager *TrafficSignalManager)
	GetName() string
}

type RedState struct{}

func (r *RedState) ChangeSignal(manager *TrafficSignalManager) {
	fmt.Println("changing red to green")
	manager.SetState(&GreenState{})
}

func (r *RedState) GetName() string {
	return "red"
}

type GreenState struct{}

func (g *GreenState) ChangeSignal(manager *TrafficSignalManager) {
	fmt.Println("Changing from GREEN to YELLOW")
	manager.SetState(&YellowState{})
}

func (g *GreenState) GetName() string {
	return "Green"
}

type YellowState struct{}

func (y *YellowState) ChangeSignal(manager *TrafficSignalManager) {
	fmt.Println("Changing from YELLOW to RED")
	manager.SetState(&RedState{})
}

func (y *YellowState) GetName() string {
	return "Yellow"
}

type TrafficSignalManager struct {
	state   SignalState
	dur     map[string]time.Duration
	eventCh chan string
	mu      sync.Mutex
}

func NewTrafficSignalManager(eventCh chan string) *TrafficSignalManager {
	return &TrafficSignalManager{
		eventCh: eventCh,
		state:   &RedState{},
		dur: map[string]time.Duration{
			"Red":    15 * time.Second,
			"Green":  15 * time.Second,
			"Yellow": 15 * time.Second,
		},
	}
}

func (m *TrafficSignalManager) SetState(state SignalState) {
	m.state = state
}

func (m *TrafficSignalManager) Start() {
	for {
		select {
		case event := <-m.eventCh:
			if event == "EMERGENCY" {
				fmt.Println("Emergency detected")
				m.SetState(&GreenState{})
				time.Sleep(m.dur["Green"])
				fmt.Println("resuming normal traffic")
			}
		default:
			fmt.Printf("current signal: %s\n", m.state.GetName())
			time.Sleep(m.dur[m.state.GetName()])
			m.state.ChangeSignal(m)
		}
	}
}

type EmergencyHandler struct {
	eventCh chan string
}

func NewEmergencyHandler(eventCh chan string) *EmergencyHandler {
	return &EmergencyHandler{
		eventCh: eventCh,
	}
}

func (e *EmergencyHandler) Start() {
	fmt.Println("starting emergency handler")
	e.eventCh <- "EMERGENCY"
}

func main() {
	eventCh := make(chan string)
	manager := NewTrafficSignalManager(eventCh)
	em := NewEmergencyHandler(eventCh)

	go manager.Start()

	time.Sleep(30 * time.Second)
	em.Start()

	time.Sleep(30 * time.Second)
}
