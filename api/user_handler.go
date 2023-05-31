package api

import (
	"fmt"
	"hotel/api/errors"
	"hotel/db"
	"hotel/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		// values bson.M
		params types.UpdateUserParams
		id     = c.Params("id")
	)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.ErrorInvalidID()
	}

	if err := c.BodyParser(&params); err != nil {
		return errors.ErrorBadRequest()
	}

	filter := bson.M{"_id": oid}
	if err = h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}

	return c.JSON(GenericResp{
		Success: true,
		Msg:     fmt.Sprintf("updated %s", id),
	})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), id); err != nil {
		return err
	}
	return c.JSON(GenericResp{
		Success: true,
		Msg:     fmt.Sprintf("user %s deleted", id),
	})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return errors.ErrorBadRequest()
	}

	if errs := params.Validate(); len(errs) > 0 {
		return errors.ErrorBadRequest()
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
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
		return errors.ErrorNotFound(id)
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}
