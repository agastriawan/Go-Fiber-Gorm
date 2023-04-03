package utils

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

const DefaultPathAssetImage = "./public/covers/"

func HandlerSingleFile(ctx *fiber.Ctx) error {
	// HANDLER FILE
	file, errFile := ctx.FormFile("cover")
	if errFile != nil {
		log.Println("Error File =", errFile)
	}

	var filename *string

	if file != nil {
		errCheckContentType := checkContentType(file, "image/jpg", "image/png", "image/jpeg")
		if errCheckContentType != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": errCheckContentType.Error(),
			})
		}

		// CUSTOM FILE NAME
		filename = &file.Filename
		extenstionFile := filepath.Ext(*filename)
		newFilename := fmt.Sprintf("gambar-satu%s", extenstionFile)

		errSaveFile := ctx.SaveFile(file, fmt.Sprintf("./public/covers/%s", newFilename)) // * menandakan tipe datanya string
		if errSaveFile != nil {
			log.Println("Fail to store file into public/covers directory")
		}

	} else {
		log.Println("Nothing file to uplodeding")
	}

	if filename != nil {
		ctx.Locals("filename", *filename)
	} else {
		ctx.Locals("filename", nil)
	}

	return ctx.Next()
}

func HandlerMultipleFile(ctx *fiber.Ctx) error {
	form, errForm := ctx.MultipartForm()
	if errForm != nil {
		log.Println("Error saat membaca Multipart Form Request, Error = ", errForm)
	}

	files := form.File["photos"]

	var filenames []string

	for i, file := range files {
		var filename string

		if file != nil {

			errCheckContentType := checkContentType(file, "image/jpg", "image/png", "image/jpeg")
			if errCheckContentType != nil {
				return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
					"message": errCheckContentType.Error(),
				})
			}

			extenstionFile := filepath.Ext(file.Filename)
			filename = fmt.Sprintf("%d-%s%s", i, "gambar", extenstionFile)

			errSaveFile := ctx.SaveFile(file, fmt.Sprintf("./public/covers/%s", filename))
			if errSaveFile != nil {
				log.Println("Gagal menyimpan file ke dalam direktori public/covers")
			}
		} else {
			log.Println("Tidak ada file yang di-upload")
		}

		if filename != "" {
			filenames = append(filenames, filename)
		}
	}

	ctx.Locals("filenames", filenames)

	return ctx.Next()
}

// HandlerRemoveFile digunakan untuk menghapus file pada direktori yang ditentukan
func HandlerRemoveFile(filename string, pathFile ...string) error {

	if len(pathFile) > 0 {
		err := os.Remove(pathFile[0] + filename)
		if err != nil {
			log.Println("Gagal menghapus file")
			return err
		}
	} else {
		err := os.Remove(DefaultPathAssetImage + filename)
		if err != nil {
			log.Println("Gagal menghapus file")
			return err
		}
	}

	return nil
}

func checkContentType(file *multipart.FileHeader, contentTypes ...string) error {
	if len(contentTypes) > 0 {
		for _, contentType := range contentTypes {
			contentTypeFile := file.Header.Get("Content-Type")
			if contentTypeFile == contentType {
				return nil
			}
		}

		return errors.New("not allowed file type")

	} else {
		return errors.New("not found content type to be checking")
	}
}
