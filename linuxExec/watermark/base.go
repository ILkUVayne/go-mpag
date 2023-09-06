package watermark

type WaterMarker interface {
	Watermark(srcPath, dstPath string, marker ...string) error
}
