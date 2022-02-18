package schedule

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type scheduleService interface {
	ScheduleJob(ctx context.Context, job Job) error
}

type Handler struct {
	scheduleService scheduleService
}

func NewHandler(scheduleService scheduleService) *Handler {
	return &Handler{scheduleService: scheduleService}
}

func (h *Handler) RegisterFastHTTPRouters(a fiber.Router) {
	a.Post("/api/v1/schedule-job", h.scheduleJob)
}

type scheduleJobArgs struct {
	Timestamp int64  `json:"timestamp"`
	QueueID   string `json:"queue_id"`
	Action    string `json:"action"`
}

func (h *Handler) scheduleJob(c *fiber.Ctx) error {
	var args scheduleJobArgs
	if err := c.BodyParser(&args); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if args.Timestamp == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "missing timestamp field")
	}
	if args.Action == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing action field")
	}

	if err := h.scheduleService.ScheduleJob(c.UserContext(), Job{
		DateTime: time.Unix(args.Timestamp, 0),
		QueueID:  args.QueueID,
		Action:   args.Action,
	}); err != nil {
		return errors.Wrap(err, "schedule job")
	}

	return nil
}
