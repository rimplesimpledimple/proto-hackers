package nogptallowed

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
)

func PrimeTime() {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 47777})
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			panic(err)
		}

		go func() {
			defer conn.Close()
			reader := bufio.NewReader(conn)
			for {
				st, err := reader.ReadString('\n')
				if err != nil {
					return
				}

				fmt.Printf("%s\n", st)
				if isMalformedReq(st) {
					conn.Write([]byte("malformed request"))
					conn.Close()
					return
				}
				isPrimeN := processReq(st)

				// Create a map to represent the JSON object
				data := map[string]interface{}{
					"method": "isPrime",
					"prime":  isPrimeN,
				}

				// Marshal the map into a JSON byte slice
				jsonData, err := json.Marshal(data)
				if err != nil {
					log.Fatalf("Failed to marshal JSON: %v", err)
				}

				log.Printf("Resp: %s\n", string(jsonData))

				// Send the JSON data, followed by a newline
				if _, err := conn.Write(append(jsonData, '\n')); err != nil {
					log.Fatalf("Failed to send data: %v", err)
				}
			}
		}()
	}
}

func isMalformedReq(data string) bool {
	var obj map[string]interface{}

	log.Printf(data)

	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		log.Printf(err.Error())
		return true
	}

	if _, ok := obj["method"]; !ok {
		return true
	}

	methodVal, ok := obj["method"].(string)
	if !ok || methodVal != "isPrime" {
		return true
	}

	if _, ok := obj["number"]; !ok {
		return true
	}

	_, ok = obj["number"].(float64)
	if !ok {
		log.Printf("not float64: %+v", obj)
		return true
	}

	return false
}

func processReq(data string) bool {
	var obj map[string]interface{}

	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		log.Printf(err.Error())
		return true
	}

	num, _ := obj["number"].(float64)
	if num != float64(int64(num)) {
		return false
	}

	return isPrime(int64(num))
}

func isPrime(num int64) bool {
	if num <= 1 {
		return false
	}

	floatNum := float64(num)
	for i := int64(2); i <= int64(math.Sqrt(floatNum)); i++ {
		if num%i == 0 {
			log.Printf("Num: %d is not prime", num)
			return false
		}
	}

	log.Printf("Num: %d is prime", num)
	return true
}
