package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//complile: go build -ldflags "-H windowsgui"

func main() {
	var out string

	c, err := os.ReadFile(os.Args[1])
	if err != nil {
		return
	}

	cnt := strings.Split(string(c), "\n")

	for _, e := range cnt {
		e = strings.TrimSpace(e)
		inf, err := getInfo(e)
		if err != nil {
			out += e + " - " + err.Error() + "\n\n"
		} else {
			out += e + inf + "\n\n"
		}
	}

	f, err := os.Create("a.txt")
	if err != nil {
		return
	}
	f.Write([]byte(out))

}

func getInfo(ind string) (string, error) {
	var IP, mask uint32 = 0, 32
	var cidrM int
	var err error

	addres := strings.Split(ind, "/")
	IP, err = getIP(addres[0])
	if err != nil {
		return "", err
	}

	if len(addres) > 1 {
		cidrM, err = strconv.Atoi(addres[1])

		if cidrM < 0 && cidrM > -32 {
			cidrM = 32 + cidrM
		}

		if err != nil || cidrM < 0 || cidrM > 30 {
			return "", errors.New("maschera non valida")
		}

		mask = getSbnetMask(uint(cidrM))
	}

	var info string
	if mask == 0xFFFFFFFC {
		info = fmt.Sprint(
			"\nBaseAddres:\t", addToStringD(IP&mask), "/", cidrM,
			"\nBroadCast:\t", addToStringD(IP|^mask),
			"\nGetaway:\t-",
			"\nPrimo IP:\t", addToStringD((IP&mask)+1),
			"\nSecondo IP:\t", addToStringD((IP|^mask)-1),
			"\nNet Mask:\t", addToStringD(mask),
			"\nWildcard:\t", addToStringD(^mask),
		)
	} else {
		info = fmt.Sprint(
			"\nBaseAddres:\t", addToStringD(IP&mask), "/", cidrM,
			"\nBroadCast:\t", addToStringD(IP|^mask),
			"\nGetaway:\t", addToStringD((IP|^mask)-1),
			"\nPrimo IP:\t", addToStringD((IP&mask)+1),
			"\nUltimo IP:\t", addToStringD((IP|^mask)-2),
			"\nNet Mask:\t", addToStringD(mask),
			"\nWildcard:\t", addToStringD(^mask),
		)
	}

	return info, nil
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

func addToStringD(addr uint32) string {

	var add [4]uint8
	add[0] = uint8(addr)
	add[1] = uint8(addr >> 8)
	add[2] = uint8(addr >> 16)
	add[3] = uint8(addr >> 24)

	return fmt.Sprintf("_ %d.%d.%d.%d", add[3], add[2], add[1], add[0])
}
