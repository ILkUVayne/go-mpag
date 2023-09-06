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

/**
添加文字水印
*/

type Option func(dt *DrawText)

type DrawText struct {
	fontFile     string
	text         string
	textFile     string
	fontcolor    string
	box          bool
	boxColor     string
	fontsize     int
	x, y         int
	transparency float64
}

type DrawTexts []*DrawText

func NewDrawText(opts ...Option) *DrawText {
	dt := &DrawText{
		fontFile:     "lazy.ttf",
		fontcolor:    "#ffffff",
		box:          false,
		boxColor:     "#ffffff",
		fontsize:     100,
		transparency: 1,
	}
	for _, v := range opts {
		v(dt)
	}
	return dt
}

func WithFontFile(fontFile string) Option {
	return func(dt *DrawText) {
		dt.fontFile = fontFile
	}
}

func WithText(text string) Option {
	return func(dt *DrawText) {
		dt.text = text
	}
}

func WithTextFile(textFile string) Option {
	return func(dt *DrawText) {
		dt.textFile = textFile
	}
}

func WithBox(box bool) Option {
	return func(dt *DrawText) {
		dt.box = box
	}
}

func WithBoxColor(boxColor string) Option {
	return func(dt *DrawText) {
		dt.boxColor = boxColor
	}
}
func WithFontsize(fontsize int) Option {
	return func(dt *DrawText) {
		dt.fontsize = fontsize
	}
}
func WithTransparency(transparency float64) Option {
	return func(dt *DrawText) {
		dt.transparency = transparency
	}
}
func WithFontcolor(fontcolor string) Option {
	return func(dt *DrawText) {
		dt.fontcolor = fontcolor
	}
}
func WithPosition(x, y int) Option {
	return func(dt *DrawText) {
		dt.x = x
		dt.y = y
	}
}

func (dt *DrawTexts) Watermark(srcPath, dstPath string) error {
	// 判断dstPath是否存在
	exists, err := common.PathExists(dstPath)
	if err != nil {
		return err
	}
	if exists {
		return errors.New(fmt.Sprintf("dstPath %s is exists", dstPath))
	}
	// 构建参数
	cmd := ""
	for _, v := range *dt {
		if cmd == "" {
			cmd = cmd + v.buildArgs()
			continue
		}
		cmd = cmd + "," + v.buildArgs()
	}

	if cmd == "" {
		return errors.New("text or textFile is empty")
	}
	println(cmd)
	c := exec.Command("ffmpeg", "-i", srcPath, "-vf", cmd, dstPath)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Run()
	if err != nil {
		log.Fatalf("failed to call cmd.Run(): %v", err)
	}
	return nil
}

func (dt *DrawText) buildArgs() string {
	cmd := "drawtext=fontfile=" + dt.fontFile + ":"
	cmd += "fontsize=" + strconv.Itoa(dt.fontsize) + ":"
	cmd += "fontcolor=" + dt.fontcolor + "@" + strconv.FormatFloat(dt.transparency, 'f', 2, 64) + ":"
	cmd += "x=" + strconv.Itoa(dt.x) + ":" + "y=" + strconv.Itoa(dt.y) + ":"
	if dt.box {
		cmd += "box=1:boxcolor=" + dt.boxColor + ":"
	}
	if dt.textFile != "" {
		cmd += "textfile=" + dt.textFile
		return cmd
	}
	if dt.text != "" {
		cmd += "text=" + dt.text
		return cmd
	}

	return ""
}
