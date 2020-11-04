package camera

import (
	"fmt"
	"github.com/TwinProduction/rpi-rasptank-pro/controller"
	pigo "github.com/esimov/pigo/core"
	"github.com/pkg/errors"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"os/exec"
)

const (
	ClassifierModelFile = "data/facefinder"
)

func TakePicture() (image.Image, error) {
	err := exec.Command("/bin/bash", "-c", "/usr/bin/raspistill -o picture.jpeg -t 1").Run()
	if err != nil {
		return nil, err
	}
	file, err := os.Open("picture.jpeg")
	if err != nil {
		return nil, err
	}
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode jpeg")
	}
	return img, nil
}

func Run(arm *controller.Arm) error {
	img, err := TakePicture()
	if err != nil {
		return err
	}
	nrgbaImage := pigo.ImgToNRGBA(img)
	pixels := pigo.RgbToGrayscale(nrgbaImage)
	cParams := pigo.CascadeParams{
		MinSize:     100,
		MaxSize:     600,
		ShiftFactor: 0.15,
		ScaleFactor: 1.1,
		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   nrgbaImage.Bounds().Max.Y,
			Cols:   nrgbaImage.Bounds().Max.X,
			Dim:    nrgbaImage.Bounds().Max.X,
		},
	}
	cascade, err := ioutil.ReadFile(ClassifierModelFile)
	if err != nil {
		return errors.Wrap(err, "error unpacking the cascade file")
	}
	p := pigo.NewPigo()
	classifier, err := p.Unpack(cascade)
	if err != nil {
		return errors.Wrap(err, "error unpacking the cascade file")
	}

	detections := classifier.RunCascade(cParams, 0.0)
	//for {
	//	time.Sleep(200 * time.Millisecond)
	faces := classifier.ClusterDetections(detections, 0)
	fmt.Printf("detected %d faces\n", len(faces))
	//	if len(faces) > 0 {
	//		// move camera until the face is in the middle of the picture
	//		targetX := (img.Bounds().Max.X / 2) - (faces[0].Size().X / 2)
	//		targetY := (img.Bounds().Max.Y / 2) - (faces[0].Size().Y / 2)
	//		fmt.Printf("targetX=%d; targetY=%d\n", targetX, targetY)
	//		//arm.LookAt()
	//	}
	//}
	return nil
}
