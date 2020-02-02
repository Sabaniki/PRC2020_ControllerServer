package main

import (
	"fmt"
	"github.com/nobonobo/joycon"
	"log"
	"math"
	"net"
	"strconv"
)

func adjustment(n float32) float32 {
	if math.Abs(float64(n)) < 0.2 {
		return 0
	} else {
		return n
	}
}

func makeLpower(x float32, y float32) float32 {
	return 100 * (y + x)
}

func makeRpower(x float32, y float32) float32 {
	return 100 * (y - x)
}

func makeArduinoCommand(Lpower float32, Rpower float32, command string) string {
	return "L:" + strconv.Itoa(int(Lpower)) + ",R:" + strconv.Itoa(int(Rpower)) + ",E:" + command
}

func main() {
	devices, err := joycon.Search()
	toRasp, _ := net.Dial("udp", "192.168.100.10:2222")
	toMac, _ := net.Dial("udp", "localhost:3333")
	defer toRasp.Close()
	defer toMac.Close()
	if err != nil {
		log.Fatalln(err)
	}
	if len(devices) == 0 {
		log.Fatalln("joycon not found")
	}
	jcR, err := joycon.NewJoycon(devices[0].Path, false)
	jcL, err := joycon.NewJoycon(devices[1].Path, false)
	if err != nil {
		log.Fatalln(err)
	}
	for true {
		var command string
		stateR := <-jcR.State()
		stateL := <-jcL.State()
		x := adjustment(stateR.RightAdj.X)
		y := adjustment(stateL.LeftAdj.Y)
		if y < 0 {
			x *= -1
		}
		//fmt.Printf("L: %#v, ", 100*(y+x))
		//fmt.Printf("R: %#v\n", 100*(y-x))
		//fmt.Printf("L_Button: %#v\n", stateL.Buttons) // Button bits
		//fmt.Printf("R_Button: %#v\n", stateR.Buttons) // Button bits
		if stateR.Buttons == 0x08 {
			//fmt.Print("wh")
			_, _ = toMac.Write([]byte("wh"))
		} else if stateR.Buttons == 0x01 {
			// toRasp
			command = "pm"
		}
		if stateL.Buttons == 0x80000 {
			// toRasp
			command = "dc"
		}
		//_, _ = toRasp.Write([]byte(makeArduinoCommand(makeLpower(x,y), makeRpower(x,y), command)))
		command =makeArduinoCommand(makeLpower(x, y), makeRpower(x, y), command)
		fmt.Print(command)
		_, _  = toRasp.Write([]byte(command))
	}
	jcR.Close()
	jcL.Close()
}
