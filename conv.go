package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var IP, mask uint32 = 0, 32
	var err error

	if len(os.Args) != 2 {
		fmt.Println("nessun argomento da linea di comando")
		return
	}

	addres := strings.Split(os.Args[1], "/")
	IP, err = getIP(addres[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(addres) > 1 {
		m, err := strconv.Atoi(addres[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		mask = getSbnetMask(uint(m))
	}

	fmt.Println("Net Mask:\t", addToString(mask), "-", addToStringD(mask))
	fmt.Println("IP Addres:\t", addToString(IP), "-", addToStringD(IP))
	fmt.Println("BaseAddres:\t", addToString(IP&mask), "-", addToStringD(IP&mask))
	fmt.Println("BroadCast:\t", addToString(IP|^mask), "-", addToStringD(IP|^mask))
	fmt.Println()
	fmt.Println("Primo IP:\t", addToString((IP&mask)+1), "-", addToStringD((IP&mask)+1))
	fmt.Println("Ultimo IP:\t", addToString((IP|^mask)-1), "-", addToStringD((IP|^mask)-1))

}

func getIP(IP string) (uint32, error) {

	addr := strings.Split(IP, ".")
	if len(addr) != 4 {
		return 0, errors.New("IP non valido")
	}

	var bin uint32
	for _, o := range addr {
		n, err := strconv.Atoi(o)
		if err != nil {
			return 0, errors.New("IP non valido")
		}

		if n > 255 || n < 0 {
			return 0, errors.New("IP non valido")
		}

		bin <<= 8
		bin += uint32(n)
	}

	return bin, nil
}

func getSbnetMask(n uint) uint32 {
	var mask uint32
	mask--
	mask <<= (32 - n)
	return mask
}

func addToString(addr uint32) string {

	var add [4]uint8
	add[0] = uint8(addr)
	add[1] = uint8(addr >> 8)
	add[2] = uint8(addr >> 16)
	add[3] = uint8(addr >> 24)

	return fmt.Sprintf("%08b.%08b.%08b.%08b", add[3], add[2], add[1], add[0])
}

func addToStringD(addr uint32) string {

	var add [4]uint8
	add[0] = uint8(addr)
	add[1] = uint8(addr >> 8)
	add[2] = uint8(addr >> 16)
	add[3] = uint8(addr >> 24)

	return fmt.Sprintf("%d.%d.%d.%d", add[3], add[2], add[1], add[0])
}
