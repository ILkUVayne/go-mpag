package main

import "go-mpeg/linuxExec/watermark"

// ffmpeg -hide_banner -i /mnt/f/video/1251203951-1-192.mp4 -i /mnt/f/video/watermark.jpg -filter_complex "overlay=x=0:y=0" out.mp4 -y
// ffmpeg -i /mnt/f/video/1251203951-1-192.mp4 -vf "drawtext=fontsize=100:fontfile=lazy.ttf:text='hello world':x=20:y=20:fontcolor=#123fff@.5:box=1:boxcolor=yellow" out1.mp4
// ffmpeg -i /mnt/f/video/1251203951-1-192.mp4 -vf "drawtext=fontsize=100:fontfile=lazy.ttf:textfile=./1.text:x=20:y=20:fontcolor=#123fff:box=1:boxcolor=yellow" out1.mp4

// ffmpeg -i /mnt/f/video/1251203951-1-192.mp4 -i /mnt/f/video/watermark.jpg -i /mnt/f/video/watermark.jpg -i /mnt/f/video/watermark.jpg -i /mnt/f/video/watermark.jpg -filter_complex "[1:v]scale=160:90,format=yuva444p,colorchannelmixer=aa=0.4[img1];[2:v]scale=160:90,format=yuva444p,colorchannelmixer=aa=0.4[img2];[3:v]scale=160:90,format=yuva444p,colorchannelmixer=aa=0.4[img3];[4:v]scale=160:90,format=yuva444p,colorchannelmixer=aa=0.4[img4];[0:v][img1]overlay=x=5:y=5[01];[01][img2]overlay=x=400:y=5[012];[012][img3]overlay=x=5:y=200[0123];[0123][img4]overlay=x=400:y=200" /mnt/f/video/out.mp4;
// 简单filtergraphs 视频和音频分别-vf和-af
func main() {
	// 文字水印
	//watermark.TextWM()
	// 图片水印
	watermark.FcWM()
}
