package alarm

import (
	"fmt"

	"cube/lib/context"
	"cube/lib/database"
	"cube/lib/utils"

	"github.com/alexflint/go-arg"
)

type CreateArgs struct {
	Channel     uint64 `arg:"--chan"`
	CronPattern string `arg:"positional"`
	Message     string `arg:"positional"`
}

func (h *Alarm) create(req *context.ChatContext, args *CreateArgs) context.IResponse {
	if len(args.CronPattern) < 1 {
		return utils.PrintUsage(arg.Config{Program: "create"}, args)
	}

	record := &database.Alarm{
		Pattern: args.CronPattern,
		Message: args.Message,
		UserID:  fmt.Sprintf("%v", req.UserID),
	}

	if args.Channel > 0 {
		var reg database.Registration
		tx := h.DB.First(&reg, map[string]interface{}{
			"reg_id":  args.Channel,
			"user_id": req.UserID,
		})
		if tx.Error != nil {
			return context.NewTextResponse(fmt.Sprintf("unknown channel(%v): %v", args.Channel, tx.Error))
		}

		record.RegID = args.Channel
	}

	tx := h.DB.Save(record)
	if tx.Error != nil {
		return context.NewTextResponse(tx.Error.Error())
	}

	h.Event.Publish(EventCreated, record)

	return context.NewTextResponse(displayCurrentTask(h.Schedule, *record))
}
