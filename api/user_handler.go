package api

import (
	"fmt"

	"github.com/thimc/hotel-backend/db"
	"github.com/thimc/hotel-backend/types"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		params types.UpdateUserParams
		id     = c.Params("id")
	)

	if err := c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}

	filter := map[string]any{"_id": id}
	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}

	return c.JSON(Response{
		Success: true,
		Message: fmt.Sprintf("updated %s", id),
	})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), id); err != nil {
		return err
	}
	return c.JSON(Response{
		Success: true,
		Message: fmt.Sprintf("user %s deleted", id),
	})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}

	if errs := params.Validate(); len(errs) > 0 {
		return ErrorBadRequest()
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return ErrorBadRequest()
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return ErrorNotFound("User")
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return ErrorNotFound("User")
	}
	return c.JSON(users)
}
