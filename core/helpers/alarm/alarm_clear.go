package alarm

import (
	"fmt"

	"cube/lib/context"
	"cube/lib/database"
)

type ClearArgs struct {
	Dummy []string `arg:"positional"`
}

func (h *Alarm) clear(req *context.ChatContext, args *ClearArgs) context.Response {
	alarms := []database.Alarm{}
	tx := h.DB.Find(&alarms, map[string]interface{}{
		"user_id": req.UserID,
	})
	if tx.Error != nil {
		return context.Response(tx.Error.Error())
	}

	alarmIDs := []uint64{}
	for _, r := range alarms {
		alarmIDs = append(alarmIDs, r.AlarmID)
	}

	rowsAffected, err := DeleteAlarms(h.DB, h.Schedule, h.Event, alarmIDs)
	if err != nil {
		return context.Response(err.Error())
	}

	return context.Response(fmt.Sprintf("%v alarms deleted", rowsAffected))
}
