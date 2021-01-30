package network

import (
	"fmt"
	"math"
	"math/bits"
	"strconv"
	"strings"
)

// Mask is just a tidy way of storing a netmask
type Mask struct {
	Decimal uint8
	Dotted  []byte
}

// NetworkInfo contains the result of the CalculateNetwork function
type NetworkInfo struct {
	Address        []uint8
	Netmask        Mask
	Wildcard       []uint8
	Network        []uint8
	Broadcast      []uint8
	HostMinAddress []uint8
	HostMaxAddress []uint8
	HostsQuantity  uint32
}

// CidrToMask converts a slash netmask to a Mask struct
func CidrToMask(cidr uint8) Mask {
	var maskStruct Mask
	maskStruct.Decimal = cidr
	var byteArr = make([]byte, 0)

	for cidr > 0 {
		if cidr >= 8 {
			byteArr = append(byteArr, uint8(255))
			cidr -= 8
		} else {
			byteArr = append(byteArr, bits.Reverse8(uint8(math.Pow(2, float64(cidr))-1)))
			break
		}
	}

	for len(byteArr) < 4 {
		byteArr = append(byteArr, uint8(0))
	}

	maskStruct.Dotted = byteArr

	return maskStruct
}

// DottedToMask converts a dotted netmask to a Mask struct
func DottedToMask(dotted []byte) Mask {
	var cidr uint8 = 0

	for i := 0; i < len(dotted); i++ {
		cidr += uint8(bits.Len8(dotted[i]))
	}

	return Mask{cidr, dotted}
}

func strToByteArr(str string) []byte {
	var byteArr []byte = make([]byte, 0)

	for _, bytePart := range strings.Split(str, ".") {
		bytePartInt, err := strconv.Atoi(bytePart)
		if err != nil {
			return make([]byte, 0)
		}
		bytePartUint8 := uint8(bytePartInt)
		byteArr = append(byteArr, bytePartUint8)
	}

	return byteArr
}

func ByteArrToStr(byteArr []byte) string {
	str := "" + fmt.Sprint(byteArr[0])
	for i := 1; i < len(byteArr); i++ {
		str += "." + fmt.Sprint(byteArr[i])
	}

	return str
}

// CalculateNetwork calculates all the infos of a given network
func CalculateNetwork(ip string, subnet string) NetworkInfo {
	var networkInfoStruct NetworkInfo

	// Convert the ip string to a byte array
	networkInfoStruct.Address = strToByteArr(ip)

	// Convert the subnet mask string to a byte array
	networkInfoStruct.Netmask = DottedToMask(strToByteArr(subnet))

	// Calculate network wildcard
	var wildcard []byte = make([]byte, 0)

	for _, subnetPart := range networkInfoStruct.Netmask.Dotted {
		wildcard = append(wildcard, uint8(255)-subnetPart)
	}

	networkInfoStruct.Wildcard = wildcard
	wildcard = nil

	// Calculate network
	var networkArr []byte = make([]byte, 0)

	for i, ipByte := range networkInfoStruct.Address {
		networkArr = append(networkArr, ipByte&networkInfoStruct.Netmask.Dotted[i])
	}

	networkInfoStruct.Network = networkArr
	networkArr = nil

	// Calculate broadcast address
	var broadcastArr []byte = make([]byte, 0)
	for i, networkByte := range networkInfoStruct.Network {
		broadcastArr = append(broadcastArr, networkByte+networkInfoStruct.Wildcard[i])
	}

	networkInfoStruct.Broadcast = broadcastArr
	broadcastArr = nil

	// Calculate minimum host IP address
	networkInfoStruct.HostMinAddress = make([]uint8, 4)
	copy(networkInfoStruct.HostMinAddress, networkInfoStruct.Network)
	networkInfoStruct.HostMinAddress[3]++ // Increment of 1 the host address (so it is network address + 1)

	// Calculate maximum host IP address
	networkInfoStruct.HostMaxAddress = make([]uint8, 4)
	copy(networkInfoStruct.HostMaxAddress, networkInfoStruct.Broadcast)
	networkInfoStruct.HostMaxAddress[3]--

	// Calculate maximum quantity of hosts in the network
	networkInfoStruct.HostsQuantity = uint32(math.Pow(2, float64(32-networkInfoStruct.Netmask.Decimal)) - 2)

	return networkInfoStruct
}
