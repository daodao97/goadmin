package user

import (
	"bytes"
	"fmt"
	"github.com/daodao97/goadmin/pkg/util"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// 生成随机数
func randNum(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func GetRandStr(n int) (randStr string) {
	chars := "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
	charsLen := len(chars)
	if n > 10 {
		n = 10
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(charsLen)
		randStr += chars[randIndex : randIndex+1]
	}
	return randStr
}

func CaptchaImgText(width, height int, text string) (b []byte) {
	textLen := len(text)
	dc := gg.NewContext(width, height)
	bgR, bgG, bgB, bgA := getRandColorRange(240, 255)
	dc.SetRGBA255(bgR, bgG, bgB, bgA)
	dc.Clear()

	// 干扰线
	for i := 0; i < 10; i++ {
		x1, y1 := getRandPos(width, height)
		x2, y2 := getRandPos(width, height)
		r, g, b, a := getRandColor(255)
		w := float64(rand.Intn(3) + 1)
		dc.SetRGBA255(r, g, b, a)
		dc.SetLineWidth(w)
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	fontSize := float64(height/2) + 5
	face := loadFontFace(fontSize)
	dc.SetFontFace(face)

	for i := 0; i < len(text); i++ {
		r, g, b, _ := getRandColor(100)
		dc.SetRGBA255(r, g, b, 255)
		fontPosX := float64(width/textLen*i) + fontSize*0.6

		writeText(dc, text[i:i+1], float64(fontPosX), float64(height/2))
	}

	buffer := bytes.NewBuffer(nil)
	dc.EncodePNG(buffer)
	b = buffer.Bytes()
	return
}

func LoginMD5(username, passwd, captcha string) string {
	return util.Md5(fmt.Sprintf("%s%s%s", username, passwd, captcha))
}

// 渲染文字
func writeText(dc *gg.Context, text string, x, y float64) {
	xfload := 5 - rand.Float64()*10 + x
	yfload := 5 - rand.Float64()*10 + y

	radians := 40 - rand.Float64()*80
	dc.RotateAbout(gg.Radians(radians), x, y)
	dc.DrawStringAnchored(text, xfload, yfload, 0.2, 0.5)
	dc.RotateAbout(-1*gg.Radians(radians), x, y)
	dc.Stroke()
}

// 随机坐标
func getRandPos(width, height int) (x float64, y float64) {
	x = rand.Float64() * float64(width)
	y = rand.Float64() * float64(height)
	return x, y
}

// 随机颜色
func getRandColor(maxColor int) (r, g, b, a int) {
	r = int(uint8(rand.Intn(maxColor)))
	g = int(uint8(rand.Intn(maxColor)))
	b = int(uint8(rand.Intn(maxColor)))
	a = int(uint8(rand.Intn(255)))
	return r, g, b, a
}

// 随机颜色范围
func getRandColorRange(miniColor, maxColor int) (r, g, b, a int) {
	if miniColor > maxColor {
		miniColor = 0
		maxColor = 255
	}
	r = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	g = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	b = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	a = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	return r, g, b, a
}

// 加载字体
func loadFontFace(points float64) font.Face {
	// 这里是将字体TTF文件转换成了 byte 数据保存成了一个 go 文件 文件较大可以到附录下
	// 通过truetype.Parse可以将 byte 类型的数据转换成TTF字体类型
	f, err := truetype.Parse(goregular.TTF)

	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
	})
	return face
}

// 生成验证码V2
func generateCaptcha(w http.ResponseWriter, r *http.Request) {
	// 随机生成 4 位数字验证码
	captcha := strconv.Itoa(randNum(1000, 9999))

	// 创建空白图片
	imgWidth, imgHeight := 120, 40
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// 背景颜色随机
	bgColor := color.RGBA{uint8(randNum(200, 255)), uint8(randNum(200, 255)), uint8(randNum(200, 255)), 255}
	for x := 0; x < imgWidth; x++ {
		for y := 0; y < imgHeight; y++ {
			img.Set(x, y, bgColor)
		}
	}

	// 验证码字符颜色随机
	fontColor := color.RGBA{uint8(randNum(0, 150)), uint8(randNum(0, 150)), uint8(randNum(0, 150)), 255}

	// 在图片上写入验证码
	fontSize := 24.0
	fontFace := "Arial"
	pt := freetype.Pt(15, 25)
	dpi := float64(72)
	fg := image.NewUniform(fontColor)
	ft, _ := freetype.ParseFont([]byte(fontFace))
	ftContext := freetype.NewContext()
	ftContext.SetDPI(dpi)
	ftContext.SetFont(ft)
	ftContext.SetFontSize(fontSize)
	ftContext.SetClip(img.Bounds())
	ftContext.SetDst(img)
	ftContext.SetSrc(fg)

	// 循环写入验证码每个字符
	for _, c := range captcha {
		ftContext.DrawString(string(c), pt)
		pt.X += ftContext.PointToFixed(fontSize * 1.5)
	}

	// 输出图片
	w.Header().Set("Content-Type", "image/png")
	err := png.Encode(w, img)
	if err != nil {
		fmt.Println(err)
	}
}
