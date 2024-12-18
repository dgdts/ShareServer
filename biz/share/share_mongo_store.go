package share

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/dgdts/ShareServer/internal/utils"
	"github.com/dgdts/ShareServer/pkg/cache"
	"github.com/dgdts/UniversalServer/pkg/mongo"
)

const (
	ShareNoteCollection    = "share_notes"
	MarkdownNoteCollection = "markdown_notes"
)

const (
	IDField = "_id"
)

const (
	NoteTypeMarkdown = "markdown"
)

type ShareNoteStatus int

const (
	ShareNoteStatusDefault ShareNoteStatus = iota
	ShareNoteStatusExpired
	ShareNoteStatusCancel
)

type ShareNoteShareType int

const (
	ShareTypeDefault ShareNoteShareType = iota
	ShareTypeCanView
	ShareTypeCanComment
	ShareTypeCanEdit
)

type ShareNote struct {
	ID        string             `bson:"_id"`
	NoteID    string             `bson:"note_id"`
	UserID    string             `bson:"user_id"`
	NoteType  string             `bson:"note_type"`
	ShareType ShareNoteShareType `bson:"share_type"`
	ShareURL  string             `bson:"share_url"`
	ViewCount int                `bson:"view_count"`
	Status    ShareNoteStatus    `bson:"status"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type MarkdownNoteData struct {
	ID      string `bson:"_id"`
	Content string `bson:"content"`
}

var _ cache.CacheStore[[]byte] = (*ShareMongoStore)(nil)

type ShareMongoStore struct {
}

func NewShareMongoStore() *ShareMongoStore {
	return &ShareMongoStore{}
}

func (s *ShareMongoStore) validate(note *ShareNote) error {
	if note.ShareType == ShareTypeDefault {
		return errors.New("share type is default")
	}

	if note.Status == ShareNoteStatusCancel {
		return errors.New("share note is canceled")
	}

	if note.NoteType != NoteTypeMarkdown {
		return errors.New("share note is only support markdown type now")
	}

	return nil
}

func (s *ShareMongoStore) Get(ctx context.Context, key string) ([]byte, error) {
	r := mongo.Finder(utils.GlobalCollection(ShareNoteCollection)).FindOne(ctx, IDField, key)
	if r.Error() != nil {
		return nil, r.Error()
	}

	var shareNote ShareNote
	err := r.Read(&shareNote)
	if err != nil {
		return nil, err
	}

	if err := s.validate(&shareNote); err != nil {
		return nil, err
	}

	r = mongo.Finder(utils.GlobalCollection(MarkdownNoteCollection)).FindOne(ctx, IDField, shareNote.NoteID)
	if r.Error() != nil {
		return nil, r.Error()
	}

	var markdownNote MarkdownNoteData
	err = r.Read(&markdownNote)
	if err != nil {
		return nil, err
	}

	return json.Marshal(markdownNote)
}

func (s *ShareMongoStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return errors.New("share service does not support set")
}
