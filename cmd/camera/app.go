package main

import (
	"camera/pkg/capture"
	"camera/pkg/sensor"
	"camera/pkg/transport"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nenavizhuleto/horizon/producer"
	"github.com/nenavizhuleto/horizon/protocol"
)

type App struct {
	config Config

	camera *sensor.Camera
	source string

	// durations
	max_duration       time.Duration
	threshold_duration time.Duration

	// transports
	detection_transport *transport.DetectionTransport
	frame_transport     *transport.FrameTransport

	// producers
	producer producer.MessageProducer
}

func NewApp(config Config) (*App, error) {
	source := config.Camera.Source

	if source == "" {
		return nil, fmt.Errorf("source is required")
	}

	max_duration, err := time.ParseDuration(config.Camera.MaxDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to parse duration: %v", err)
	}

	threshold_duration, err := time.ParseDuration(config.Camera.Threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to parse duration: %v", err)
	}

	camera := sensor.NewCamera(
		config.Camera.Name,
		config.Camera.Source,
		config.Camera.Location,
	)

	log.Println("initialized camera with: ", camera)

	frame_transport, err := transport.NewFrameTransport(
		config.Camera.Name,
		[]string{config.Kafka.Broker},
		config.Kafka.Frame,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize frame transport: %v", err)
	}

	detection_transport, err := transport.NewDetectionTransport(
		config.Camera.Name,
		[]string{config.Kafka.Broker},
		config.Kafka.Detection,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize detection transport: %v", err)
	}
	producer := producer.NewMessageProducer(
		camera.ID, camera.Name,
		protocol.ProducerOptions{
			Camera: &protocol.Camera{
				Name:  "camera name goes here",
				Group: "and camera group goes here",
			},
		},
	)

	return &App{
		config: config,

		camera: camera,
		source: source,

		max_duration:       max_duration,
		threshold_duration: threshold_duration,

		detection_transport: detection_transport,
		frame_transport:     frame_transport,

		producer: producer,
	}, nil
}

func (a *App) ExceededDuration(last_detection time.Time) bool {
	return time.Since(last_detection) > a.max_duration
}

func (a *App) ExceededThreshold(last_frame time.Time) bool {
	return time.Since(last_frame) > a.threshold_duration
}

func (a *App) ShouldSendDetection(last_detection, last_frame time.Time) bool {
	return a.ExceededDuration(last_detection) || a.ExceededThreshold(last_frame)
}

func (a *App) Loop() (error, bool) {
	var (
		ctx = context.Background()

		last_detection time.Time
		last_frame     time.Time
	)

	log.Println("connecting to camera")

	results, err := capture.CaptureCamera(ctx, a.camera)
	if err != nil {
		return fmt.Errorf("failed to capture camera: %v", err), true
	}

	log.Println("connected")

	for result := range results {
		if a.ShouldSendDetection(last_detection, last_frame) {
			motions := make([]protocol.Motion, 0)

			for _, d := range result.Detections {
				motion := protocol.Motion{
					Source: a.camera.GetSource(),
					Position: protocol.Position{
						X:      d.Min.X,
						Y:      d.Min.Y,
						Width:  d.Dx(),
						Height: d.Dy(),
					},
				}
				motions = append(motions, motion)
			}

			msg := a.producer.NewMotionDetectionMessage(result.Timestamp, motions)
			if err := a.detection_transport.Send(msg); err != nil {
				return fmt.Errorf("error sending motion: %v\n", err), true
			}

			last_detection = time.Now()
		}

		msg := a.producer.NewFrameMessage(result.Frame)

		if err := a.frame_transport.Send(msg); err != nil {
			return fmt.Errorf("error sending frame: %v\n", err), true
		}

		last_frame = time.Now()
	}

	return fmt.Errorf("results channel closed"), true
}
