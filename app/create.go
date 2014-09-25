package app

import (
	"fmt"

	car "github.com/mohae/carchivum"
	"github.com/mohae/contour"
)

func Create(destination string, sources ...string) (string, error) {
	var err error
	var message string

	fmt.Printf("\nCreate %q from %v\n", destination, sources)

	switch contour.GetString(CfgFormat) {
	case "tar":
		message, err = createTar(destination, sources...)
	case "zip":
		message, err = createZip(destination, sources...)
	default:
		err = fmt.Errorf("%s not supported", contour.GetString(CfgFormat))
	}

	if err != nil {
		logger.Error(err)
		return "", err
	}

	return message, nil
}

func createZip(destination string, sources ...string) (string, error) {
	var err error

	logger.Debugf("Creating zip: %s from %s", destination, sources)
	zipper := car.NewZipArchive()
	zipper.Name = destination
	zipper.UseFullpath = contour.GetBool("usefullpath")
	_, err = zipper.CreateFile(destination, sources...)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	return zipper.Message(), nil
}

func createTar(destination string, sources ...string) (string, error) {
	var err error

	logger.Debugf("Creating tar: %s from %s", destination, sources)
	tballer := car.NewTar()
	tballer.Name = destination
	tballer.UseFullpath = contour.GetBool("usefullpath")
	_, err = tballer.CreateFile(destination, sources...)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	return tballer.Message(), nil
}
