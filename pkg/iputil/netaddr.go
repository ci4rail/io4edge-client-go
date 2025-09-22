package iputil

import (
	"errors"
	"strconv"
	"strings"
)

// NetAddressSplit splits addr to host and port
// example: addr="myhost.example.com:1234" -> host="myhost.example.com", port=1234
func NetAddressSplit(addr string) (host string, port int, err error) {
	fields := strings.Split(addr, ":")
	if len(fields) != 2 {
		return "", 0, errors.New("invalid address " + addr)
	}
	port, err = strconv.Atoi(fields[1])
	if err != nil {
		return "", 0, errors.New("invalid port in " + addr)
	}
	return fields[0], port, nil
}
