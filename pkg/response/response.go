package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ybalcin/wallet-service/pkg/errr"
)

type Response struct {
	data interface{}
	err  *errr.Error

	c *fiber.Ctx
}

// New creates new instance of Response
func New(c *fiber.Ctx) *Response {
	return &Response{c: c}
}

// Data provides to set data of Response
func (r *Response) Data(data interface{}) *Response {
	r.data = data
	return r
}

// Error provides to set Error of Response
func (r *Response) Error(err *errr.Error) *Response {
	r.err = err
	return r
}

// JSON calls fiber.Ctx.JSON() and returns error
func (r *Response) JSON() error {
	if r.err != nil {
		return r.c.Status(r.err.Code).JSON(r.err)
	}

	return r.c.JSON(r.data)
}
