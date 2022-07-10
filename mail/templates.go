package mail

import (
	"os"
	"time"
)

//TODO: Add link to reservation of desired cluster	fmt.Println("router")

func ReservationNotification(to string) string {
	return "To: " + to + ", \r\n" +
		"From: " + os.Getenv("RELAY_USERNAME") + ", \r\n" +
		"Subject: Cluster Reservierung wieder möglich \r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		"<head>    " +
		"<style>" +
		"        div {" +
		"            text-align: center;" +
		"        }" +
		"        h1 {" +
		"            padding: 15px;" +
		"        }" +
		"        a {" +
		"            text-decoration: none;" +
		"            padding: 5px;" +
		"            color: black;" +
		"            background-color: red;" +
		"            border-radius: 8px;" +
		"        }" +
		"    </style>" +
		"</head>" +
		"<body>" +
		"    <div>" +
		"        <h1>Das von ihnen gewünschte Cluster ist wieder verfügbar!</h1>" +
		"        <p>Klicken sie hier um das Cluster zu reservieren:</p>" +
		"        <br>" +
		"        <a href=\"" + os.Getenv("CLIENT_URL") + "\">Reservieren</a>" +
		"    </div>" +
		"</body>"
}

func ReservationReminder(to string, time time.Time) string {
	return "To: " + to + ", \r\n" +
		"From: " + os.Getenv("RELAY_USERNAME") + ", \r\n" +
		"Subject: Cluster Reservierung läuft bald aus! \r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		"<head>    " +
		"<style>" +
		"        div {" +
		"            text-align: center;" +
		"        }" +
		"        h1 {" +
		"            padding: 15px;" +
		"        }" +
		"        a {" +
		"            text-decoration: none;" +
		"            padding: 5px;" +
		"            color: black;" +
		"            background-color: red;" +
		"            border-radius: 8px;" +
		"        }" +
		"    </style>" +
		"</head>" +
		"<body>" +
		"    <div>" +
		"        <h1>Ihre Cluster Reservierung läuft um " + time.String() + " aus</h1>" +
		"        <p>Falls Sie ihre Reservierung verlängern wollen, oder die Resource nicht mehr weiter benötigen, klicken Sie bitte hier;:</p>" +
		"        <br>" +
		"        <a href=\"" + os.Getenv("CLIENT_URL") + "\">Reservierung ändern</a>" +
		"    </div>" +
		"</body>"
}
