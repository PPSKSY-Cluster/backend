package mail

import (
	"fmt"
	"github.com/PPSKSY-Cluster/backend/db"
	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var scheduler = gocron.NewScheduler(time.Local)

var sendMail = func(to string, message string) {
	err := SendMail(to, message)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func ScheduleMail(r db.Reservation) error {
	RemoveIfExists(r.ID)

	user, err := db.GetUserById(r.UserID)
	if err != nil {
		return err
	}

	diff := float64(r.EndTime-r.StartTime) * 0.9
	t := time.Unix(r.StartTime+int64(diff), 0)

	_, err = scheduler.
		Every(1).
		Minute().
		StartAt(t).
		LimitRunsTo(1).
		Tag(r.ID.String()).
		Do(sendMail,
			user.Mail, //Params to SendMail are given here
			ReservationReminder(user.Mail, time.Unix(r.EndTime, 0)))
	if err != nil {
		return err
	}

	scheduler.StartAsync()
	return nil
}

func RemoveIfExists(id primitive.ObjectID) {
	jobs, _ := scheduler.FindJobsByTag(id.String())
	if len(jobs) > 0 {
		_ = scheduler.RemoveByTag(id.String())
	}
}
