package action

import (
	"cube/lib/context"
)

type ShowArgs struct {
	ActionID uint64 `arg:"positional"`
}

func (h *Action) show(req *context.ChatContext, args *ShowArgs) context.IResponse {
	resp, err := NewActionResponse(h.DB, args.ActionID, req.UserID)
	if err != nil {
		return context.NewErrorResponse(err)
	}

	resp.Title = "action"
	return resp
}
