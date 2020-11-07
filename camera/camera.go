package camera

import (
	"fmt"
	"github.com/TwinProduction/rpi-rasptank-pro/controller"
	pigo "github.com/esimov/pigo/core"
	"github.com/pkg/errors"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"math"
	"os/exec"
	"time"
)

const (
	ClassifierModelFile = "data/facefinder"
)

func TakePicture() (*image.NRGBA, error) {
	err := exec.Command("/bin/bash", "-c", "/usr/bin/raspistill --quality 80 --timeout 350 -w 1600 -h 1200 --nopreview --output picture.jpeg").Run()
	if err != nil {
		return nil, err
	}
	//file, err := os.Open("picture.jpeg")
	//if err != nil {
	//	return nil, err
	//}
	//img, err := jpeg.Decode(file)
	//if err != nil {
	//	return nil, errors.Wrap(err, "failed to decode jpeg")
	//}
	img, err := pigo.GetImage("picture.jpeg")
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode jpeg")
	}
	return img, nil
}

func Run(arm *controller.Arm) error {
	cascade, err := ioutil.ReadFile(ClassifierModelFile)
	if err != nil {
		return errors.Wrap(err, "error unpacking the cascade file")
	}
	classifier, err := pigo.NewPigo().Unpack(cascade)
	if err != nil {
		return errors.Wrap(err, "error unpacking the cascade file")
	}
	targetX := arm.BaseHorizontalServo.Default
	targetY := arm.CameraVerticalServo.Default
	for {
		//time.Sleep(200 * time.Millisecond)
		faces, img, err := detectFaces(classifier)
		if err != nil {
			return err
		}
		fmt.Printf("detected %d faces\n", len(faces))
		if len(faces) > 0 {
			fmt.Printf("x=%d; xCenter=%d; y=%d; yCenter=%d; faces[0].Scale=%d\n", faces[0].Col, img.Bounds().Max.X/2, faces[0].Row, img.Bounds().Max.Y/2, faces[0].Scale)

			if math.Abs(float64((img.Bounds().Max.X/2)-faces[0].Col)) > 200 {
				if faces[0].Col > img.Bounds().Max.X/2 {
					fmt.Println("<--------")
					targetX -= 10
				} else {
					fmt.Println("-------->")
					targetX += 10
				}
			} else {
				fmt.Println("not moving bc close enough")
			}
			if math.Abs(float64((img.Bounds().Max.Y/2)-faces[0].Row)) > 100 {
				if faces[0].Row > img.Bounds().Max.Y/2 {
					fmt.Println("^")
					targetY -= 10
				} else {
					fmt.Println("v")
					targetY += 10
				}
			} else {
				fmt.Println("not moving bc close enough")
			}
		} else {
			targetX = arm.BaseHorizontalServo.Default
			targetY = arm.CameraVerticalServo.Default
		}
		arm.LookAt(targetX, targetY)
	}
	return nil
}

func detectFaces(classifier *pigo.Pigo) ([]pigo.Detection, image.Image, error) {
	start := time.Now()
	img, err := TakePicture()
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("picture taken at %s\n", time.Since(start))
	pixels := pigo.RgbToGrayscale(img)
	fmt.Printf("picture converted to grayscale at %s\n", time.Since(start))
	cParams := pigo.CascadeParams{
		MinSize:     100,
		MaxSize:     800,
		ShiftFactor: 0.15,
		ScaleFactor: 1.1,
		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   img.Bounds().Max.Y,
			Cols:   img.Bounds().Max.X,
			Dim:    img.Bounds().Max.X,
		},
	}
	detections := classifier.RunCascade(cParams, 0.0)
	fmt.Printf("cascade ran at %s\n", time.Since(start))
	faces := classifier.ClusterDetections(detections, 0)
	fmt.Printf("cluster detection finished at %s\n", time.Since(start))
	return faces, img, nil
}
