package capture

import (
	"camera/pkg/sensor"
	"context"
	"image"
	"time"

	"gocv.io/x/gocv"
)

type Result struct {
	Timestamp  time.Time
	Frame      []byte
	Detections []image.Rectangle
}

func CaptureCamera(ctx context.Context, camera *sensor.Camera) (chan Result, error) {
	var (
		source    = camera.GetSource()
		threshold = camera.GetThreshold()
		results   = make(chan Result, 128)
	)

	capture, err := gocv.OpenVideoCapture(source.URI)
	if err != nil {
		return nil, err
	}

	width := capture.Get(gocv.VideoCaptureFrameWidth)
	height := capture.Get(gocv.VideoCaptureFrameHeight)
	camera.SetDimensions(int(width), int(height))

	go func() {
		var (
			frame    = gocv.NewMat()
			roi      = gocv.NewMat()
			detector = NewDetector(threshold)
		)

		defer frame.Close()
		defer detector.Close()
		defer capture.Close()
		defer func() {
			close(results)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if ok := capture.Read(&frame); !ok {
					return
				}

				if frame.Empty() {
					continue
				}
				var (
					timestamp = time.Now()
					resize    = camera.GetResize()
					region    = camera.GetRegion()
				)

				if resize != nil {
					gocv.Resize(frame, &frame, *resize, 0, 0, gocv.InterpolationLinear)
				}

				if region != nil {
					roi = frame.Region(*region)
				} else {
					roi = frame
				}

				detections := detector.Detect(roi)
				if len(detections) == 0 {
					continue
				}

				bytes, err := GetFrameBytes(roi)
				if err != nil {
					continue
				}

				results <- Result{
					Timestamp:  timestamp,
					Frame:      bytes,
					Detections: detections,
				}
			}
		}
	}()

	return results, nil
}

func GetFrameBytes(mat gocv.Mat) ([]byte, error) {
	buf, err := gocv.IMEncodeWithParams(
		gocv.JPEGFileExt,
		mat,
		[]int{
			gocv.IMWriteJpegQuality, 60,
			gocv.IMWriteJpegOptimize, 1,
		},
	)

	if err != nil {
		return nil, err
	}

	defer buf.Close()
	return buf.GetBytes(), nil
}
