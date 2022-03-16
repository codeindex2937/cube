package alarm

import (
	"bytes"
	"fmt"

	"cube/lib/context"
	"cube/lib/database"
)

type ListArgs struct {
	Dummy []string `arg:"positional"`
}

func (h *Alarm) list(req *context.ChatContext, args *ListArgs) context.Response {
	records := []database.Alarm{}

	tx := h.DB.Where(
		"user_id=? OR user_id LIKE ?", req.UserID, fmt.Sprintf("%v_%%", req.UserID),
	).Find(&records)
	if tx.Error != nil {
		return context.Response(tx.Error.Error())
	}

	buf := new(bytes.Buffer)
	for _, r := range records {
		buf.Write([]byte(displayCurrentTask(h.Schedule, r) + "\n"))
	}

	return context.Response(buf.String())
}
