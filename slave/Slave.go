// slave
package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
)

func main() {
	// Connect to the master server
	conn, err := net.Dial("tcp", "192.168.15.250:9050")
	if err != nil {
		fmt.Println("Error connecting to master server:", err)
		return
	}
	defer conn.Close()

	var userInput, masterResponse string
	sendRequest("login", conn)
	var adminName, adminPassword string

	fmt.Println("--->> Please Enter your Admin Name : ")
	fmt.Scan(&adminName)
	sendRequest(adminName, conn)

	fmt.Println("--->> Please Enter your password : ")
	fmt.Scan(&adminPassword)
	sendRequest(adminPassword, conn)

	masterResponse, err = receiveResponse(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if masterResponse != "1" {
		fmt.Println(masterResponse)
	}
	if masterResponse == "1" {
		fmt.Println("--------------------Successful login-------------------------")
		fmt.Println("--------------------------Welcome to Hotel Management-------------------------------")
		for {

			fmt.Println("--->> Press 1 to Clients. ")
			fmt.Println("--->> Press 2 to Rooms. ")
			fmt.Println("--->> Press 3 to Reservations. ")
			fmt.Println("--->> Press 4 to Logout. ")

			fmt.Scan(&userInput)

			if userInput == "1" {
				for {
					fmt.Println("--->> Please Select what you want to do.")
					fmt.Println("--->> Press 1 to Add a Client.")
					fmt.Println("--->> Press 2 to Edit an Existing Client.")
					fmt.Println("--->> Press 3 to Delete a Client.")
					fmt.Println("--->> Press 4 to Show all Clients.")
					fmt.Println("--->> Press 5 to Go Backward")
					var commandResponse string
					fmt.Scan(&userInput)
					if userInput == "5" {
						break
					}
					commandResponse = Clients(userInput, conn)
					fmt.Println(commandResponse)

				}

			} else if userInput == "2" {
				for {
					fmt.Println("--->> Please Select what you want to do.")
					fmt.Println("--->> Press 1 to Add a Room.")
					fmt.Println("--->> Press 2 to Edit an Existing Room.")
					fmt.Println("--->> Press 3 to Delete a Room.")
					fmt.Println("--->> Press 4 to Show all Rooms.")
					fmt.Println("--->> Press 5 to Go Backward")
					var commandResponse string
					fmt.Scan(&userInput)
					if userInput == "5" {
						break
					}
					commandResponse = Rooms(userInput, conn)
					fmt.Println(commandResponse)
				}

			} else if userInput == "3" {
				for {
					fmt.Println("--->> Please Select what you want to do.")
					fmt.Println("--->> Press 1 to Add a Reservation.")
					fmt.Println("--->> Press 2 to Edit an Existing Reservation.")
					fmt.Println("--->> Press 3 to Delete a Reservation.")
					fmt.Println("--->> Press 4 to Show all Reservation.")
					fmt.Println("--->> Press 5 to Go Backward")
					var commandResponse string
					fmt.Scan(&userInput)
					if userInput == "5" {
						break
					}
					commandResponse = Reservations(userInput, conn)
					fmt.Println(commandResponse)
				}

			} else if userInput == "4" {
				sendRequest("quit", conn)
				fmt.Println("----------------------------------------------------------------------------")
				return
			}

		}
	}

}

func sendRequest(str string, conn net.Conn) {
	var err error
	_, err = conn.Write([]byte(str))
	if err != nil {
		fmt.Println("Error sending query to master server:", err)
		return
	}
}
func receiveResponse(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		if err == io.EOF {
			return "", nil // End of response
		}
		return "", err // or handle the error appropriately
	}
	response := string(buf[:n])

	return response, nil
}
func receiveResponseArray(conn net.Conn) ([]string, error) {
	var receivedArray []string
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&receivedArray)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return make([]string, 0), err
	}
	return receivedArray, nil
}

var message string
var messageArray []string

//client

type ClientForm struct {
	ClientID      string
	ClientName    string
	ClientPhone   string
	ClientCountry string
}

