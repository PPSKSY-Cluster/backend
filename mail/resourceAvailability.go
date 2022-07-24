package mail

import (
	"github.com/PPSKSY-Cluster/backend/db"
	"github.com/go-co-op/gocron"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

var s = gocron.NewScheduler(time.Local)

var check = func() {
	if err := checkNotifications(); err != nil {
		log.Println(err.Error())
	}
}

func InitSchedule() error {
	every, err := strconv.Atoi(os.Getenv("NOTIFICATION_INTERVAL"))
	if err != nil {
		return err
	}

	_, err = s.Every(every).
		Second().
		Do(check)
	if err != nil {
		return err
	}

	s.StartAsync()

	return nil
}

//Gets called from the schedule
func checkNotifications() error {
	notifications, err := db.GetAllNotifications()
	if err != nil {
		return err
	}

	for _, n := range notifications {
		if err := checkAvailability(n); err != nil {
			return err
		}
	}

	return nil
}

func checkAvailability(notification db.ReservationNotification) error {
	user, err := db.GetUserById(notification.UserID)
	if err != nil {
		return err
	}

	cluster, err := db.GetCResourceById(notification.ClusterID)
	if err != nil {
		return err
	}

	reservations, err := db.GetReservationsByClusterId(notification.ClusterID)
	if err != nil {
		return err
	}

	clusterReservations := make(map[int64]int64) //Maps startTime to nodes used
	var keyList []int64                          //Keylist for keeping correct order of days
	for _, r := range reservations {

		//Iterate over a span of 3 Months from now  (Increments one day)
		for start := time.Now().Unix(); start < time.Now().AddDate(0, 3, 0).Unix(); start += 86400 {
			if _, ok := clusterReservations[start]; !ok {
				keyList = append(keyList, start)
			}
			if start < r.EndTime {
				clusterReservations[start] += r.Nodes
			} else {
				clusterReservations[start] += 0
			}
		}
	}

	sort.Slice(keyList, func(i, j int) bool { return keyList[i] < keyList[j] })
	t := 0 //The number of successive days with free nodes
	for _, k := range keyList {
		if clusterReservations[k] <= cluster.Nodes-2 { //Assumes that 2 nodes warrant a reservation
			t += 1
		} else {
			t = 0
			continue
		}
		if t >= 3 { //If number of successive days >= 3 -> Send notification
			message := ReservationNotification(user.Mail)
			if err := SendMail(user.Mail, message); err != nil {
				return err
			}

			if err := db.DeleteNotification(notification.ID); err != nil {
				return err
			}
			break
		}
	}

	return nil
}
