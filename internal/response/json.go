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