func Clients(userInput string, conn net.Conn) string {
	client := &ClientForm{}
	var message string
	switch userInput {
	case "1":
		message = client.insert(conn)
	case "2":
		message = client.update(conn)
	case "3":
		message = client.delete(conn)
	case "4":
		message = client.Select(conn)
	}
	return message
}

func (c *ClientForm) insert(conn net.Conn) string {
	fmt.Print("Client ID : ")
	fmt.Scan(&c.ClientID)
	fmt.Scanln()

	fmt.Print("Client Name : ")
	fmt.Scan(&c.ClientName)
	fmt.Scanln()

	fmt.Print("Phone Number : ")
	fmt.Scan(&c.ClientPhone)
	fmt.Scanln()

	fmt.Print("Country : ")
	fmt.Scan(&c.ClientCountry)
	fmt.Scanln()

	sendRequest("insertClient", conn)
	sendRequest(c.ClientID, conn)
	sendRequest(c.ClientName, conn)
	sendRequest(c.ClientPhone, conn)
	sendRequest(c.ClientCountry, conn)

	message, err := receiveResponse(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return message
}
func (c *ClientForm) update(conn net.Conn) string {
	fmt.Print("Client ID : ")
	fmt.Scan(&c.ClientID)
	fmt.Scanln()

	fmt.Print("Client Name : ")
	fmt.Scan(&c.ClientName)
	fmt.Scanln()

	fmt.Print("Phone Number : ")
	fmt.Scan(&c.ClientPhone)
	fmt.Scanln()

	fmt.Print("Country : ")
	fmt.Scan(&c.ClientCountry)
	fmt.Scanln()

	sendRequest("editClient", conn)
	sendRequest(c.ClientID, conn)
	sendRequest(c.ClientName, conn)
	sendRequest(c.ClientPhone, conn)
	sendRequest(c.ClientCountry, conn)

	message, err := receiveResponse(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return message
}
func (c *ClientForm) delete(conn net.Conn) string {
	fmt.Print("Client ID : ")
	fmt.Scan(&c.ClientID)
	fmt.Scanln()

	sendRequest("deleteClient", conn)
	sendRequest(c.ClientID, conn)
	message, err := receiveResponse(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return message
}
func (c *ClientForm) Select(conn net.Conn) string {

	fmt.Println("----------------------------------------------------------------")
	sendRequest("showClient", conn)
	messageArray, err := receiveResponseArray(conn)
	if err != nil {
		fmt.Println("Error decoding:", err)
	}
	for _, message = range messageArray {
		fmt.Print(message)
	}
	fmt.Println("----------------------------------------------------------------")

	return "Success"
}

// Room
type RoomForm struct {
	RoomId string
	IsFree string
}

func Rooms(userInput string, conn net.Conn) string {
	Room := &RoomForm{}
	var message string
	switch userInput {
	case "1":
		message = Room.insert(conn)
	case "2":
		message = Room.update(conn)
	case "3":
		message = Room.delete(conn)
	case "4":
		message = Room.Select(conn)
	}
	return message
}

func (r *RoomForm) insert(conn net.Conn) string {
	var isFree string
	fmt.Print("Room ID : ")
	fmt.Scan(&r.RoomId)
	fmt.Scanln()
	fmt.Print("Room Status : ")
	fmt.Scan(&isFree)
	fmt.Scanln()

	if isFree == "free" {
		r.IsFree = "1"
	} else {
		r.IsFree = "0"
	}

	sendRequest("insertRoom", conn)
	sendRequest(r.RoomId, conn)
	sendRequest(r.IsFree, conn)
	message, err := receiveResponse(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return message
}
func (r *RoomForm) update(conn net.Conn) string {
	var isFree string
	fmt.Print("Room ID : ")
	fmt.Scan(&r.RoomId)
	fmt.Scanln()
	fmt.Print("Room Status : ")
	fmt.Scan(&isFree)
	fmt.Scanln()

	if isFree == "free" {
		r.IsFree = "1"
	} else {
		r.IsFree = "0"
	}

	sendRequest("editRoom", conn)
	sendRequest(r.RoomId, conn)
	sendRequest(r.IsFree, conn)
	message, err := receiveResponse(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return message
}
func (r *RoomForm) delete(conn net.Conn) string {
	fmt.Print("Room ID : ")
	fmt.Scan(&r.RoomId)
	fmt.Scanln()

	sendRequest("deleteRoom", conn)
	sendRequest(r.RoomId, conn)
	message, err := receiveResponse(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return message
}
func (r *RoomForm) Select(conn net.Conn) string {
	fmt.Println("----------------------------------------------------------------")
	sendRequest("showRoom", conn)
	messageArray, err := receiveResponseArray(conn)
	if err != nil {
		fmt.Println("Error decoding:", err)
	}
	for _, message = range messageArray {
		fmt.Print(message)
	}

	fmt.Println("----------------------------------------------------------------")
	return "success"
}

// Reservations
type ReservationForm struct {
	ReservationID string
	RoomNumber    string
	ClientID      string
	CheckInDate   string
	CheckOutDate  string
}

func Reservations(userInput string, conn net.Conn) string {
	var message string
	Reservation := &ReservationForm{}

	switch userInput {
	case "1":
		message = Reservation.insert(conn)
	case "2":
		message = Reservation.update(conn)
	case "3":
		message = Reservation.delete(conn)
	case "4":
		message = Reservation.Select(conn)
	}

	return message
}

func (r *ReservationForm) insert(conn net.Conn) string {
	fmt.Print("Reservation ID : ")
	fmt.Scan(&r.ReservationID)
	fmt.Scanln()
	fmt.Print("Room ID : ")
	fmt.Scan(&r.RoomNumber)
	fmt.Scanln()
	fmt.Print("Client ID : ")
	fmt.Scan(&r.ClientID)
	fmt.Scanln()
	fmt.Print("Check-in-date : ")
	fmt.Scan(&r.CheckInDate)
	fmt.Scanln()
	fmt.Print("Check-out-date : ")
	fmt.Scan(&r.CheckOutDate)
	fmt.Scanln()

	sendRequest("insertReservation", conn)
	sendRequest(r.ReservationID, conn)
	sendRequest(r.RoomNumber, conn)
	sendRequest(r.ClientID, conn)
	sendRequest(r.CheckInDate, conn)
	sendRequest(r.CheckOutDate, conn)
	message, err := receiveResponse(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return message
}
func (r *ReservationForm) update(conn net.Conn) string {
	fmt.Print("Reservation ID : ")
	fmt.Scan(&r.ReservationID)
	fmt.Scanln()
	fmt.Print("Room ID : ")
	fmt.Scan(&r.RoomNumber)
	fmt.Scanln()
	fmt.Print("Client ID : ")
	fmt.Scan(&r.ClientID)
	fmt.Scanln()
	fmt.Print("Check-in-date : ")
	fmt.Scan(&r.CheckInDate)
	fmt.Scanln()
	fmt.Print("Check-out-date : ")
	fmt.Scan(&r.CheckOutDate)
	fmt.Scanln()

	sendRequest("editReservation", conn)
	sendRequest(r.ReservationID, conn)
	sendRequest(r.RoomNumber, conn)
	sendRequest(r.ClientID, conn)
	sendRequest(r.CheckInDate, conn)
	sendRequest(r.CheckOutDate, conn)
	message, err := receiveResponse(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return message
}
func (r *ReservationForm) delete(conn net.Conn) string {
	fmt.Print("Reservation ID : ")
	fmt.Scan(&r.ReservationID)
	fmt.Scanln()
	fmt.Print("Room ID : ")
	fmt.Scan(&r.RoomNumber)
	fmt.Scanln()

	sendRequest("deleteReservation", conn)
	sendRequest(r.ReservationID, conn)
	sendRequest(r.RoomNumber, conn)
	message, err := receiveResponse(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return message
}
func (r *ReservationForm) Select(conn net.Conn) string {
	fmt.Println("----------------------------------------------------------------")
	sendRequest("showReservation", conn)
	messageArray, err := receiveResponseArray(conn)
	if err != nil {
		fmt.Println("Error decoding:", err)
	}
	for _, message = range messageArray {
		fmt.Print(message)
	}
	fmt.Println("----------------------------------------------------------------")
	return "Success"
}
