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
		return
	}
}

func ScheduleMail(r db.Reservation) error {
	if nil == RemoveIfExists(r.ID) {
		fmt.Println("Deleted job")
	}

	user, err := db.GetUserById(r.UserID)
	if err != nil {
		return err
	}

	diff := float64(r.EndTime-r.StartTime) * 0.9
	t := time.Unix(r.StartTime+int64(diff), 0)

	fmt.Println(t)

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

func RemoveIfExists(id primitive.ObjectID) error {
	return scheduler.RemoveByTag(id.String())
}
