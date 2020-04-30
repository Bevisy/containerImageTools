package main

import (
	"encoding/json"
	"github.com/bevisy/imageTool/v1/layers"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	log.SetOutput(os.Stdout)
	funds := map[string]struct{}{
		"bcf2f368fe234217249e00ad9d762d8f1a3156d60c442ed92079fa5b120634a1": {},
		"344fb4b275b72fa2f835af4a315fa3c10e6b14086126adc01471eaa57659f7a5": {},
	}
	currentdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println("get current path failed!")
		os.Exit(1)
	}
	fpath := path.Join(currentdir, os.Args[1])
	if _, err := os.Stat(fpath); err != nil {
		if os.IsNotExist(err) {
			log.Println(os.Args[1], " is NOT existed!")
			os.Exit(1)
		} else {
			log.Println("Unknown err: ", err)
			os.Exit(1)
		}
	} else {
		log.Println(os.Args[1], " existed!")
	}

	destdir := path.Join(currentdir, "temp")
	if layers.IsPathExist(destdir) {
		os.RemoveAll(destdir)
		os.Mkdir(destdir, 0744)
	} else {
		os.Mkdir(destdir, 0744)
	}

	images := layers.NewImage(fpath, destdir)
	err = images.Unzip()
	//if err != nil {
	//    log.Fatal("tar extract failed\t", err)
	//    os.Exit(1)
	//}
	log.Println("Untar success")

	var manifest layers.Manifests
	manifestPath := path.Join(destdir, "manifest.json")
	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Println("decode manifest.json success")

	// compute layer's hash256
	for _, layerpath := range manifest[0].Layers {
		layerHash256, err := layers.HashSHA256(path.Join(destdir, layerpath))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		if _, ok := funds[layerHash256]; ok {
			images.RemoveDir(path.Join(destdir, filepath.Dir(layerpath)))
		}
	}

	err = images.Zip()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Println("Tar success")

	tmppath := path.Join(os.TempDir(), os.Args[1])
	err = layers.Move(tmppath, fpath)
	if err != nil {
		log.Printf("move %s to %s failed!\n", tmppath, fpath)
		os.Exit(1)
	}
	log.Printf("move %s to %s success!\n", tmppath, fpath)

	os.Remove(tmppath)
	log.Printf("remove %s success!\n", tmppath)
	os.RemoveAll(destdir)
	log.Printf("remove %s success!\n", destdir)
}
