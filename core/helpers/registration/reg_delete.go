package registration

import (
	"fmt"

	"cube/lib/context"
	"cube/lib/utils"

	"github.com/alexflint/go-arg"
)

type DeleteArgs struct {
	IDs []uint64 `arg:"positional"`
}

func (h *Reg) delete(req *context.ChatContext, args *DeleteArgs) context.IResponse {
	if len(args.IDs) < 1 {
		return utils.PrintUsage(arg.Config{Program: "delete"}, args)
	}

	alarmDeletedCount, err := deleteRegs(h.DB, h.Schedule, h.Event, fmt.Sprintf("%v", req.UserID), args.IDs)
	if err != nil {
		return context.NewTextResponse(err.Error())
	}

	return context.NewTextResponse(fmt.Sprintf("%v alarms deleted", alarmDeletedCount))
}
