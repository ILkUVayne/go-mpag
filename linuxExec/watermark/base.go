package watermark

type WaterMarker interface {
	Watermark() error
}
