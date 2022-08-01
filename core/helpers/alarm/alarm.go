package alarm

import (
	"fmt"

	"cube/lib/context"
	"cube/lib/database"
	"cube/lib/utils"
	"cube/service/event"
	"cube/service/schedule"

	"github.com/alexflint/go-arg"
	"gorm.io/gorm"
)

type Args struct {
	Create *CreateArgs `arg:"subcommand:create"`
	List   *ListArgs   `arg:"subcommand:list"`
	Delete *DeleteArgs `arg:"subcommand:delete"`
	Clear  *ClearArgs  `arg:"subcommand:clear"`
}

type Alarm struct {
	DB       *gorm.DB
	Schedule schedule.IService
	Event    event.IService
	Time     *utils.TimeService
}

var config = arg.Config{Program: "alarm"}

const EventCreated = "alarm_created"
const EventDeleted = "alarm_deleted"

func (h *Alarm) Handle(req *context.ChatContext, args *Args) context.IResponse {
	switch {
	case args.Create != nil:
		return h.create(req, args.Create)
	case args.List != nil:
		return h.list(req, args.List)
	case args.Delete != nil:
		return h.delete(req, args.Delete)
	case args.Clear != nil:
		return h.clear(req, args.Clear)
	}
	return utils.PrintHelp("alarm", args)
}

func displayCurrentTask(schedule schedule.IService, alarm database.Alarm, ts *utils.TimeService) string {
	var nextSchedule string
	task, nextSched := schedule.SearchTask(alarm.AlarmID)
	if task == nil {
		nextSchedule = "invalid"
	} else {
		nextSchedule = ts.LocalTimeString(nextSched)
	}

	if alarm.RegID < 1 {
		return fmt.Sprintf(`ID=%v "%v" "%v" Next=%v`,
			alarm.AlarmID, alarm.Pattern, alarm.Message, nextSchedule)
	} else {
		return fmt.Sprintf(`ID=%v "%v" "%v" Next=%v (to %v)`,
			alarm.AlarmID, alarm.Pattern, alarm.Message, nextSchedule, alarm.RegID)
	}
}

func DeleteAlarms(db *gorm.DB, schedule schedule.IService, event event.IService, alarmIDs []uint64) (count int64, err error) {
	records := []database.Alarm{}
	tx := db.Find(&records, map[string]interface{}{
		"alarm_id": alarmIDs,
	})
	if tx.Error != nil {
		return 0, tx.Error
	}

	tx = db.Delete(records)
	if tx.Error != nil {
		return 0, tx.Error
	}

	IDs := []uint64{}
	for _, alarm := range records {
		IDs = append(IDs, alarm.AlarmID)
	}

	event.Publish(EventDeleted, IDs)

	return tx.RowsAffected, nil
}
