package app

import (
	"fmt"
	"os"
	_ "strconv"
	"strings"

	car "github.com/mohae/carchivum"
	"github.com/mohae/contour"
	jww "github.com/spf13/jwalterweatherman"
)

func Create(dst string, sources ...string) (string, error) {
	var err error
	var message string
	fmt.Printf("\nCreate %q from %v\n", dst, sources)
	switch contour.GetString(Format) {
	case "zip":
		message, err = createZip(dst, sources...)
	default:
		message, err = createTar(dst, sources...)
	}
	if err != nil {
		jww.ERROR.Print(err)
		return "", err
	}
	return message, nil
}

func createZip(dst string, sources ...string) (string, error) {
	jww.INFO.Printf("Creating zip: %s from %s", dst, sources)
	zipper := car.NewZip()
	zipper.Car.Name = dst
	zipper.UseFullpath = contour.GetBool("usefullpath")
	_, err := zipper.Create(dst, sources...)
	if err != nil {
		jww.ERROR.Print(err)
		return "", err
	}
	return zipper.Message(), nil
}

func createTar(dst string, sources ...string) (string, error) {
	jww.INFO.Printf("Creating tar: %s from %s", dst, sources)
	tballer := car.NewTar()
	t := contour.GetString("type")
	if t != "" {
		f := car.FormatFromString(t)
		if f == car.UnsupportedFmt {
			err := fmt.Errorf("Unuspported format: %s", t)
			jww.ERROR.Print(err)
			return "", err
		}
		tballer.Format = f
	}
	tballer.Name = dst
	tballer.Owner = contour.GetInt("owner")
	tballer.Group = contour.GetInt("group")
	tballer.FileMode = os.FileMode(contour.GetInt64("mode"))
	//	tabller.Exclude = contour.GetString("exclude")
	tballer.ExcludeAnchored = contour.GetString("exclude-anchored")
	temp := contour.GetString("exclude-ext")
	if temp != "" {
		tballer.ExcludeExt = strings.Split(temp, ",")
		tballer.ExcludeExtCount = len(tballer.ExcludeExt)
	}
	tballer.IncludeAnchored = contour.GetString("include-anchored")
	temp = contour.GetString("include-ext")
	if temp != "" {
		tballer.IncludeExt = strings.Split(temp, ",")
		tballer.IncludeExtCount = len(tballer.IncludeExt)
	}
	// TODO figure out how to convert the incoming time info to time.Time
	// tballer.NewerMTime = contour.GetTime("newer-mtime")
	//	tballer.UseFullpath = contour.GetBool("usefullpath")
	_, err := tballer.Create(sources...)
	if err != nil {
		jww.ERROR.Print(err)
		return "", err
	}
	return tballer.Message(), nil
}
