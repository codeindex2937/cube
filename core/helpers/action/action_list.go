package action

import (
	"cube/lib/context"
	"cube/lib/database"
	"strings"
)

type ListArgs struct {
	Dummy []uint64 `arg:"positional"`
}

func (h *Action) list(req *context.ChatContext, args *ListArgs) context.IResponse {
	records := []database.Action{}
	tx := h.DB.Find(&records, map[string]interface{}{
		"user_id": req.UserID,
	})
	if tx.Error != nil {
		return context.NewErrorResponse(tx.Error)
	}

	var b strings.Builder
	for _, r := range records {
		b.WriteString(displayAction(r))
		b.WriteString("\n")
	}

	return context.NewTextResponse(b.String())
}
