package watermark

type WaterMarker interface {
	Watermark(srcPath, dstPath string) error
}
