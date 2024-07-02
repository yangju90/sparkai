package subImage

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func SaveImagePic(base64Data string, uuid string) {

	img, err := DecodeBase64Image(base64Data)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	path := "D:/goconfig/static/object/" + uuid + ".png"
	err = EncodeImage(path, img, "png")
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}
}

func ToResizeSubImage(base64Data string, uuid string, box []int) {

	img, err := DecodeBase64Image(base64Data)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// 指定截取的范围并截取图片
	croppedImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(box[0], box[1], box[2], box[3]))

	// 调整截取后的图片大小
	resizedImg := resize.Resize(100, 0, croppedImg, resize.Lanczos3)

	path := "D:/goconfig/static/object/" + uuid + ".png"

	err = EncodeImage(path, resizedImg, "png")
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}
}

// DecodeBase64Image 解码 Base64 编码的图片数据
func DecodeBase64Image(base64Data string) (image.Image, error) {
	parts := strings.Split(base64Data, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid data URI: expected 2 parts, got %d", len(parts))
	}
	base64Str := parts[1]

	// 将Base64编码解码为字节流
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	// 创建一个图片对象
	img, _, err := image.Decode(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	return img, nil
}

// EncodeImage 将图片编码为指定格式，并写入输出流
func EncodeImage(path string, img image.Image, format string) error {
	out, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return err
	}
	defer out.Close()

	switch format {
	case "jpeg":
		return jpeg.Encode(out, img, nil)
	case "png":
		return png.Encode(out, img)
	case "gif":
		return gif.Encode(out, img, nil)
	default:
		return fmt.Errorf("unsupported image format: %s", format)
	}
}

func DrawRectanglePic(base64Data string, uuid string, boxes map[string][]interface{}) error {
	img, err := DecodeBase64Image(base64Data)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return err
	}
	dc := gg.NewContextForImage(img)

	// red := color.RGBA{255, 0, 0, 255} // 红色
	var colors = []color.RGBA{
		{255, 0, 0, 255},   // 红色
		{0, 255, 0, 255},   // 绿色
		{0, 0, 255, 255},   // 蓝色
		{255, 255, 0, 255}, // 黄色
		{255, 0, 255, 255}, // 品红
		{0, 255, 255, 255}, // 青色
	}

	imageWidth := float64(img.Bounds().Dx())
	borderWidth := imageWidth * 0.002

	indexColor := 0

	for _, v := range boxes {
		color := colors[indexColor]
		indexColor++
		if indexColor == 6 {
			break
		}
		for _, v1 := range v {
			v2 := v1.([]interface{})
			x := v2[0].(float64)
			y := v2[1].(float64)
			w := v2[2].(float64) - v2[0].(float64)
			h := v2[3].(float64) - v2[1].(float64)
			dc.SetLineWidth(borderWidth)
			dc.SetColor(color)
			dc.DrawRectangle(x, y, w, h)
			dc.Stroke()
		}
	}
	// dc.DrawRectangle(float64(box[0]), float64(box[1]), float64(box[2])-float64(box[0]), float64(box[3])-float64(box[1]))

	path := "D:/goconfig/static/object/" + uuid + ".png"

	if err := dc.SavePNG(path); err != nil {
		return err
	}

	fmt.Println("Red rectangle drawn successfully.")
	return nil
}
