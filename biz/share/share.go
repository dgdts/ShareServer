package share

import (
	"context"

	"github.com/dgdts/ShareServer/biz/model/share"
)

func GetShareNote(ctx context.Context, req *share.GetShareNoteRequest) ([]byte, error) {
	html, err := shareCache.Get(ctx, req.ShareId)
	if err != nil {
		return nil, err
	}

	return html, nil
}
