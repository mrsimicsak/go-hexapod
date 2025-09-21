package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/joystick"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	pca9685 := i2c.NewPCA9685Driver(r)

	joystickAdaptor := joystick.NewAdaptor()
	stick := joystick.NewDriver(joystickAdaptor, "dualshock4")

	FrontLeft := Leg{Servos: [3]*gpio.ServoDriver{gpio.NewServoDriver(pca9685, "0"), gpio.NewServoDriver(pca9685, "1"), gpio.NewServoDriver(pca9685, "2")}}

	work := func() {
		pca9685.SetPWMFreq(60)

		stick.On(joystick.LeftY, func(data interface{}) {
			rawData := data.(int16)

			scaledData := float32(rawData) / 32768.0

			// scaledData2 := uint8((() * 255) + 128)

			scaledData2 := (scaledData * 90) + 90

			fmt.Println("Raw: ", rawData, " Scaled: ", scaledData2)

			FrontLeft.Servos[0].Move(uint8(scaledData2))
		})
	}

	robot := gobot.NewRobot("Hexapod",
		[]gobot.Connection{r, joystickAdaptor},
		[]gobot.Device{pca9685, FrontLeft.Servos[0], FrontLeft.Servos[1], FrontLeft.Servos[2], stick},
		work,
	)

	robot.Start()

}
