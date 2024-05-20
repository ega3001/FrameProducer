package main

import (
	"camera/pkg/sensor"
	"image"
	"log"

	"github.com/nenavizhuleto/sadm"
)

func ServeSadm(camera *sensor.Camera, port string) {
	s := sadm.New("camera")

	s.AddCommand("region", "set camera region", func(c *sadm.Connection) error {
		var region image.Rectangle

		if err := c.Scanf("%d", &region.Min.X); err != nil {
			return err
		}
		if err := c.Scanf("%d", &region.Min.Y); err != nil {
			return err
		}
		if err := c.Scanf("%d", &region.Max.X); err != nil {
			return err
		}
		if err := c.Scanf("%d", &region.Max.Y); err != nil {
			return err
		}

		c.Println(region)

		camera.SetCaptureRegion(&region)

		return c.Println("camera region set successfully")
	})

	s.AddCommand("resize", "set camera size", func(c *sadm.Connection) error {
		var size image.Point

		if err := c.Scanf("%dx%d", &size.X, &size.Y); err != nil {
			return err
		}

		camera.SetCaptureResize(&size)
		return c.Println("camera size set successfully")
	})

	s.AddCommand("info", "print camera info", func(c *sadm.Connection) error {
		return c.Println(camera.String())
	})

	s.AddCommand("reset", "reset camera setting", func(c *sadm.Connection) error {
		if err := camera.Reset(); err != nil {
			return err
		}
		return c.Println("camera setting reseted")
	})

	for {
		log.Println("serving SAdm at", port)
		if err := s.Listen(port); err != nil {
			log.Println("SADM:", err)
		}
	}
}
