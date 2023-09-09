package linuxExec

type Path interface {
	GetPath() string
}

type WaterMarker interface {
	Watermark(path Path, dstPath string, marker ...string) error
}
