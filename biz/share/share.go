package share

import (
	"context"

	"github.com/dgdts/ShareServer/biz/model/share"
)

type ShareNoteResponse struct {
	Data []byte `json:"data"`
}

func GetShareNote(ctx context.Context, req *share.GetShareNoteRequest) (*ShareNoteResponse, error) {
	note, err := shareCache.Get(ctx, req.ShareId)
	if err != nil {
		return nil, err
	}

	return &ShareNoteResponse{
		Data: note,
	}, nil
}
