package handler

import (
	"belajar/database"
	"belajar/model/entity"
	"belajar/model/request"
	"belajar/utils"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func PhotoHandlerCreate(ctx *fiber.Ctx) error {
	photo := new(request.PhotoCreateRequest)
	if err := ctx.BodyParser(photo); err != nil {
		return err
	}

	//Validasi Request
	validate := validator.New()
	errValidate := validate.Struct(photo)
	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   errValidate.Error(),
		})
	}

	// Validation Require Image
	filenames := ctx.Locals("filenames")
	if filenames == nil {
		return ctx.Status(422).JSON(fiber.Map{
			"message": "image cover is required!",
		})
	} else {
		filenamesData := filenames.([]string)
		for _, filename := range filenamesData {
			// Mengirim data
			newPhoto := entity.Photo{
				Image:      filename,
				CategoryID: photo.CategoryID,
			}

			errCreatePhoto := database.DB.Create(&newPhoto).Error
			if errCreatePhoto != nil {
				log.Println("Ada file yang gagal")
			}
		}
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}

func PhotoHandlerDelete(ctx *fiber.Ctx) error {
	photoId := ctx.Params("id")

	var photo entity.Photo

	// CHECK AVAILABLE PHOTO
	err := database.DB.Debug().First(&photo, "id=?", photoId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "photo not found",
		})
	}

	// HANDLE REMOVE PHOTO
	errDeleteFile := utils.HandlerRemoveFile(photo.Image)
	if errDeleteFile != nil {
		log.Println("Fail to delete some file")
	}

	errDelete := database.DB.Debug().Delete(&photo).Error
	if errDelete != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "photo was deleted",
	})
}
