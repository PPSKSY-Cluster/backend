package mail

import (
	"fmt"
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

	fmt.Println("Initialized schedule")

	return nil
}

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
	var keyList []int64
	for _, r := range reservations {
		for start := time.Now().Unix(); start < time.Now().AddDate(0, 3, 0).Unix(); start += 86400 { //1 Day = 86400 seconds
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

	//It might be worthwhile to keep track of the amount of usable nodes over the observed available time
	sort.Slice(keyList, func(i, j int) bool { return keyList[i] < keyList[j] })
	t := 0
	for _, k := range keyList {
		if clusterReservations[k] <= cluster.Nodes-2 { //Assumes that 2 nodes warrant a reservation
			t += 1
		} else {
			t = 0
			continue
		}
		if t >= 3 {
			fmt.Println("Sending Notification!")
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
