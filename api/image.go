package api

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"text/template"

	"github.com/nfnt/resize"
)

type imageBundle struct {
	bundle
}

func newImageBundle(b bundle) imageBundle {
	return imageBundle{b}
}

func (i *imageBundle) Process(r *http.Request, id int) string {
	switch r.Method {
	case getMethod:
		return i.Get(r, id)
	case postMethod:
		return i.Post(r, id)
	case putMethod:
		return jsonError("Method not implemented")
	case deleteMethod:
		return jsonError("Method not implemented")
	default:
		return jsonError("Unknown HTTP Method")
	}
}

func (i *imageBundle) Get(r *http.Request, id int) string {
	sqlReq := fmt.Sprintf("SELECT img FROM images WHERE msg_id = %d", id)
	rows, err := i.db.Query(sqlReq)
	if err != nil {
		return jsonError(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		// Get img
		var img string
		if err := rows.Scan(&img); err != nil {
			return jsonError(err.Error())
		}

		return fmt.Sprintf("{\"img\":%q}", img)
	}

	return jsonError("No image corresponding to this message ID")
}

func (i *imageBundle) Post(r *http.Request, id int) string {
	var err error

	// Get image from request
	image := r.FormValue("image")
	if image == "" {
		return jsonError("No image provided")
	}

	// Process img
	if image, err = processImg(i.data, image); err != nil {
		return jsonError(err.Error())
	}

	// Save image into DB
	if err := saveImg(i.db, image, id); err != nil {
		return jsonError(err.Error())
	}

	return jsonResponseOk()
}

// Private functions
////////////////////////////////////////////////////////////////////////////////

// Process image
func processImg(data *data, dataURL string) (string, error) {
	b64Img, mime := parseDataURL(dataURL)

	// Decode
	rawImg, err := base64.StdEncoding.DecodeString(b64Img)
	if err != nil {
		return "", err
	}

	// Read image by its type
	img, config, err := imgRawDecode(rawImg, mime)
	if err != nil {
		return "", err
	}

	// Resize if width is to high
	if config.Width > data.imgMaxWidth {
		img = resize.Resize(uint(data.imgMaxWidth), 0, img, resize.Lanczos3)
	}

	// Resize if height is to high
	if config.Height > data.imgMaxHeight {
		img = resize.Resize(0, uint(data.imgMaxHeight), img, resize.Lanczos3)
	}

	// Encode
	rawImg, err = imgRawEncode(img, mime)
	if err != nil {
		return "", err
	}
	b64Img = base64.StdEncoding.EncodeToString(rawImg)

	return formatDataURL(b64Img, mime), nil
}

// Save image
func saveImg(db *sql.DB, img string, id int) error {
	// TODO: Check if there is already an image

	// Prepare add to database
	stmt, err := db.Prepare("INSERT images SET time=?,msg_id=?,img=?")
	if err != nil {
		return err
	}

	// Exec request
	_, err = stmt.Exec(sqlTimeNow(), id, img)
	if err != nil {
		return err
	}

	// Prepare request
	stmt, errPrep := db.Prepare("UPDATE messages SET img=? WHERE id=?")
	if errPrep != nil {
		return errPrep
	}

	// Exec request
	_, errExec := stmt.Exec(true, id)
	if errExec != nil {
		return errExec
	}

	return nil
}

func parseDataURL(dataURL string) (string, string) {
	// Check XSS
	dataURL = template.HTMLEscapeString(dataURL)

	// Remove infos
	infos, b64Img := splitString(dataURL, ";base64,")
	_, mime := splitString(infos, ":")

	// Format mimes
	if mime == "image/jpg" {
		mime = jpegMIME
	}

	return b64Img, mime
}

func formatDataURL(b64Img, mime string) string {
	return "data:" + mime + ";base64," + b64Img
}

func imgRawDecode(rawImg []byte, mime string) (image.Image, image.Config, error) {
	var img image.Image
	var config image.Config
	var err error

	imgReader := bytes.NewReader(rawImg)
	configReader := bytes.NewReader(rawImg)
	switch mime {
	case gifMIME:
		if img, err = gif.Decode(imgReader); err != nil {
			return img, config, err
		} else if config, err = gif.DecodeConfig(configReader); err != nil {
			return img, config, err
		}
	case jpegMIME:
		if img, err = jpeg.Decode(imgReader); err != nil {
			return img, config, err
		} else if config, err = jpeg.DecodeConfig(configReader); err != nil {
			return img, config, err
		}
	case pngMIME:
		if img, err = png.Decode(imgReader); err != nil {
			return img, config, err
		} else if config, err = png.DecodeConfig(configReader); err != nil {
			return img, config, err
		}
	default:
		return img, config, errors.New("Unrecognized image format")
	}

	return img, config, nil
}

func imgRawEncode(img image.Image, mime string) ([]byte, error) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	switch mime {
	case gifMIME:
		if err := gif.Encode(writer, img, nil); err != nil {
			return nil, err
		}
	case jpegMIME:
		if err := jpeg.Encode(writer, img, nil); err != nil {
			return nil, err
		}
	case pngMIME:
		if err := png.Encode(writer, img); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Unrecognized image format")
	}

	return buffer.Bytes(), nil
}
