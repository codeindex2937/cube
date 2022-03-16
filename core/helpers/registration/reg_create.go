package registration

import (
	"cube/lib/context"
	"cube/lib/database"
)

type CreateArgs struct {
	Description string   `arg:"positional"`
	Token       string   `arg:"positional"`
	Dummy       []string `arg:"positional"`
}

func (h *Reg) create(req *context.ChatContext, args *CreateArgs) context.Response {
	if 1 > len(args.Token) {
		return "require token"
	}

	record := database.Registration{
		Token:       args.Token,
		Description: args.Description,
		UserID:      req.UserID,
		UserName:    req.UserName,
	}
	tx := h.DB.Create(&record)
	if tx.Error != nil {
		return context.Response(tx.Error.Error())
	}

	if tx.RowsAffected > 0 {
		return context.Response(displayRegistration(record))
	} else {
		return context.Response("nothing changed")
	}
}
