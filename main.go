package main

import (
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

		var targetPos Vector3

		stick.On(joystick.LeftX, func(data interface{}) {
			rawData := data.(int16)

			scaledData := float32(rawData) / 32768.0

			scaledData2 := scaledData * 50.0

			targetPos.X = float64(scaledData2)

			angle0, angle1, angle2 := calcLegServoAngles(targetPos, 0, false)
			FrontLeft.Servos[0].Move(uint8(angle0))
			FrontLeft.Servos[1].Move(uint8(angle1))
			FrontLeft.Servos[2].Move(uint8(angle2))
		})

		stick.On(joystick.LeftY, func(data interface{}) {
			rawData := data.(int16)

			scaledData := float32(rawData) / 32768.0

			scaledData2 := scaledData * 100.0

			targetPos.Y = float64(scaledData2)

			angle0, angle1, angle2 := calcLegServoAngles(targetPos, 0, false)
			FrontLeft.Servos[0].Move(uint8(angle0))
			FrontLeft.Servos[1].Move(uint8(angle1))
			FrontLeft.Servos[2].Move(uint8(angle2))
		})
	}

	robot := gobot.NewRobot("Hexapod",
		[]gobot.Connection{r, joystickAdaptor},
		[]gobot.Device{pca9685, FrontLeft.Servos[0], FrontLeft.Servos[1], FrontLeft.Servos[2], stick},
		work,
	)

	robot.Start()

}
