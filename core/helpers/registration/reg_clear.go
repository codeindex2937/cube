package registration

import (
	"fmt"

	"cube/lib/context"
	"cube/lib/database"
)

type ClearArgs struct {
	Dummy []string `arg:"positional"`
}

func (h *Reg) clear(req *context.ChatContext, args *ClearArgs) context.Response {
	regs := []database.Registration{}
	tx := h.DB.Find(&regs, map[string]interface{}{"user_id": req.UserID})
	if tx.Error != nil {
		return context.Response(tx.Error.Error())
	}

	regIDs := []uint64{}
	for _, r := range regs {
		regIDs = append(regIDs, r.RegID)
	}

	alarmDeletedCount, err := deleteRegs(h.DB, h.Schedule, h.Event, req.UserID, regIDs)
	if err != nil {
		return context.Response(err.Error())
	}

	return context.Response(fmt.Sprintf("%v registrations deleted, %v alarms deleted",
		len(regIDs), alarmDeletedCount))
}
