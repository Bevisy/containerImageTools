package layers

import (
	"archive/tar"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

type Image struct {
	Fpath   string
	Destdir string
}

type Manifests []Manifest

type Manifest struct {
	Config  string
	RepoTag string
	Layers  []string
}

func NewImage(fpath string, destdir string) *Image {
	i := &Image{
		Fpath:   fpath,
		Destdir: destdir,
	}
	return i
}

func (i *Image) Unzip() error {
	// 解压镜像tar包
	f, err := os.Open(i.Fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	fread := tar.NewReader(f)

	for {
		h, err := fread.Next()
		if err != nil {
			return err
		}
		if h.FileInfo().IsDir() {
			if err := os.MkdirAll(path.Join(i.Destdir, h.Name), h.FileInfo().Mode()); err != nil {
				return err
			}
		} else {
			f, err := os.OpenFile(path.Join(i.Destdir, h.Name), os.O_RDWR|os.O_CREATE, h.FileInfo().Mode())
			if err != nil {
				return err
			} else {
				if _, err := io.Copy(f, fread); err != nil {
					return err
				}
			}
			f.Close()
		}
	}
}

func (i *Image) RemoveDir(path string) error {
	// 删除目录
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

func HashSHA256(fpath string) (string, error) {
	// 计算文件sha256值
	var hashValue string
	f, err := os.Open(fpath)
	if err != nil {
		return hashValue, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return hashValue, err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (i *Image) Zip() error {
	// 打包镜像文件,替换原文件
	tarfd, _ := os.Stat(i.Fpath)
	tmptar := path.Join(os.TempDir(), path.Base(i.Fpath))
	if _, err := os.Stat(tmptar); err == nil {
		os.Remove(tmptar)
	}
	tmpfd, err := os.OpenFile(tmptar, os.O_RDWR|os.O_CREATE, tarfd.Mode())
	if err != nil {
		return err
	}
	defer tmpfd.Close()

	w := tar.NewWriter(tmpfd)
	err = filepath.Walk(i.Destdir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relpath, err := filepath.Rel(i.Destdir, p)
		if err != nil {
			return err
		}
		if relpath == "." {
			return nil
		}
		if info.IsDir() {
			h, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			} else {
				h.Name = relpath
				if err := w.WriteHeader(h); err != nil {
					return err
				}
			}
		} else {
			if f, err := os.Open(p); err != nil {
				return err
			} else {
				defer f.Close()
				if h, err := tar.FileInfoHeader(info, relpath); err != nil {
					return err
				} else {
					h.Name = relpath
					if err := w.WriteHeader(h); err != nil {
						return err
					} else {
						if _, err := io.Copy(w, f); err != nil {
							return err
						}
					}
				}
			}
		}
		return nil
	})
	defer w.Close()
	if err != nil {
		return err
	} else {
		w.Flush()
	}
	return nil
}

func Move(oldFile, newFile string) error {
	fold, err := os.Open(oldFile)
	if err != nil {
		log.Printf("open %s failed! error: %s\n", oldFile, err)
		return err
	}
	defer fold.Close()
	if _, err = os.Stat(newFile); err == nil {
		os.Remove(newFile) // if newfile exists, delete it
	}
	fnew, err := os.OpenFile(newFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("open %s failed! error: %s \n", newFile, err)
		return err
	}
	defer fnew.Close()
	_, err = io.Copy(fnew, fold)
	if err != nil {
		log.Printf("copy from %s to %s failed! error: %s\n", oldFile, newFile, err)
		return err
	}
	return nil
}

func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
