package handler

import (
	"belajar/database"
	"belajar/model/entity"
	"belajar/model/request"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func BookHandlerCreate(ctx *fiber.Ctx) error {
	book := new(request.BookCreateRequest)
	if err := ctx.BodyParser(book); err != nil {
		return err
	}

	//Validasi Request
	validate := validator.New()
	errValidate := validate.Struct(book)
	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   errValidate.Error(),
		})
	}

	// Validation Require Image
	var filenameString string

	filename := ctx.Locals("filename")
	log.Println("filename =", filename)
	if filename == nil {
		return ctx.Status(422).JSON(fiber.Map{
			"message": "image cover is required!",
		})
	} else {
		filenameString = fmt.Sprintf("%v", filename)
	}

	// Mengirim data
	newBook := entity.Book{
		Title:  book.Title,
		Author: book.Author,
		Cover:  filenameString,
	}

	errCreateBook := database.DB.Create(&newBook).Error
	if errCreateBook != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    newBook,
	})
}
