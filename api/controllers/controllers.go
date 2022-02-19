package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/toshkentov01/alif-tech-task/api-gateway/api/models"
	pb "github.com/toshkentov01/alif-tech-task/api-gateway/genproto/user-service"
	client "github.com/toshkentov01/alif-tech-task/api-gateway/grpc_client"
	"github.com/toshkentov01/alif-tech-task/api-gateway/pkg/utils"
)

// SignUpFully ...
// @Description Creates an identified user.
// @Summary creates an identified user
// @Tags register
// @Accept json
// @Produce json
// @Param signup body models.SignUpModel true "Sign Up"
// @Success 200 {object} models.SignUpResponseModel
// @Failure 400 {object} models.StandardErrorModel
// @Failure 404 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /create-identified-user/ [post]
func SignUpFully(c *fiber.Ctx) error {
	var (
		body                      models.SignUpModel
		accessToken, refreshToken string
	)

	err := c.BodyParser(&body)
	if err != nil {
		log.Println("Error parsing body: ", err)
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: err.Error(),
		})
	}

	err = body.Validate()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: err.Error(),
		})
	}

	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	body.Username = strings.TrimSpace(body.Username)
	body.Username = strings.ToLower(body.Username)

	id, err := uuid.NewRandom()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: "Error while generating id",
		})
	}

	//Creating access and refresh tokens
	tokens, err := utils.GenerateNewTokens(id.String(), map[string]string{"role": "user"})

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: "Error while generating tokens",
		})
	}
	accessToken, refreshToken = tokens.Access, tokens.Refresh

	result, err := client.UserService().CheckFields(context.Background(), &pb.CheckfieldsRequest{
		Username: body.Username,
		Email:    body.Email,
	})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.StandardErrorModel{
			ErrorMessage: err.Error(),
		})
	}

	if result.EmailExists {
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: "User with this email already exists",
		})
	} else if result.UsernameExists {
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: "User with this username already exists",
		})
	}

	_, serviceError := client.UserService().CreateIdentifiedUser(context.Background(), &pb.CreateIdentifiedUserRequest{
		Id:           id.String(),
		Username:     body.Username,
		FullName:     body.FullName,
		Email:        body.Email,
		Password:     body.Password,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	if serviceError != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.StandardErrorModel{
			ErrorMessage: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(models.SignUpResponseModel{
		UserID:       id.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// CreateUnidentifiedUser ...
// @Description Creates an unidentified user.
// @Summary creates an unidentified user
// @Tags register
// @Accept json
// @Produce json
// @Param signup body models.SignUpModelForUnidentifiedUser true "Sign Up"
// @Success 200 {object} models.SignUpResponseModelForUnidentifiedUser
// @Failure 400 {object} models.StandardErrorModel
// @Failure 404 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /create-unidentified-user/ [post]
func CreateUnidentifiedUser(c *fiber.Ctx) error {
	var (
		body models.SignUpModelForUnidentifiedUser
	)

	err := c.BodyParser(&body)
	if err != nil {
		log.Println("Error parsing body: ", err)
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: err.Error(),
		})
	}

	err = body.Validate()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: err.Error(),
		})
	}

	body.Username = strings.TrimSpace(body.Username)
	body.Username = strings.ToLower(body.Username)

	id, err := uuid.NewRandom()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: "Error while generating id",
		})
	}

	result, err := client.UserService().CheckFields(context.Background(), &pb.CheckfieldsRequest{
		Username: body.Username,
		Email:    "",
	})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.StandardErrorModel{
			ErrorMessage: err.Error(),
		})
	}

	if result.UsernameExists {
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: "User with this username already exists",
		})
	}

	//Creating access and refresh tokens
	tokens, err := utils.GenerateNewTokens(id.String(), map[string]string{"role": "user"})

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: "Error while generating tokens",
		})
	}
	accessToken, refreshToken := tokens.Access, tokens.Refresh

	_, serviceError := client.UserService().CreateUnIdentifiedUser(context.Background(), &pb.CreateUnIdentifiedUserRequest{
		Id:           id.String(),
		Username:     body.Username,
		Password:     body.Password,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	if serviceError != nil {
		log.Println("Error while creating unidentified user. ERROR: ", err)
		return c.Status(http.StatusInternalServerError).JSON(models.StandardErrorModel{
			ErrorMessage: "Internal Server Error",
		})
	}

	return c.Status(http.StatusOK).JSON(models.SignUpResponseModelForUnidentifiedUser{
		ID:           id.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// CheckUserAccount ...
// @Description CheckUserAccount API checks whether user has an account or not.
// @Tags user
// @Accept json
// @Produce json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Success 200 {object} models.CheckUserAccountResponseModel
// @Failure 400 {object} models.StandardErrorModel
// @Failure 404 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /check-user-account/ [get]
func CheckUserAccount(c *fiber.Ctx) error {

	result, err := client.UserService().CheckUserAccount(context.Background(), &pb.CheckUserAccountRequest{
		Username: c.Query("username"),
		Password: c.Query("password"),
	})

	if err != nil {
		log.Println("Error while checking user account. Error: ", err)
		return c.Status(http.StatusInternalServerError).JSON(models.StandardErrorModel{
			ErrorMessage: "Internal Server Error",
		})
	}

	if result.Exists {
		return c.Status(http.StatusOK).JSON(models.CheckUserAccountResponseModel{
			Exists: true,
		})
	}

	return c.Status(http.StatusOK).JSON(models.CheckUserAccountResponseModel{
		Exists: false,
	})
}

// Income ...
// @Description Income API used for topping up a balance.
// @Security ApiKeyAuth
// @Tags user
// @Accept json
// @Produce json
// @Param income body models.IncomeModel true "Income"
// @Success 200 {object} models.Success
// @Failure 400 {object} models.StandardErrorModel
// @Failure 404 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /user/income/ [post]
func Income(c *fiber.Ctx) error {
	var (
		body models.IncomeModel
	)

	err := c.BodyParser(&body)
	if err != nil {
		log.Println("Error parsing body: ", err)
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: err.Error(),
		})
	}

	user, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		log.Println("Error taking user id! ", err)
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: "Failed to extract id from token",
		})
	}

	_, serviceErr := client.UserService().Income(context.Background(), &pb.IncomeRequest{
		UserId:       user.UserID.String(),
		IncomeAmount: body.IncomeAmount,
	})

	st, ok := status.FromError(serviceErr)
	if !ok || st.Code() == codes.Internal {
		return c.Status(http.StatusInternalServerError).JSON(models.StandardErrorModel{
			ErrorMessage: "Internal Server Error",
		})
	} else if st.Code() == codes.PermissionDenied {
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: "Permission Denied. If you top up balance with this amount, your balance will be above the maximum allowed cash",
		})
	}

	return c.Status(http.StatusOK).JSON(models.Success{
		Success: true,
	})
}

