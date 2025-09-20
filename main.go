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
	leg1_1 := gpio.NewServoDriver(pca9685, "0")
    leg1_2 := gpio.NewServoDriver(pca9685, "1")
    leg1_3 := gpio.NewServoDriver(pca9685, "2")

	work := func() {
		pca9685.SetPWMFreq(60)

		for i := 10; i < 150; i += 10 {
			fmt.Println("Turning", i)
			leg1_1.Move(uint8(i))
            leg1_2.Move(uint8(i))
            leg1_3.Move(uint8(i))
			time.Sleep(1 * time.Second)
		}

		for i := 150; i > 10; i -= 10 {
			fmt.Println("Turning", i)
            leg1_1.Move(uint8(i))
            leg1_2.Move(uint8(i))
            leg1_3.Move(uint8(i))
			time.Sleep(1 * time.Second)
		}
	}

	robot := gobot.NewRobot("servoBot",
		[]gobot.Connection{r},
		[]gobot.Device{pca9685, leg1_1, leg1_2, leg1_3},
		work,
	)

	robot.Start()
}
