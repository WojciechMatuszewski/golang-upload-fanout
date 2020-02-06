package formdata

import (
	"encoding/base64"
	"fmt"
	"mime"
	"mime/multipart"
	"strings"
)

type Reader struct{}

//func createFormData(p string) (*io.PipeReader, string) {
//	pr, pw := io.Pipe()
//	mw := multipart.NewWriter(pw)
//
//	go func() {
//		defer pw.Close()
//		defer mw.Close()
//		fmt.Println("doiung stuff 1")
//		pwd, err := os.Getwd()
//		if err != nil {
//			fmt.Println("error", err)
//			pw.CloseWithError(err)
//		}
//		fmt.Println("doiung stuff 2")
//
//		f, err := os.Open(pwd + p)
//		if err != nil {
//			fmt.Println("error", err)
//			pw.CloseWithError(err)
//		}
//
//		fmt.Println(pwd + p)
//
//		defer f.Close()
//		fmt.Println("doiung stuff 3")
//		fw, err := mw.CreateFormFile("image", f.Name())
//		if err != nil {
//			fmt.Println("error", err)
//			pw.CloseWithError(err)
//		}
//		fmt.Println("doiung stuff")
//		_, err = io.Copy(fw, f)
//		if err != nil {
//			fmt.Println("error", err)
//			pw.CloseWithError(err)
//		}
//
//		fmt.Println("doiung stuff")
//
//		//buf := new(bytes.Buffer)
//		//_, err = buf.ReadFrom(pr)
//		//if err != nil {
//		//	pw.CloseWithError(err)
//		//}
//		//
//		//spew.Dump(buf)
//	}()
//
//	//spew.Dump(bytes)
//	//fmt.Println("im here")
//	//
//	//var bytes []byte
//	//fmt.Println("before write")
//	//
//	//n, err := pw.Write(bytes)
//	//
//	//if err != nil {
//	//	panic(err)
//	//	//fmt.Println("could not write the bytes")
//	//	//return nil, ""
//	//}
//	//fmt.Println(n)
//	//defer pw.Close()
//
//	//err = ioutil.WriteFile("/raw", bytes, 0644)
//	//if err != nil {
//	//	fmt.Println("could not do stuff with raw")
//	//	return nil, ""
//	//}
//
//	buf := new(bytes.Buffer)
//	_, err := buf.ReadFrom(pr)
//	if err != nil {
//		panic(err)
//	}
//
//	//fmt.Println(buf.String())
//
//	ioutil.WriteFile("/Users/wn.matuszewski/Desktop/golang/testing-stuff/pkg/formdata/raw", []byte(buf.String()), 0640)
//
//	return pr, mw.Boundary()
//}

func (fdr Reader) readRaw(d, ct string) (*multipart.Form, error) {
	_, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(d)
	mr := multipart.NewReader(sr, params["boundary"])

	f, err := mr.ReadForm(32 << 20)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (fdr Reader) ReadBase64Encoded(d string, ct string) (*multipart.Form, error) {
	fmt.Println("trying too decode", d)
	p, err := base64.StdEncoding.DecodeString(d)
	if err != nil {
		return nil, err
	}

	return fdr.readRaw(string(p), ct)
}

func NewReader() *Reader {
	return &Reader{}
}
