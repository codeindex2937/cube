package alarm

import (
	"fmt"

	"cube/lib/context"
	"cube/lib/database"
)

type DeleteArgs struct {
	IDs []uint64 `arg:"positional"`
}

func (h *Alarm) delete(req *context.ChatContext, args *DeleteArgs) context.IResponse {
	alarms := []database.Alarm{}
	tx := h.DB.Find(&alarms, map[string]interface{}{
		"user_id":  req.UserID,
		"alarm_id": args.IDs,
	})
	if tx.Error != nil {
		return context.NewTextResponse(tx.Error.Error())
	}

	alarmIDs := []uint64{}
	for _, r := range alarms {
		alarmIDs = append(alarmIDs, r.AlarmID)
	}

	rowsAffected, err := DeleteAlarms(h.DB, h.Schedule, h.Event, alarmIDs)
	if err != nil {
		return context.NewTextResponse(err.Error())
	}

	return context.NewTextResponse(fmt.Sprintf("%v alarms deleted", rowsAffected))
}
