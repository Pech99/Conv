package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var noGetaway, binario bool = false, false

func main() {
	var IP, mask uint32 = 0, 32
	var err error
	var ind string

	if len(os.Args) <= 1 {
		fmt.Println("nessun argomento da linea di comando")
		return
	}

	for _, arg := range os.Args[1:] {
		if arg == "-g" {
			noGetaway = true
			continue

		} else if arg == "-b" {
			binario = true
			continue

		} else if arg == "--help" {
			println("Fornisce alcune inforazioni riguardo l'indirizzo in notazione puntata (x.y.z.q/m)\nformito da linea di comando. Può anche essere aggiunta una maschera CIDR,\nin assenza della quale verrà usata la suddivisione in classi.\n-g\tnon tiene conto del Getaway.\n-b\tfornisce anche la codifica in binario.")
			return

		} else {
			ind = arg
			break
		}
	}

	addres := strings.Split(ind, "/")
	IP, err = getIP(addres[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(addres) > 1 {
		m, err := strconv.Atoi(addres[1])
		if err != nil || m < 0 || m > 30 {
			fmt.Println("maschera non valida")
			return
		}
		mask = getSbnetMask(uint(m))
	} else {
		if uint8(IP>>28) == 15 { // class E
			fmt.Println("IP di classe E\nRiservato per usi futuri")
			return

		} else if uint8(IP>>28) == 14 { // class D
			fmt.Println("IP di classe D\nIndirizzo Multicast")
			return

		} else if uint8(IP>>29) == 6 { // class C
			fmt.Println("IP di classe C")
			mask = getSbnetMask(uint(24))

		} else if uint8(IP>>30) == 2 { // class B
			fmt.Println("IP di classe B")
			mask = getSbnetMask(uint(16))

		} else if uint8(IP>>31) == 0 { // class A
			mask = getSbnetMask(uint(8))
			fmt.Println("IP di classe A")
		}
	}

	if mask == 0xFFFFFFFC {
		fmt.Println("Connessione punto punto")
		print("BaseAddres:\t", (IP & mask))
		print("BroadCast:\t", (IP | ^mask))
		print("Primo IP:\t", ((IP & mask) + 1))
		print("Secondo IP:\t", ((IP | ^mask) - 1))
		print("Net Mask:\t", (mask))
		print("Wildcard:\t", (^mask))
		return
	}

	if noGetaway {
		print("BaseAddres:\t", (IP & mask))
		print("BroadCast:\t", (IP | ^mask))
		print("Primo IP:\t", ((IP & mask) + 1))
		print("Ultimo IP:\t", ((IP | ^mask) - 1))
		print("Net Mask:\t", (mask))
		print("Wildcard:\t", (^mask))
	} else {
		print("BaseAddres:\t", (IP & mask))
		print("BroadCast:\t", (IP | ^mask))
		print("Getaway:\t", ((IP | ^mask) - 1))
		print("Primo IP:\t", ((IP & mask) + 1))
		print("Ultimo IP:\t", ((IP | ^mask) - 2))
		print("Net Mask:\t", (mask))
		print("Wildcard:\t", (^mask))
	}

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

func print(lable string, num uint32) {
	if binario {
		fmt.Println(lable, addToStringB(num), "-", addToStringD(num))
	} else {
		fmt.Println(lable, addToStringD(num))
	}
}

func addToStringB(addr uint32) string {

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
