package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	pca9685 := i2c.NewPCA9685Driver(r)

	FrontLeft := Leg{Servos: [3]*gpio.ServoDriver{gpio.NewServoDriver(pca9685, "0"), gpio.NewServoDriver(pca9685, "1"), gpio.NewServoDriver(pca9685, "2")}}

	work := func() {
		pca9685.SetPWMFreq(60)

		for i := 0; i < 10; i += 1 {
			fmt.Println("Moving", i)
			angle0, angle1, angle2 := calcLegServoAngles(Vector3{0, 0, float64(i)}, 0, false)
			FrontLeft.Servos[0].Move(uint8(angle0))
			FrontLeft.Servos[1].Move(uint8(angle1))
			FrontLeft.Servos[2].Move(uint8(angle2))
			time.Sleep(1 * time.Second)
		}

		for i := 10; i > 0; i -= 1 {
			fmt.Println("Moving", i)
			angle0, angle1, angle2 := calcLegServoAngles(Vector3{0, 0, float64(i)}, 0, false)
			FrontLeft.Servos[0].Move(uint8(angle0))
			FrontLeft.Servos[1].Move(uint8(angle1))
			FrontLeft.Servos[2].Move(uint8(angle2))
			time.Sleep(1 * time.Second)
		}
	}

	robot := gobot.NewRobot("Hexapod",
		[]gobot.Connection{r},
		[]gobot.Device{pca9685, FrontLeft.Servos[0], FrontLeft.Servos[1], FrontLeft.Servos[2]},
		work,
	)

	robot.Start()

}
