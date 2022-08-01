package registration

import (
	"fmt"

	"cube/core/helpers/alarm"
	"cube/lib/context"
	"cube/lib/database"
	"cube/lib/utils"
	"cube/service/event"
	"cube/service/schedule"

	"gorm.io/gorm"
)

type Args struct {
	Create *CreateArgs `arg:"subcommand:create"`
	List   *ListArgs   `arg:"subcommand:list"`
	Delete *DeleteArgs `arg:"subcommand:delete"`
	Clear  *ClearArgs  `arg:"subcommand:clear"`
}

type Reg struct {
	DB       *gorm.DB
	Schedule schedule.IService
	Event    event.IService
}

type RegisteredAlarm struct {
	AlarmID uint64
	RegID   uint64
}

func (h *Reg) Handle(req *context.ChatContext, args *Args) context.IResponse {
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

func displayRegistration(reg database.Registration) string {
	return fmt.Sprintf(`ID=%v "%v" Token=%v`,
		reg.RegID, reg.Description, reg.Token)
}

func listRegisteredAlarms(db *gorm.DB, userID string, regIDs []uint64) ([]RegisteredAlarm, error) {
	pairs := []RegisteredAlarm{}
	conditions := map[string]interface{}{
		"alarms.user_id": userID,
	}

	if len(regIDs) > 0 {
		conditions["alarms.reg_id"] = regIDs
	}

	tx := db.Model(&database.Alarm{}).Where(conditions).Select(
		"alarms.alarm_id, alarms.reg_id",
	).Scan(&pairs)

	return pairs, tx.Error
}

func deleteRegs(db *gorm.DB, schedule schedule.IService, event event.IService, userID string, ids []uint64) (alarmDeletedCount int64, err error) {
	pairs, err := listRegisteredAlarms(db, userID, ids)
	if err != nil {
		return 0, err
	}

	if len(pairs) > 0 {
		alarmIDs := []uint64{}
		for _, r := range pairs {
			alarmIDs = append(alarmIDs, r.AlarmID)
		}

		alarmDeletedCount, err = alarm.DeleteAlarms(db, schedule, event, alarmIDs)
		if err != nil {
			return alarmDeletedCount, err
		}
	}

	tx := db.Where(map[string]interface{}{
		"user_id": userID,
		"reg_id":  ids,
	}).Delete(&database.Registration{})
	if tx.Error != nil {
		return alarmDeletedCount, tx.Error
	}
	return alarmDeletedCount, nil
}
