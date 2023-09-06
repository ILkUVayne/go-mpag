package main

import (
	"fmt"
	"go-mpeg/linuxExec/watermark"
)

// ffmpeg -hide_banner -i /mnt/f/video/1251203951-1-192.mp4 -i /mnt/f/video/watermark.jpg -filter_complex "overlay=x=0:y=0" out.mp4 -y
// ffmpeg -i /mnt/f/video/1251203951-1-192.mp4 -vf "drawtext=fontsize=100:fontfile=lazy.ttf:text='hello world':x=20:y=20:fontcolor=#123fff@.5:box=1:boxcolor=yellow" out1.mp4
// ffmpeg -i /mnt/f/video/1251203951-1-192.mp4 -vf "drawtext=fontsize=100:fontfile=lazy.ttf:textfile=./1.text:x=20:y=20:fontcolor=#123fff:box=1:boxcolor=yellow" out1.mp4

// 简单filtergraphs 视频和音频分别-vf和-af
func main() {
	teXtWM()
}

func teXtWM() {
	srcPath := "/mnt/f/video/1251203951-1-192.mp4"
	dstPath := "/mnt/f/video/out6.mp4"
	dt1 := watermark.NewDrawText(
		watermark.WithText("happy!"),
		watermark.WithTransparency(0.3),
		watermark.WithFontsize(50),
		watermark.WithPosition(10, 20),
	)
	dt2 := watermark.NewDrawText(
		watermark.WithText("hello!"),
		watermark.WithTransparency(0.6),
		watermark.WithFontsize(100),
		watermark.WithPosition(20, 100),
		watermark.WithFontcolor("red"),
	)
	dts := watermark.DrawTexts{
		dt1,
		dt2,
	}
	err := dts.Watermark(srcPath, dstPath)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
