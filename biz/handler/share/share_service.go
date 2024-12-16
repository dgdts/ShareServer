// Code generated by hertz generator.

package share

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/dgdts/ShareServer/biz/model/share"
	biz_share "github.com/dgdts/ShareServer/biz/share"
	"github.com/dgdts/ShareServer/internal/response"
)

// GetShareNote .
// @router /api/v1/share/note [GET]
func GetShareNote(ctx context.Context, c *app.RequestContext) {
	response.JSON(ctx, c, biz_share.GetShareNote, func(ctx context.Context, c *app.RequestContext, req *share.GetShareNoteRequest) error {
		shareId := c.Param("share_id")
		if shareId == "" {
			return errors.New("share_id is required")
		}
		req.ShareId = shareId
		return nil
	})
}

// ListShareNoteComments .
// @router /api/v1/share/note/comments [GET]
func ListShareNoteComments(ctx context.Context, c *app.RequestContext) {
	var err error
	var req share.ListShareNoteCommentsRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(share.ListShareNoteCommentsResponse)

	c.JSON(consts.StatusOK, resp)
}

// CreateShareNoteComment .
// @router /api/v1/share/note/comment [POST]
func CreateShareNoteComment(ctx context.Context, c *app.RequestContext) {
	var err error
	var req share.CreateShareNoteCommentRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(share.CreateShareNoteCommentResponse)

	c.JSON(consts.StatusOK, resp)
}