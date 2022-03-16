package registration

import (
	"bytes"

	"cube/lib/context"
	"cube/lib/database"
)

type ListArgs struct {
	Dummy []string `arg:"positional"`
}

func (h *Reg) list(req *context.ChatContext, args *ListArgs) context.Response {
	records := []database.Registration{}
	tx := h.DB.Find(&records, map[string]interface{}{"user_id": req.UserID})
	if tx.Error != nil {
		return context.Response(tx.Error.Error())
	}

	buf := new(bytes.Buffer)
	for _, r := range records {
		buf.Write([]byte(displayRegistration(r) + "\n"))
	}

	return context.Response(buf.String())
}
