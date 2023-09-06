package watermark

import "fmt"

func TextWM() {
	srcPath := "/mnt/f/video/1251203951-1-192.mp4"
	dstPath := "/mnt/f/video/out1.mp4"
	dt1 := NewDrawText(
		WithText("happy!"),
		WithAlpha(0.3),
		WithFontsize(50),
		WithPosition(10, 20),
	)
	dt2 := NewDrawText(
		WithText("hello!"),
		WithAlpha(0.6),
		WithFontsize(100),
		WithPosition(20, 100),
		WithFontcolor("red"),
	)
	dts := DrawTexts{
		dt1,
		dt2,
	}
	err := dts.Watermark(srcPath, dstPath)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}

func FcWM() {
	srcPath := "/mnt/f/video/1251203951-1-192.mp4"
	dstPath := "/mnt/f/video/f_out3.mp4"
	fc1 := NewFilterComplex(
		WithOverlay("10", "10"),
	)
	fc2 := NewFilterComplex(
		WithOverlay("main_w-overlay_w-10", "10"),
	)
	fc3 := NewFilterComplex(
		WithOverlay("10", "main_h-overlay_h-10"),
	)
	fc4 := NewFilterComplex(
		WithOverlay("main_w-overlay_w-10", "main_h-overlay_h-10"),
	)
	fcs := FilterComplexes{
		fc1,
		fc2,
		fc3,
		fc4,
	}
	mark := []string{
		"/mnt/f/video/watermark.jpg",
		"/mnt/f/video/watermark.jpg",
		"/mnt/f/video/watermark.jpg",
		"/mnt/f/video/watermark.jpg",
	}
	err := fcs.Watermark(srcPath, dstPath, mark...)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
