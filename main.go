package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/tarm/serial"
	"go.bug.st/serial/enumerator"
)

var paths []string
var port *serial.Port

const Sep = "\r\n"

// var (
// 	ErrWriteFailed = errors.New("at: write failed")
// )

func main() {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		time.Sleep(10 * time.Second)
		return
	}
	for _, port := range ports {
		fmt.Printf("Found port: %s\n", port.Name)
		if port.IsUSB {
			fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
			fmt.Printf("   USB serial %s\n", port.SerialNumber)
		}
		paths = append(paths, port.Name)
	}

	var qs = []*survey.Question{
		{
			Name: "path",
			Prompt: &survey.Select{
				Message: "Choose a path:",
				Options: paths,
			},
		},
		{
			Name:   "baudRate",
			Prompt: &survey.Input{Message: "input a baudRate", Default: "115200"},
		},
	}

	answers := struct {
		Path     string // survey will match the question and field names
		BaudRate int    // survey will match the question and field names
	}{}

	err = survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	config := &serial.Config{Name: answers.Path, Baud: answers.BaudRate, ReadTimeout: time.Second * 2}
	s, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	port = s

	// mode := &serial.Mode{
	// 	BaudRate: answers.BaudRate,
	// }
	// selectPort, err := serial.Open(answers.Path, mode)
	// port = selectPort
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("main: Port opened")

	imputer("AT")
}

func checkAt(at string) string {
	if at == "exit" {
		os.Exit(3)
	}
	if at == "text_mode" {
		at = "AT+CMGF=0"
	}
	if at == "imsi" {
		at = "AT+CIMI"
	}
	if at == "iccid" {
		at = "AT+ICCID"
	}
	if at == "manufacturer" {
		at = "AT+CGMI"
	}
	if at == "model" {
		at = "AT+CGMI"
	}
	if at == "version" {
		at = "AT+CGMR"
	}
	if at == "imei" {
		at = "AT+CGSN"
	}
	if at == "messages" {
		at = "AT+CMGL=4"
	}
	if at == "signal" {
		at = "AT+CSQ"
	}
	if at == "network" {
		at = "AT+CREG?"
	}
	if at == "number" {
		at = "AT+CNUM"
	}

	return at
}

func imputer(at string) {
	if at == "" {
		var qs = []*survey.Question{
			{
				Name: "at",
				Prompt: &survey.Input{
					Message: "input a AT command",
					Help:    "Binded: text_mode, imsi, iccid, manufacturer, model, version, imei, messages, signal, network, number, exit",
				},
			},
		}
		answers := struct {
			At string // survey will match the question and field names
		}{}
		var err = survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		at = checkAt(answers.At)

	}

	fmt.Printf(">>> %v\n", at)
	var err error
	_, err = port.Write([]byte(at + Sep))
	if err != nil {
		log.Fatal(err)
	}

	buff := make([]byte, 100)
	for {
		n, _ := port.Read(buff)
		// if n == 0 {
		// 	fmt.Println("\nEOF")
		// 	break
		// }
		text := strings.TrimSpace(string(buff[:n]))
		if len(text) < 1 {
			break
		}
		if strings.HasPrefix(at, text) {
			continue
		}
		fmt.Printf("<<< %v\n", text)

		// fmt.Printf("%v", string(buff[:n]))

		// buf := bufio.NewReader(port)
		// if line, err = buf.ReadString('\r'); err != nil {
		// 	continue
		// }
		// text := strings.TrimSpace(line)
		// fmt.Printf("TEXT: %v\n", text)
		// if !strings.HasPrefix("AT", text) {
		// 	return
		// }

		// fmt.Println(">> reading")
		// buff := make([]byte, 1000)
		// n, err := port.Read(buff)
		// fmt.Printf(">> reading completed n=%d err=%+v\n", n, err)
		// fmt.Printf("%v", string(buff[:n]))

	}

	imputer("")
}
