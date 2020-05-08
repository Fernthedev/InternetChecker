package network

import (
	"fmt"
	"net"
	"strings"
)

type Iface struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	StatusBool bool `json:"status_bool"`
	Multicast bool   `json:"multicast"`
	Broadcast bool   `json:"broadcast"`
}

func GetInterfaces() []Iface {
	interfaces, err := net.Interfaces()

	if err != nil {
		fmt.Println(err)
		return nil
	}

	var ifaceset []Iface
	var ifc Iface

	for _, i := range interfaces {
		ifc.Name = i.Name
		if strings.Contains(i.Flags.String(), "up") {
			ifc.Status = "UP"
			ifc.StatusBool = true
		} else {
			ifc.Status = "DOWN"
			ifc.StatusBool = false
		}
		if strings.Contains(i.Flags.String(), "multicast") {
			ifc.Multicast = true
		} else {
			ifc.Multicast = false
		}
		if strings.Contains(i.Flags.String(), "broadcast") {
			ifc.Broadcast = true
		} else {
			ifc.Broadcast = false
		}


		ifaceset = append(ifaceset, ifc)
	}

	return ifaceset
}