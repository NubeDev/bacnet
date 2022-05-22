package ndpu

type NetworkMessageType uint8

//go:generate stringer -type=NetworkMessageType
const (
	WhoIsRouterToNetwork NetworkMessageType = 0x00
	WhatIsNetworkNumber  NetworkMessageType = 0x12
)
