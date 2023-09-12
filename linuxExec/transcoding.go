package linuxExec

import (
	"errors"
	"fmt"
	"go-mpeg/common"
	"os"
	"os/exec"
	"strconv"
)

// ffmpeg -codecs  命令查看可用编码器
const h264 = "libx264"
const h265 = "libx265"

const aac = "aac"

// cv -vcodec -c:v
type cv struct {
	// 编码器
	vCodec string
	// 0 - 51 0无损 数字越大，质量越差，文件体积越小
	crf uint8
	// 帧率
	fps float64
	// 分辨率 ：1280x720
	resolutionRatio string
	// 去除视频内容 default false
	vn bool
}

// ca -acodec -c:a 默认直接拷贝原音频，不做编码: -c:a copy
type ca struct {
	aCodec string
	// 音频的采样率
	ar float64
	// 视频转码的时候，是否将音频给去除，default false
	an bool
	// 用来指定相对与原来的文件的音量大小,default 256
	vol int
	// copy
	copy bool
}

type Transcoding struct {
	cv
	ca
	v *Video
}

func NewTranscoding(v *Video, opts ...common.Option) *Transcoding {
	t := new(Transcoding)
	t.v = v
	// cv
	t.vCodec = h264
	t.crf = 18
	t.fps = v.video.fps
	t.resolutionRatio = v.video.resolutionRatio
	t.vn = false
	// ca
	t.aCodec = aac
	t.ar = v.audio.hz
	t.an = false
	t.vol = 256
	t.copy = true
	for _, v := range opts {
		v(t)
	}
	return t
}

func TVCodec(vCodec string) common.Option {
	return func(st interface{}) {
		st.(*Transcoding).vCodec = vCodec
	}
}

func TCrf(crf uint8) common.Option {
	return func(st interface{}) {
		st.(*Transcoding).crf = crf
	}
}

func TFps(fps float64) common.Option {
	return func(st interface{}) {
		st.(*Transcoding).fps = fps
	}
}

func TResolutionRatio(resolutionRatio string) common.Option {
	return func(st interface{}) {
		st.(*Transcoding).resolutionRatio = resolutionRatio
	}
}

func TVn(vn bool) common.Option {
	return func(st interface{}) {
		st.(*Transcoding).vn = vn
	}
}

func TACodec(aCodec string) common.Option {
	return func(st interface{}) {
		st.(*Transcoding).aCodec = aCodec
		st.(*Transcoding).copy = false
	}
}

func TAr(ar float64) common.Option {
	return func(st interface{}) {
		st.(*Transcoding).ar = ar
		st.(*Transcoding).copy = false
	}
}

func TAn(an bool) common.Option {
	return func(st interface{}) {
		st.(*Transcoding).an = an
		st.(*Transcoding).copy = false
	}
}

func TVol(vol int) common.Option {
	return func(st interface{}) {
		st.(*Transcoding).vol = vol
		st.(*Transcoding).copy = false
	}
}

func (t *Transcoding) Trans(dstPath string) error {
	// 判断dstPath是否存在
	exists, err := common.PathExists(dstPath)
	if err != nil {
		return err
	}
	if exists {
		return errors.New(fmt.Sprintf("dstPath %s is exists", dstPath))
	}
	c := exec.Command("ffmpeg", "-i")
	// 构建参数
	c.Args = append(c.Args, t.buildArgs(dstPath)...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func (t *Transcoding) buildArgs(dstPath string) []string {
	// -c:v
	args := []string{t.v.GetPath()}
	args = append(args, "-c:v", t.vCodec)
	args = append(args, "-r", strconv.FormatFloat(t.fps, 'f', 0, 64))
	args = append(args, "-s", t.resolutionRatio)
	args = append(args, "-crf", strconv.Itoa(int(t.crf)))
	if t.vn {
		args = append(args, "-vn")
	}
	// -c:a
	if t.copy {
		return append(args, "-c:a", "copy", dstPath)
	}
	args = append(args, "-c:a", t.aCodec)
	args = append(args, "-ar", strconv.FormatFloat(t.ar, 'f', 0, 64))
	args = append(args, "-vol", strconv.Itoa(t.vol))
	if t.an {
		args = append(args, "-an")
	}
	return append(args, dstPath)
}
