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

// GET ALL
func UserHandlerGetAll(ctx *fiber.Ctx) error {
	// Middleware
	userInfo := ctx.Locals("userInfo")
	log.Println("user info data ::", userInfo)

	var users []entity.User
	result := database.DB.Find(&users)
	if result.Error != nil {
		log.Println(result.Error)
	}

	return ctx.JSON(users)
}

// CREATE
func UserHandlerCreate(ctx *fiber.Ctx) error {
	user := new(request.UserCreateRequest)
	if err := ctx.BodyParser(user); err != nil {
		return err
	}

	//Validasi Request
	validate := validator.New()
	errValidate := validate.Struct(user)
	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   errValidate.Error(),
		})
	}

	newUser := entity.User{
		Name:    user.Name,
		Email:   user.Email,
		Address: user.Address,
		Phone:   user.Phone,
	}

	// hashedPassword, err := utils.HashingPassword(user.Password)
	hashedPassword, err := utils.HashingPassword(user.Password)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
		})
	}
	newUser.Password = hashedPassword

	errCreateUser := database.DB.Create(&newUser).Error
	if errCreateUser != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    newUser,
	})
}

// GET DATA BY ID
func UserHandlerGetById(ctx *fiber.Ctx) error {

	userId := ctx.Params("id")

	var user entity.User
	err := database.DB.First(&user, "id = ? ", userId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	//userResponse := response.UserResponse{
	//	ID:        user.ID,
	//	Name:      user.Name,
	//	Address:   user.Address,
	//	Phone:     user.Phone,
	//	CreatedAt: user.CreatedAt,
	//	UpdatedAt: user.UpdatedAt,
	//}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    user,
	})

}

// UPDATE
func UserHandlerUpdate(ctx *fiber.Ctx) error {
	var user entity.User

	userRequest := new(request.UserUpdateRequest)
	if err := ctx.BodyParser(userRequest); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	userId := ctx.Params("id")
	//CHECK AVAILABLE USER
	err := database.DB.First(&user, "id = ?", userId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	//UPDATE USER DATA
	if userRequest.Name != "" {
		user.Name = userRequest.Name
	}
	user.Address = userRequest.Address
	user.Phone = userRequest.Phone
	errUpdate := database.DB.Save(&user).Error
	if errUpdate != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    user,
	})

}

// EMAIL UPDATE
func UserHandlerUpdateEmail(ctx *fiber.Ctx) error {
	var user entity.User
	var isEmailUserExist entity.User

	userRequest := new(request.UserEmailRequest)
	if err := ctx.BodyParser(userRequest); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	userId := ctx.Params("id")

	//CHECK AVAILABLE USER
	err := database.DB.First(&user, "id = ?", userId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	//CHECK AVAILABLE EMAIL
	errCheckEmail := database.DB.First(&isEmailUserExist, "email = ?", userRequest.Email).Error
	if errCheckEmail == nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Email Sudah Ada",
		})
	}

	//UPDATE USER DATA
	user.Email = userRequest.Email

	errUpdate := database.DB.Save(&user).Error
	if errUpdate != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    user,
	})

}

// DELETE
func UserHandlerDelete(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	var user entity.User

	// CHECK AVAILABLE USER
	err := database.DB.Debug().First(&user, "id=?", userId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	errDelete := database.DB.Debug().Delete(&user).Error
	if errDelete != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "user was deleted",
	})
}
