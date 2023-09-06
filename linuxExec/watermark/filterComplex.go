package watermark

import (
	"errors"
	"fmt"
	"go-mpeg/common"
	"log"
	"os"
	"os/exec"
	"strconv"
)

/*
左上角 10:10
右上角 main_w-overlay_w-10:10
左下角 10:main_h-overlay_h-10
右下角 main_w-overlay_w-10 : main_h-overlay_h-10
*/

type Scale struct {
	Width, Height int
}

type Overlay struct {
	x, y string
}

type FilterComplex struct {
	scale             Scale
	overlay           Overlay
	format            string
	colorChannelMixer float64
}

type FilterComplexes []*FilterComplex

func NewFilterComplex(opts ...common.Option) *FilterComplex {
	fc := &FilterComplex{
		scale: Scale{
			Width:  160,
			Height: 90,
		},
		overlay: Overlay{
			x: "10",
			y: "10",
		},
		format:            "yuva444p",
		colorChannelMixer: 0.4,
	}
	for _, v := range opts {
		v(fc)
	}
	return fc
}

func WithScale(width, height int) common.Option {
	return func(st interface{}) {
		st.(*FilterComplex).scale.Width = width
		st.(*FilterComplex).scale.Height = height
	}
}

func WithOverlay(x, y string) common.Option {
	return func(st interface{}) {
		st.(*FilterComplex).overlay.x = x
		st.(*FilterComplex).overlay.y = y
	}
}

func WithFormat(format string) common.Option {
	return func(st interface{}) {
		st.(*FilterComplex).format = format
	}
}

func WithColorChannelMixer(ccm float64) common.Option {
	return func(st interface{}) {
		st.(*FilterComplex).colorChannelMixer = ccm
	}
}

func (fc *FilterComplexes) Watermark(srcPath, dstPath string, marker ...string) error {
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
	c.Args = append(c.Args, fc.buildArgs(srcPath, dstPath, marker...)...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Run()
	if err != nil {
		log.Fatalf("failed to call cmd.Run(): %v", err)
	}
	return nil
}

func (fc *FilterComplexes) buildArgs(srcPath, dstPath string, marker ...string) []string {
	// 构建filter_complex 内容
	cmd1, cmd2 := "", ""
	// [1:v]scale=160:90,format=yuva444p,colorchannelmixer=aa=0.4[img1]
	// [0:v][img1]overlay=x=5:y=5[01];[01][img2]overlay=x=400:y=5[012];[012][img3]overlay=x=5:y=200[0123];[0123][img4]overlay=x=400:y=200
	splitNum := 0
	for idx, v := range *fc {
		cmd11, cmd22, idxNum := v.buildArgs(idx)
		splitNum = idxNum
		if cmd1 == "" && cmd2 == "" {
			cmd1, cmd2 = cmd11, cmd22
			continue
		}
		cmd1 = cmd1 + ";" + cmd11
		cmd2 = cmd2 + ";" + cmd22
	}
	cmd2 = cmd2[:len(cmd2)-(splitNum+2)]
	cmd := cmd1 + ";" + cmd2
	println(cmd)
	// 构建命令参数
	args := []string{srcPath}
	for _, v := range marker {
		args = append(args, "-i", v)
	}
	args = append(args, "-filter_complex", cmd, dstPath)
	return args
}

func (fc *FilterComplex) buildArgs(idx int) (string, string, int) {
	nextIdx := strconv.Itoa(idx + 1)
	cmd1 := "[" + nextIdx + ":v]"
	cmd1 += "scale=" + strconv.Itoa(fc.scale.Width) + ":" + strconv.Itoa(fc.scale.Height) + ","
	cmd1 += "format=" + fc.format + ","
	cmd1 += "colorchannelmixer=aa=" + strconv.FormatFloat(fc.colorChannelMixer, 'f', 2, 64)
	cmd1 += "[img" + nextIdx + "]"

	sIdx := getIdx(idx)
	cmd2 := "[" + sIdx + "][img" + nextIdx + "]"
	if idx == 0 {
		cmd2 = "[0:v][img" + nextIdx + "]"
	}
	cmd2 += "overlay=x=" + fc.overlay.x + ":y=" + fc.overlay.y + "[" + sIdx + nextIdx + "]"
	return cmd1, cmd2, len(sIdx + nextIdx)
}

func getIdx(idx int) string {
	if idx == 0 {
		return "0"
	}
	sIdx := ""
	for i := 0; i <= idx; i++ {
		sIdx += strconv.Itoa(i)
	}
	return sIdx
}
