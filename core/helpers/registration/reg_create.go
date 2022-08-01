package registration

import (
	"cube/lib/context"
	"cube/lib/database"
	"fmt"
)

type CreateArgs struct {
	Description string   `arg:"positional"`
	Token       string   `arg:"positional"`
	Dummy       []string `arg:"positional"`
}

func (h *Reg) create(req *context.ChatContext, args *CreateArgs) context.IResponse {
	if 1 > len(args.Token) {
		return context.NewTextResponse("require token")
	}

	record := database.Registration{
		Token:       args.Token,
		Description: args.Description,
		UserID:      fmt.Sprintf("%v", req.UserID),
		UserName:    req.Username,
	}
	tx := h.DB.Create(&record)
	if tx.Error != nil {
		return context.NewTextResponse(tx.Error.Error())
	}

	if tx.RowsAffected > 0 {
		return context.NewTextResponse(displayRegistration(record))
	} else {
		return context.NewTextResponse("nothing changed")
	}
}
