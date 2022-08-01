package action

import (
	"cube/lib/context"
	"cube/lib/database"
)

type CreateArgs struct {
	Message string `arg:"positional"`
	Items   string `arg:"positional"`
}

func (h *Action) create(req *context.ChatContext, args *CreateArgs) context.IResponse {
	if 1 > len(args.Items) {
		return context.NewTextResponse("require items")
	}

	record := database.Action{
		UserID:  req.UserID,
		Message: args.Message,
		Items:   args.Items,
	}
	tx := h.DB.Create(&record)
	if tx.Error != nil {
		return context.NewTextResponse(tx.Error.Error())
	}

	if tx.RowsAffected > 0 {
		return context.NewTextResponse(displayAction(record))
	} else {
		return context.NewTextResponse("nothing changed")
	}
}
