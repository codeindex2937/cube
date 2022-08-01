package action

import (
	"fmt"

	"cube/lib/context"
	"cube/lib/database"
)

type DeleteArgs struct {
	IDs []uint64 `arg:"positional"`
}

func (h *Action) delete(req *context.ChatContext, args *DeleteArgs) context.IResponse {
	actions := []database.Action{}
	tx := h.DB.Find(&actions, map[string]interface{}{
		"user_id":   req.UserID,
		"action_id": args.IDs,
	})
	if tx.Error != nil {
		return context.NewTextResponse(tx.Error.Error())
	}

	tx = h.DB.Delete(actions)
	if tx.Error != nil {
		return context.NewErrorResponse(tx.Error)
	}
	return context.NewTextResponse(fmt.Sprintf("%v actions deleted", tx.RowsAffected))
}
