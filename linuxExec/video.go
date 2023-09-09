package linuxExec

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

const Tool = "ffmpeg"

type metaData struct {
	majorBrand       string
	minorVersion     string
	compatibleBrands string
	encoder          string
	description      string
	duration         int
	// kb/s
	bitrate int
	err     error
}

type video struct {
	coding          string
	resolutionRatio string
	fps             float64
}

type audio struct {
	hz float64
}

type Video struct {
	metaData
	video
	audio
	File string
}

func NewVideo(file string) *Video {
	v := new(Video)
	v.File = file
	v.set()
	if v.err != nil {
		log.Fatal(v.err.Error())
	}
	return v
}

func (v *Video) GetPath() string {
	return v.File
}

// -------------------------------------- set --------------------------------------

func (v *Video) set() {
	v.setMajorBrand()
	v.setMinorVersion()
	v.setCompatibleBrands()
	v.setEncoder()
	v.setDescription()
	v.setDuration()
	v.setBitrate()
	// video
	v.setCoding()
	v.setResolutionRatio()
	v.setFps()
	// audio
	v.setHz()
}

func (v *Video) setMajorBrand() {
	v.majorBrand = v.handle(` | grep 'major_brand' | awk -F ":" '{print $2}' | awk -F " " '{print $1}'`)
}

func (v *Video) setMinorVersion() {
	v.minorVersion = v.handle(` | grep 'minor_version' | awk -F ":" '{print $2}' | awk -F " " '{print $1}'`)
}
func (v *Video) setCompatibleBrands() {
	v.compatibleBrands = v.handle(` | grep 'compatible_brands' | awk -F ":" '{print $2}' | awk -F " " '{print $1}'`)
}

func (v *Video) setEncoder() {
	v.encoder = v.handle(` | grep 'encoder' | awk -F ":" '{print $2}' | awk -F " " '{print $1}'`)
}

func (v *Video) setDescription() {
	v.description = v.handle(` | grep 'description' | awk -F ":" '{print $2}' | sed s/' '//`)
}

func (v *Video) setDuration() {
	res := v.handle(` | grep -m 1 'Duration:' | cut -d ' ' -f 4 | sed s/,// | awk -F ":" '{print(($1 * 3600) + ($2 * 60 + $3)) * 1000}'`)
	if len(res) == 0 {
		return
	}

	duration, err := strconv.Atoi(res)
	if err != nil {
		v.err = err
		return
	}
	v.duration = duration
}

func (v *Video) setBitrate() {
	res := v.handle(` | grep -m 1 'bitrate:' | cut -d ' ' -f 8`)
	if len(res) == 0 {
		return
	}

	bitrate, err := strconv.Atoi(res)
	if err != nil {
		v.err = err
		return
	}
	v.bitrate = bitrate
}

func (v *Video) setCoding() {
	v.coding = v.handle(` | grep 'Video:' | awk -F "," '{print $1}' | awk -F ":" '{print $4}' | awk -F " " '{print $1}'`)
}

func (v *Video) setResolutionRatio() {
	v.resolutionRatio = v.handle(` | grep 'Video:' | grep -m 1 'SAR' | awk -F "," '{print $4}' | awk -F " " '{print $1}'`)
}

func (v *Video) setFps() {
	res := v.handle(` | grep 'Video:' | grep -m 1 'fps' | awk -F "," '{print $6}' | awk -F " " '{print $1}'`)
	if len(res) == 0 {
		return
	}

	fps, err := strconv.ParseFloat(res, 2)
	if err != nil {
		v.err = err
		return
	}
	v.fps = fps
}

func (v *Video) setHz() {
	res := v.handle(` | grep -m 1 'Audio:'| cut -d ' ' -f 13`)
	if len(res) == 0 {
		return
	}

	hz, err := strconv.ParseFloat(res, 2)
	if err != nil {
		v.err = err
		return
	}
	v.hz = hz
}

func (v *Video) handle(args string) string {
	cmd := exec.Command("sh", "-c")
	arg := fmt.Sprintf(`%s -i %s 2>&1`, Tool, v.File) + args
	cmd.Args = append(cmd.Args, arg)
	res, err := cmd.CombinedOutput()
	if err != nil {
		v.err = err
		return ""
	}
	return strings.Trim(string(res), "\n")
}

// -------------------------------------- get --------------------------------------

func (v *Video) GetMajorBrand() string {
	return v.majorBrand
}

func (v *Video) GetMinorVersion() string {
	return v.minorVersion
}

func (v *Video) GetCompatibleBrands() string {
	return v.compatibleBrands
}

func (v *Video) GetEncoder() string {
	return v.encoder
}

func (v *Video) GetDuration() int {
	return v.duration
}

func (v *Video) GetBitrate() int {
	return v.bitrate
}

func (v *Video) GetCoding() string {
	return v.coding
}

func (v *Video) GetResolutionRatio() string {
	return v.resolutionRatio
}

func (v *Video) GetFps() float64 {
	return v.fps
}

func (v *Video) GetHz() float64 {
	return v.hz
}
