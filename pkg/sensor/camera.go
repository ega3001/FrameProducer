package sensor

import (
	"fmt"
	"image"
	"sync"

	"github.com/google/uuid"
	"github.com/nenavizhuleto/horizon/protocol"
)

var (
	DefaultDebug     = false
	DefaultThreshold = 500
)

type Camera struct {
	ID       string
	Name     string
	source   protocol.Source
	location string

	debug     bool
	threshold int

	mx        sync.Mutex
	region    *image.Rectangle
	resize_to *image.Point
}

func NewCamera(name string, source string, location string) *Camera {
	return &Camera{
		ID:   uuid.NewString(),
		Name: name,
		source: protocol.Source{
			URI:        source,
			Dimensions: protocol.Dimensions{},
		},
		location:  location,
		debug:     DefaultDebug,
		threshold: DefaultThreshold,
	}
}

// --- Getters

func (c *Camera) GetSource() protocol.Source {
	c.mx.Lock()
	defer c.mx.Unlock()
	return c.source
}

func (c *Camera) GetDebug() bool {
	c.mx.Lock()
	defer c.mx.Unlock()
	return c.debug
}

func (c *Camera) GetThreshold() int {
	c.mx.Lock()
	defer c.mx.Unlock()
	return c.threshold
}

func (c *Camera) GetResize() *image.Point {
	c.mx.Lock()
	defer c.mx.Unlock()
	return c.resize_to
}

func (c *Camera) GetRegion() *image.Rectangle {
	c.mx.Lock()
	defer c.mx.Unlock()
	return c.region
}

// --- Setters

func (c *Camera) SetDimensions(width, height int) {
	c.source.Dimensions = protocol.Dimensions{
		Width:  width,
		Height: height,
	}
}

func (c *Camera) SetCaptureRegion(region *image.Rectangle) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.region = region
}

func (c *Camera) SetCaptureResize(resize *image.Point) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.resize_to = resize
}

func (c *Camera) SetCaptureThreshold(threhold int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.threshold = threhold
}

func (c *Camera) SetCaptureDebug(debug bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.debug = debug
}

func (c *Camera) Reset() error {
	c.SetCaptureDebug(DefaultDebug)
	c.SetCaptureThreshold(DefaultThreshold)
	c.SetCaptureRegion(nil)
	c.SetCaptureResize(nil)
	return nil
}

func (c *Camera) String() string {
	const PRINT_FORMAT = `
| --- Source ---
URI:		%s
Dimensions:	%dx%d
| --- Parameters ---
Debug:		%v
Threshold:	%d
Region: 	%+v
Resize: 	%+v
`
	return fmt.Sprintf(PRINT_FORMAT,
		c.source.URI,
		c.source.Dimensions.Width,
		c.source.Dimensions.Height,
		c.debug,
		c.threshold,
		c.region,
		c.resize_to,
	)
}
