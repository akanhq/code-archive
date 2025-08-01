package main

import (
	"gocv.io/x/gocv"
)

func main() {
	// 打开默认摄像头（通常是 0）
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		panic(err)
	}
	defer webcam.Close()

	// 创建一个窗口来显示视频流
	window := gocv.NewWindow("Camera")
	defer window.Close()

	// 创建一个 Mat 对象来存储每一帧的图像
	img := gocv.NewMat()
	defer img.Close()

	for {
		// 读取一帧图像
		if ok := webcam.Read(&img); !ok || img.Empty() {
			continue
		}

		// 在窗口中显示图像
		window.IMShow(img)

		// 等待 30 毫秒，或者直到用户按下 ESC 键
		if window.WaitKey(30) >= 0 {
			break
		}
	}
}
