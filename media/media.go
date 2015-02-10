package media

import (
    "crypto/md5"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

// file types that I found in my library..
var extpatterns = []string{
    "JPG", "jpg",
    "JPEG", "jpeg",
    "PNG", "png",
    "MOV", "mov",
    "MP4", "mp4",
    "3GP", "3gp",
    "PSD", "psd",
    "MPG", "mpg",
    "MPEG", "mpeg",
    "TIF", "tif",
    "GIF", "gif",
}

type Media struct {
    FullPath string
    Name     string
    Size     int64
    Mode     os.FileMode
    ModTime  time.Time
    Md5Sum   []byte
}

func NewMedia(path string) (media *Media, err error) {
    f, err := os.Open(path)
    if err != nil {
        fmt.Println(err)
        return &Media{}, err
    }
    fstat, err := f.Stat()
    if err != nil {
        fmt.Println(err)
    }
    if fstat.Mode().IsRegular() && !fstat.Mode().IsDir() {
        return &Media{
            FullPath: path,
            Name:     fstat.Name(),
            Size:     fstat.Size(),
            Mode:     fstat.Mode(),
            ModTime:  fstat.ModTime(),
        }, err
    }
    return &Media{}, err
}

func (m *Media) CalcMd5() {
    var md5sum []byte

    fp := m.FullPath

    data, err := os.Open(fp)
    if err != nil {
        fmt.Println("could not open:", fp)
        return
    }
    defer data.Close()

    hash := md5.New()
    if _, err := io.Copy(hash, data); err != nil {
        fmt.Println("could not copy from", data, "to", hash)
    }

    md5sum = hash.Sum(md5sum)
    m.Md5Sum = md5sum
}

// Is the media valid or bogus?
func (m *Media) IsValid() bool {
    if ext := filepath.Ext(m.FullPath); len(ext) > 0 {
        for _, pattern := range extpatterns {
            match, err := filepath.Match("."+pattern, ext)
            if err != nil {
                fmt.Println("malfoemd pattern?")
                return false
            }
            if match {
                return true
            }
        }
    }

    return false
}
