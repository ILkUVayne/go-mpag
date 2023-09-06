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

type DrawText struct {
	fontFile  string
	text      string
	textFile  string
	fontcolor string
	box       bool
	boxColor  string
	fontsize  int
	x, y      int
	alpha     float64
}

type DrawTexts []*DrawText

func NewDrawText(opts ...common.Option) *DrawText {
	dt := &DrawText{
		fontFile:  "lazy.ttf",
		fontcolor: "#ffffff",
		box:       false,
		boxColor:  "#ffffff",
		fontsize:  100,
		alpha:     1,
	}
	for _, v := range opts {
		v(dt)
	}
	return dt
}

func WithFontFile(fontFile string) common.Option {
	return func(st interface{}) {
		st.(*DrawText).fontFile = fontFile
	}
}

func WithText(text string) common.Option {
	return func(st interface{}) {
		st.(*DrawText).text = text
	}
}

func WithTextFile(textFile string) common.Option {
	return func(st interface{}) {
		st.(*DrawText).textFile = textFile
	}
}

func WithBox(box bool) common.Option {
	return func(st interface{}) {
		st.(*DrawText).box = box
	}
}

func WithBoxColor(boxColor string) common.Option {
	return func(st interface{}) {
		st.(*DrawText).boxColor = boxColor
	}
}
func WithFontsize(fontsize int) common.Option {
	return func(st interface{}) {
		st.(*DrawText).fontsize = fontsize
	}
}
func WithAlpha(transparency float64) common.Option {
	return func(st interface{}) {
		st.(*DrawText).alpha = transparency
	}
}
func WithFontcolor(fontcolor string) common.Option {
	return func(st interface{}) {
		st.(*DrawText).fontcolor = fontcolor
	}
}
func WithPosition(x, y int) common.Option {
	return func(st interface{}) {
		st.(*DrawText).x = x
		st.(*DrawText).y = y
	}
}

func (dt *DrawText) buildArgs() string {
	cmd := "drawtext=fontfile=" + dt.fontFile + ":"
	cmd += "fontsize=" + strconv.Itoa(dt.fontsize) + ":"
	cmd += "fontcolor=" + dt.fontcolor + ":"
	cmd += "alpha=" + strconv.FormatFloat(dt.alpha, 'f', 2, 64) + ":"
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

func (dt *DrawTexts) Watermark(srcPath, dstPath string, _ ...string) error {
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
	c.Args = append(c.Args, dt.buildArgs(srcPath, dstPath)...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Run()
	if err != nil {
		log.Fatalf("failed to call cmd.Run(): %v", err)
	}
	return nil
}

func (dt *DrawTexts) buildArgs(srcPath, dstPath string) []string {
	cmd := ""
	for _, v := range *dt {
		if cmd == "" {
			cmd = cmd + v.buildArgs()
			continue
		}
		cmd = cmd + "," + v.buildArgs()
	}
	if cmd == "" {
		log.Fatal("text or textFile is empty")
	}
	return []string{srcPath, "-vf", cmd, dstPath}
}
