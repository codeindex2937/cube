package core

import (
	"strconv"
	"strings"

	"cube/core/helpers/action"
	"cube/core/helpers/alarm"
	"cube/core/helpers/food"
	"cube/core/helpers/registration"
	"cube/lib/context"
	"cube/lib/database"
	"cube/lib/fake"
	"cube/lib/logger"
	"cube/lib/utils"
	"cube/service/event"
	"cube/service/message"
	"cube/service/schedule"

	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"
	"github.com/google/shlex"
	"gorm.io/gorm"
)

type rootArgs struct {
	Alarm   *alarm.Args        `arg:"subcommand:alarm"`
	Reg     *registration.Args `arg:"subcommand:reg"`
	Shuffle *ShuffleArgs       `arg:"subcommand:shuffle"`
	Food    *food.Args         `arg:"subcommand:food"`
	Action  *action.Args       `arg:"subcommand:action"`
}

type Core struct {
	Schedule schedule.IService
	DB       *gorm.DB
	Alarm    *alarm.Alarm
	Reg      *registration.Reg
	Shuffle  *Shuffle
	Food     *food.Food
	Action   *action.Action
}

var log = logger.Log

func PrintHelp() context.IResponse {
	return utils.PrintHelp("", &rootArgs{})
}

func subscribeEvents(c *Core, e event.IService) {
	e.Subscribe(alarm.EventCreated, func(event string, ctx interface{}) {
		record := ctx.(*database.Alarm)
		sched, err := schedule.Parse(record.Pattern)
		if err != nil {
			logger.Log.Error(err.Error())
		}

		userID, _ := strconv.Atoi(record.UserID)
		c.Schedule.AddTask(&schedule.Task{
			Sched: sched,
			ID:    record.AlarmID,
			Run:   func() { c.SendMessage(userID, record.RegID, record.Message) },
		})
	})

	e.Subscribe(alarm.EventDeleted, func(event string, ctx interface{}) {
		IDs := ctx.([]uint64)

		c.Schedule.RemoveTasks(IDs)
	})
}

func New(
	db *gorm.DB,
	scheduleService schedule.IService,
) *Core {
	e := event.NewService()
	c := &Core{
		Schedule: scheduleService,
		DB:       db,
		Alarm: &alarm.Alarm{
			DB:       db,
			Schedule: scheduleService,
			Event:    e,
		},
		Reg: &registration.Reg{
			DB:       db,
			Schedule: scheduleService,
			Event:    e,
		},
		Shuffle: new(Shuffle),
		Food: &food.Food{
			DB: db,
		},
		Action: &action.Action{
			DB: db,
		},
	}

	subscribeEvents(c, e)

	return c
}

func NewFake() *Core {
	e := event.NewService()
	scheduleService := fake.NewScheduleService()
	db, err := database.New(":memory:")
	if err != nil {
		log.Errorf("NewFakeCore: %v\n", err)
		return nil
	}

	c := &Core{
		Schedule: scheduleService,
		DB:       db,
		Alarm: &alarm.Alarm{
			DB:       db,
			Schedule: scheduleService,
			Event:    e,
		},
		Reg: &registration.Reg{
			DB:       db,
			Schedule: scheduleService,
			Event:    e,
		},
		Food: &food.Food{
			DB: db,
		},
	}

	subscribeEvents(c, e)

	return c
}

func (c *Core) InitAlarms() error {
	alarms := []database.Alarm{}
	tx := c.DB.Find(&alarms)
	if tx.Error != nil {
		return tx.Error
	}

	for _, alarm := range alarms {
		record := alarm
		sched, err := schedule.Parse(record.Pattern)
		if err != nil {
			return err
		}
		c.Schedule.AddTask(&schedule.Task{
			Sched: sched,
			ID:    record.AlarmID,
			Run: func() {
				userID, _ := strconv.Atoi(record.UserID)
				c.SendMessage(userID, record.RegID, record.Message)
			},
		})
	}
	return nil
}

func (c *Core) Handle(req context.ChatContext, args []string) context.IResponse {
	root := &rootArgs{}
	p, err := arg.NewParser(arg.Config{}, root)
	if err != nil {
		return context.NewErrorResponse(err)
	}

	err = p.Parse(args)
	if err != nil {
		return utils.PrintHelp("", root)
	}

	switch {
	case root.Alarm != nil:
		return c.Alarm.Handle(&req, root.Alarm)
	case root.Reg != nil:
		return c.Reg.Handle(&req, root.Reg)
	case root.Shuffle != nil:
		return c.Shuffle.Handle(&req, root.Shuffle)
	case root.Food != nil:
		return c.Food.Handle(&req, root.Food)
	case root.Action != nil:
		return c.Action.Handle(&req, root.Action)
	}

	return utils.PrintHelp("", root)
}

func (p *Core) SendMessage(userID int, regID uint64, m string) {
	if !strings.HasPrefix(m, "/") {
		message.Service().Send(p.DB, userID, regID, gin.H{
			"text": m,
		})
	} else {
		args, err := shlex.Split(m[1:])
		if err != nil {
			logger.Log.Errorf("unknown command: %v", m)
			return
		}

		resp := p.Handle(context.ChatContext{
			UserID: userID,
		}, args)
		message.Service().Send(p.DB, userID, regID, resp.Normalize())
	}
}

func (p *Core) CreateReg(req context.ChatContext, description, token string) context.IResponse {
	return p.Handle(req, []string{"reg", "create", description, token})
}

func (p *Core) ListReg(req context.ChatContext) context.IResponse {
	return p.Handle(req, []string{"reg", "list"})
}

func (p *Core) DeleteReg(req context.ChatContext, channels ...uint64) context.IResponse {
	ids := []string{}
	for _, c := range channels {
		ids = append(ids, strconv.FormatUint(c, 10))
	}

	return p.Handle(req, append([]string{"reg", "delete"}, ids...))
}

func (p *Core) ClearReg(req context.ChatContext) context.IResponse {
	return p.Handle(req, []string{"reg", "clear"})
}

func (c *Core) CreateAlarmToChannel(req context.ChatContext, channel uint64, pattern, message string) context.IResponse {
	return c.Handle(req, []string{"alarm", "create", "--chan", strconv.FormatUint(channel, 10), pattern, message})
}

func (c *Core) CreateAlarm(req context.ChatContext, pattern, message string) context.IResponse {
	return c.Handle(req, []string{"alarm", "create", pattern, message})
}

func (c *Core) DeleteAlarm(req context.ChatContext, IDs ...string) context.IResponse {
	return c.Handle(req, append([]string{"alarm", "delete"}, IDs...))
}

func (c *Core) ClearAlarm(req context.ChatContext) context.IResponse {
	return c.Handle(req, []string{"alarm", "clear"})
}

func (c *Core) ListAlarm(req context.ChatContext) context.IResponse {
	return c.Handle(req, []string{"alarm", "list"})
}

func (c *Core) ListFood(req context.ChatContext) context.IResponse {
	return c.Handle(req, []string{"food", "list"})
}

func (c *Core) SetFoodTag(req context.ChatContext, foodID, foodTagID uint64) context.IResponse {
	return c.Handle(req, []string{"food", "attach_tag", strconv.FormatUint(foodID, 10), strconv.FormatUint(foodTagID, 10)})
}
