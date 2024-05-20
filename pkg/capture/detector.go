package capture

import (
	"image"

	"gocv.io/x/gocv"
)

type Detector struct {
	prev      gocv.Mat
	threshold int
}

func NewDetector(threshold int) *Detector {
	return &Detector{
		prev:      gocv.NewMat(),
		threshold: threshold,
	}

}

func (d *Detector) Motion(f gocv.Mat) bool {
	return len(d.Detect(f)) > 0
}

func (d *Detector) Detect(f gocv.Mat) []image.Rectangle {
	// Prepare frame
	frame := gocv.NewMatWithSize(f.Rows(), f.Cols(), f.Type())
	gocv.CvtColor(f, &frame, gocv.ColorBGRToGray)

	if d.prev.Empty() || d.prev.Rows() != f.Rows() || d.prev.Cols() != f.Cols() {
		d.prev.Close()
		d.prev = frame
		return []image.Rectangle{}
	}

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(9, 9))
	mask := d.get_mask(d.prev, frame, kernel)
	kernel.Close()

	detections := d.get_contour_detections(mask, d.threshold)
	mask.Close()

	d.prev.Close()
	d.prev = frame

	return detections
}

func (d *Detector) Close() error {
	d.prev.Close()
	return nil
}

func (d *Detector) get_contour_detections(mask gocv.Mat, threshold int) (detections []image.Rectangle) {
	contours := gocv.FindContours(
		mask,
		gocv.RetrievalExternal,
		gocv.ChainApproxTC89L1,
	)

	for i := 0; i < contours.Size(); i++ {
		detection := gocv.BoundingRect(contours.At(i))
		area := detection.Dx() * detection.Dy()
		if area <= threshold {
			continue
		}
		detections = append(detections, detection)
	}
	detections = unique(detections)
	return
}

func (d *Detector) get_mask(f1, f2 gocv.Mat, kernel gocv.Mat) gocv.Mat {
	diff := gocv.NewMat()
	defer diff.Close()
	gocv.Subtract(f2, f1, &diff)
	gocv.MedianBlur(diff, &diff, 3)

	mask := gocv.NewMat()
	gocv.AdaptiveThreshold(
		diff,
		&mask,
		255,
		gocv.AdaptiveThresholdGaussian,
		gocv.ThresholdBinaryInv,
		11,
		3,
	)

	gocv.MedianBlur(mask, &mask, 3)
	gocv.MorphologyEx(mask, &mask, gocv.MorphClose, kernel)

	return mask
}

func unique(rs []image.Rectangle) []image.Rectangle {
	u := make([]image.Rectangle, 0)
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		overlapped := false
		for j := 0; j < len(u); j++ {
			if r == u[j] {
				continue
			}

			if r.In(u[j]) {
				goto skip
			}

			if u[j].Overlaps(r.Inset(-100)) {
				overlapped = true
				u[j] = u[j].Union(r)
			}
		}

		if !overlapped {
			u = append(u, r)
		}
	skip:
	}

	return u
}
