package main

// 人脸检测，基于pico方法

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"gocv.io/x/gocv"
)

func main() {

	webcam, err := gocv.VideoCaptureDevice(0)

	if nil != err {
		fmt.Println("VideoCaptureDevice err ", err)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("pigo")
	defer window.Close()

	window.ResizeWindow(640, 480)

	img := gocv.NewMat()
	defer img.Close()

	//green := color.RGBA{0, 255, 0, 0}

	cascadeFile, err := os.ReadFile("D:\\data\\file\\code\\demo\\third-party-libraries\\face_recognition\\cascade\\facefinder")

	if err != nil {
		log.Fatalf("Error reading the cascade file: %v", err)
	}

	for {

		if ok := webcam.Read(&img); !ok {
			fmt.Println("Read err ")
			return
		}

		if img.Empty() {
			continue
		}

		//// 预处理：灰度转换
		//gocv.CvtColor(img, &img, gocv.ColorBGRToGray)
		//
		//// 预处理：直方图均衡化
		//gocv.EqualizeHist(img, &img)

		//去噪：去除图像中的噪点，有助于减少误报
		gocv.GaussianBlur(img, &img, image.Pt(5, 5), 0, 0, gocv.BorderDefault)

		//通过边缘检测技术（如 Canny 边缘检测），可以增强图像中的结构信息，帮助更准确地识别人脸。
		//gocv.Canny(img, &img, 100, 200)

		goImg, err := img.ToImage()

		if nil != err {
			fmt.Println("ToImage err ")
			return
		}

		pixels := pigo.RgbToGrayscale(goImg)
		cols, rows := goImg.Bounds().Max.X, goImg.Bounds().Max.Y

		cParams := pigo.CascadeParams{
			MaxSize:     1000, // 最大人脸尺寸
			MinSize:     200,  // 最小人脸尺寸
			ShiftFactor: 0.1,
			ScaleFactor: 1.1,

			ImageParams: pigo.ImageParams{
				Pixels: pixels,
				Rows:   rows,
				Cols:   cols,
				Dim:    cols,
			},
		}

		pPigo := pigo.NewPigo()

		classifier, err := pPigo.Unpack(cascadeFile)
		if err != nil {
			log.Fatalf("Error reading the cascade file: %s", err)
		}

		angle := 30.0
		iouThreshold := 0.3

		dets := classifier.RunCascade(cParams, angle)

		dets = classifier.ClusterDetections(dets, iouThreshold)

		faceDetected := false // 标记是否检测到人脸

		for _, face := range dets {

			if face.Q > 7 {
				//x := face.Col - face.Scale/2
				//y := face.Row - face.Scale/2
				//r := image.Rect(x, y, x+face.Scale, y+face.Scale)
				//gocv.Rectangle(&img, r, green, 3)
				faceDetected = true
			} else {
				continue
			}
		}

		// 如果检测到人脸，则保存图片
		if faceDetected {
			// 获取当前时间作为文件名
			fileName := fmt.Sprintf("D:\\data\\sucess\\files\\image\\image_%d.jpg", time.Now().Unix())
			if ok := gocv.IMWrite(fileName, img); ok {
				fmt.Println("保存图片: ", fileName)
			} else {
				fmt.Println("保存图片失败")
			}
		}

		window.IMShow(img)

		if 27 == window.WaitKey(1) {
			break
		}

	}

}
