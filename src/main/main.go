package main

import (
	"fmt"
	"github.com/nobonobo/joycon"
	"log"
)

func main() {
	devices, err := joycon.Search()
	if err != nil {
		log.Fatalln(err)
	}
	if len(devices) == 0 {
		log.Fatalln("joycon not found")
	}
	jc, err := joycon.NewJoycon(devices[0].Path, false)
	if err != nil {
		log.Fatalln(err)
	}
	for true {
		s := <-jc.State()
		fmt.Printf("%#v\n", s.Buttons)  // Button bits
		fmt.Printf("%#v\n", s.RightAdj) // Right Analog Stick State
	}
	//a := <-jc.Sensor()
	//fmt.Printf("%#v\n", a.Accel) // Acceleration Sensor State
	//fmt.Printf("%#v\n", a.Gyro)  // Gyro Sensor State

	jc.Close()
}
