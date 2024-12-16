package response

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

type JSONHandler[Req any, Res any] func(c context.Context, req *Req) (*Res, error)
type JSONBeforeHandler[Req any] func(ctx context.Context, c *app.RequestContext, req *Req) error

func JSONError(c *app.RequestContext, err error) {
	resp := NewResultFromError(err)
	c.JSON(http.StatusOK, resp)
}

func JSONSuccess(c *app.RequestContext, data interface{}) {
	resp := NewResultWithData(data)
	c.JSON(http.StatusOK, resp)
}

func JSON[Req any, Res any](ctx context.Context, c *app.RequestContext, handler JSONHandler[Req, Res], before ...JSONBeforeHandler[Req]) {
	var req Req
	err := c.BindAndValidate(&req)
	if err != nil {
		JSONError(c, err)
		return
	}

	for _, beforeHandler := range before {
		err = beforeHandler(ctx, c, &req)
		if err != nil {
			JSONError(c, err)
			return
		}
	}

	resp, err := handler(ctx, &req)
	if err != nil {
		JSONError(c, err)
		return
	}

	JSONSuccess(c, resp)
}
