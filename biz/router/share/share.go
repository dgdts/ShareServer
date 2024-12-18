// Code generated by hertz generator. DO NOT EDIT.

package share

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	share "github.com/dgdts/ShareServer/biz/handler/share"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	{
		_share := root.Group("/share", _shareMw()...)
		_share.GET("/:share_id", append(_getsharenoteMw(), share.GetShareNote)...)
		_share_id := _share.Group("/:share_id", _share_idMw()...)
		_share_id.POST("/comment", append(_createsharenotecommentMw(), share.CreateShareNoteComment)...)
		_share_id.GET("/comments", append(_listsharenotecommentsMw(), share.ListShareNoteComments)...)
	}
}