// Expense ...
// @Description Expense API used for reducing a balance.
// @Security ApiKeyAuth
// @Tags user
// @Accept json
// @Produce json
// @Param income body models.ExpenseModel true "Income"
// @Success 200 {object} models.Success
// @Failure 400 {object} models.StandardErrorModel
// @Failure 404 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /user/expense/ [post]
func Expense(c *fiber.Ctx) error {
	var (
		body models.ExpenseModel
	)

	err := c.BodyParser(&body)
	if err != nil {
		log.Println("Error parsing body: ", err)
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: err.Error(),
		})
	}

	user, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		log.Println("Error taking user id! ", err)
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: "Failed to extract id from token",
		})
	}

	_, serviceErr := client.UserService().Expense(context.Background(), &pb.ExpenseRequest{
		UserId:        user.UserID.String(),
		ExpenseAmount: body.ExpenseAmount,
	})

	st, ok := status.FromError(serviceErr)
	if !ok || st.Code() == codes.Internal {
		return c.Status(http.StatusInternalServerError).JSON(models.StandardErrorModel{
			ErrorMessage: "Internal Server Error",
		})
	} else if st.Code() == codes.PermissionDenied {
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: "Permission Denied. If you reduce a balance with this amount, your balance will be under the minimum allowed cash",
		})
	}

	return c.Status(http.StatusOK).JSON(models.Success{
		Success: true,
	})
}

// GetBalance ...
// @Description GetBalance API used for getting a user balance.
// @Security ApiKeyAuth
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.GetBalanceResponseModel
// @Failure 400 {object} models.StandardErrorModel
// @Failure 404 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /user/balance/ [get]
func GetBalance(c *fiber.Ctx) error {

	user, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		log.Println("Error taking user id! ", err)
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: "Failed to extract id from token",
		})
	}

	result, serviceErr := client.UserService().GetBalance(context.Background(), &pb.GetBalanceRequest{
		UserId: user.UserID.String(),
	})

	if serviceErr != nil {
		log.Println("Error while getting a balance, error: ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(models.StandardErrorModel{
			ErrorMessage: "Internal Server Errr",
		})
	}

	return c.Status(http.StatusOK).JSON(models.GetBalanceResponseModel{
		Balance: result.Balance,
	})
}

// ListOperationsByType ...
// @Description GetBalance API used for getting a user balance.
// @Security ApiKeyAuth
// @Tags user
// @Accept json
// @Produce json
// @Param OperationType header string false "OperationType" Enums(income_operations, expense_operations)
// @Success 200 {object} models.ListOperationsByTypeResponseModel
// @Failure 400 {object} models.StandardErrorModel
// @Failure 404 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /user/operations/ [get]
func ListOperationsByType(c *fiber.Ctx) error {

	user, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		log.Println("Error taking user id! ", err)
		return c.Status(http.StatusBadRequest).JSON(models.StandardErrorModel{
			ErrorMessage: "Failed to extract id from token",
		})
	}

	operaionType := c.Get("OperationType")

	result, serviceErr := client.UserService().ListTotalOperationsByType(context.Background(), &pb.ListTotalOperationsByTypeRequest{
		UserId:        user.UserID.String(),
		OperationType: operaionType,
	})

	if serviceErr != nil {
		log.Println("Error while getting a balance, error: ", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(models.StandardErrorModel{
			ErrorMessage: "Internal Server Errr",
		})
	}

	return c.Status(http.StatusOK).JSON(result)
}
