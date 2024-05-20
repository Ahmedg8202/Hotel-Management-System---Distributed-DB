// Master
package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "hotel",
		AllowNativePasswords: true,
	}
	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Database Connection Failed")
	}

	fmt.Println("Database Connected")
	db.SetMaxOpenConns(10) // Maximum number of open connections
	db.SetMaxIdleConns(10) // Maximum number of idle connections
	defer db.Close()

	listener, err := net.Listen("tcp", "localhost:9050")
	if err != nil {
		fmt.Println("Error starting master server:", err)
		return
	}

	fmt.Println("Master server listening for connections on port 9050...")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection : ", err)
			continue
		}
		fmt.Println("New connection accepted from : ", conn.RemoteAddr())
		go handleSlaveRequest(db, conn)
	}
}

func handleSlaveRequest(db *sql.DB, conn net.Conn) {
	defer conn.Close()
	for {
		req, err := getRequest(conn)
		if err != nil {
			fmt.Println(req, err)
			return
		}

		if req == "quit" {
			// Check for the special termination message
			fmt.Println("Closing connection with slave")
			return
		} else if req == "login" {
			adminName, err := getRequest(conn)
			adminPassword, err2 := getRequest(conn)
			if err != nil || err2 != nil {
				fmt.Println("Error in Get login request")
				sendResponse("Error", conn)
			}
			var isFound int
			isFound, _ = findAdmin(adminName, adminPassword)
			if isFound != 1 {
				sendResponse("Wrong Username or Password", conn)
				continue
			}
			sendResponse(strconv.Itoa(isFound), conn)
			fmt.Println("---> successful login.")
		} else if req == "insertClient" {
			var ClientID, ClientName, ClientPhone, ClientCountry, message string

			ClientID, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client ID")
				continue
			}
			ClientName, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Name")
				continue
			}
			ClientPhone, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Phone")
				continue
			}
			ClientCountry, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Country")
				continue
			}

			message, err = insertClient(ClientID, ClientName, ClientPhone, ClientCountry, db)

			if err != nil {
				fmt.Println("Error in Get customer balance")
				sendResponse("something went wrong. please try again", conn)
				continue
			} else {
				sendResponse(message, conn)
			}
			// DONE
			fmt.Println("Sucssesful Opreation")
		} else if req == "editClient" {
			var ClientID, ClientName, ClientPhone, ClientCountry, message string

			ClientID, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client ID")
				continue
			}
			ClientName, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Name")
				continue
			}
			ClientPhone, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Phone")
				continue
			}
			ClientCountry, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Country")
				continue
			}

			message, err = editClient(ClientID, ClientName, ClientPhone, ClientCountry, db)

			if err != nil {
				fmt.Println("Error in Get customer balance")
				sendResponse("something went wrong. please try again", conn)
				continue
			} else {
				sendResponse(message, conn)
			}
			// DONE
			fmt.Println("Sucssesful Opreation")
		} else if req == "deleteClient" {
			var ClientID, message string

			ClientID, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client ID")
				continue
			}

			message, err = deleteClient(ClientID, db)

			if err != nil {
				fmt.Println("Error in Get customer balance")
				sendResponse("something went wrong. please try again", conn)
				continue
			} else {
				sendResponse(message, conn)
			}
			// DONE
			fmt.Println("Sucssesful Opreation")
		} else if req == "showClient" {
			var message []string

			message, err = showClient(db)
			if err != nil {
				fmt.Println("Error in Get customer balance")
				sendResponse("something went wrong. please try again", conn)
				continue
			} else {
				sendResponsearray(message, conn)
			}

			// DONE
			fmt.Println("Successful Operation")
		} else if req == "insertRoom" {
			var RoomID, isFree, message string

			RoomID, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client ID")
				continue
			}
			isFree, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Name")
				continue
			}

			message, err = insertRoom(RoomID, isFree, db)

			if err != nil {
				fmt.Println("Error in Get customer balance")
				sendResponse("something went wrong. please try again", conn)
				continue
			} else {
				sendResponse(message, conn)
			}
			// DONE
			fmt.Println("Sucssesful Opreation")
		} else if req == "editRoom" {
			var RoomID, isFree, message string

			RoomID, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client ID")
				continue
			}
			isFree, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Name")
				continue
			}

			message, err = editRoom(RoomID, isFree, db)

			if err != nil {
				fmt.Println("Error in Get customer balance")
				sendResponse("something went wrong. please try again", conn)
				continue
			} else {
				sendResponse(message, conn)
			}
			// DONE
			fmt.Println("Sucssesful Opreation")
		} else if req == "deleteRoom" {
			var RoomID, message string

			RoomID, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client ID")
				continue
			}

			message, err = deleteRoom(RoomID, db)

			if err != nil {
				fmt.Println("Error in Get customer balance")
				sendResponse("something went wrong. please try again", conn)
				continue
			} else {
				sendResponse(message, conn)
			}
			// DONE
			fmt.Println("Sucssesful Opreation")
		} else if req == "showRoom" {
			var message []string

			message, err = showRoom(db)
			if err != nil {
				fmt.Println("Error in Get customer balance")
				sendResponse("something went wrong. please try again", conn)
				continue
			} else {
				sendResponsearray(message, conn)
			}

			// DONE
			fmt.Println("Successful Operation")
		} else if req == "insertReservation" {
			var ReserID, RoomNumber, ClientID, CheckIN, CheckOUT, message string

			ReserID, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client ID")
				continue
			}
			RoomNumber, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Name")
				continue
			}
			ClientID, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Phone")
				continue
			}
			CheckIN, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Country")
				continue
			}
			CheckOUT, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Country")
				continue
			}

			message, err = insertReservation(ReserID, RoomNumber, ClientID, CheckIN, CheckOUT, db)

			if err != nil {
				fmt.Println("Error in Get customer balance")
				sendResponse("something went wrong. please try again", conn)
				continue
			} else {
				sendResponse(message, conn)
			}
			// DONE
			fmt.Println("Sucssesful Opreation")
		} else if req == "editReservation" {
			var ReserID, RoomNumber, ClientID, CheckIN, CheckOUT, message string

			ReserID, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client ID")
				continue
			}
			RoomNumber, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Name")
				continue
			}
			ClientID, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Phone")
				continue
			}
			CheckIN, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Country")
				continue
			}
			CheckOUT, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Country")
				continue
			}

			message, err = editReservation(ReserID, RoomNumber, ClientID, CheckIN, CheckOUT, db)

			if err != nil {
				fmt.Println("Error in Get customer balance")
				sendResponse("something went wrong. please try again", conn)
				continue
			} else {
				sendResponse(message, conn)
			}
			// DONE
			fmt.Println("Sucssesful Opreation")
		} else if req == "deleteReservation" {
			var ReserID, RoomNumber, message string

			ReserID, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client ID")
				continue
			}
			RoomNumber, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Client Name")
				continue
			}

			message, err = deleteReservation(ReserID, RoomNumber, db)

			if err != nil {
				fmt.Println("Error in Get customer balance")
				sendResponse("something went wrong. please try again", conn)
				continue
			} else {
				sendResponse(message, conn)
			}
			// DONE
			fmt.Println("Sucssesful Opreation")
		} else if req == "showReservation" {
			var message []string

			message, err = showReservation(db)
			if err != nil {
				fmt.Println("Error in Get customer balance")
				sendResponse("something went wrong. please try again", conn)
				continue
			} else {
				sendResponsearray(message, conn)
			}

			// DONE
			fmt.Println("Successful Operation")
		} else {
			response := "invalid choose"
			_, err = conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Slave connection error : ", err)
			}
		}
	}
}

// /////////////////Dealing with Slaves/////////////////////
func sendResponse(msg string, conn net.Conn) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("Slave connection error : ", err)
	}
}
func sendResponsearray(msg []string, conn net.Conn) {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(msg)
	if err != nil {
		fmt.Println("Error encoding:", err)
		return
	}
}
func getRequest(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "Error reading from slave : ", err
	}
	req := string(buf[:n])
	return req, err
}

//////////////////////////////////////////SQL//////////////////

func findAdmin(adminName string, adminPassword string) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Admins WHERE AdminName = ? AND AdminPassword = ?", adminName, adminPassword).Scan(&count)
	if err != nil {
		fmt.Println("Wrong Username or Password")
		return 0, err
	}
	if count == 1 {
		fmt.Println("Login successful")
		return count, err
	} else {
		fmt.Println("Wrong Username or Password")
		return 0, err
	}
}
func insertClient(clientId string, clientName string, clientPhone string, clientCountry string, db *sql.DB) (string, error) {
	_, err := db.Exec("INSERT INTO Clients (ClientID, Name, Phone, Country) VALUES (?, ?, ?, ?)", clientId, clientName, clientPhone, clientCountry)
	if err != nil {
		fmt.Println(err)
		return "Add Failed", err
	}
	return "Client Successfully Added", nil
}
func editClient(clientId string, clientName string, clientPhone string, clientCountry string, db *sql.DB) (string, error) {
	_, err := db.Exec("UPDATE Clients SET Name = ?, Phone = ?, Country = ? WHERE ClientID = ?", clientName, clientPhone, clientCountry, clientId)
	if err != nil {
		fmt.Println(err)
		return "Edit Failed", err
	}
	return "Client Successfully Edited", nil
}
func deleteClient(clientId string, db *sql.DB) (string, error) {
	_, err := db.Exec("DELETE FROM Clients WHERE ClientID = ?", clientId)
	if err != nil {
		fmt.Println(err)
		return "Delete Failed", err
	}
	return "Client Successfully Deleted", nil
}
func showClient(db *sql.DB) ([]string, error) {
	var clientID, clientName, clientPhone, clientCountry string
	rows, err := db.Query("SELECT * FROM Clients")
	if err != nil {
		return make([]string, 0), err
	}
	defer rows.Close()

	var clients []string

	for rows.Next() {

		if err := rows.Scan(&clientID, &clientName, &clientPhone, &clientCountry); err != nil {
			fmt.Println(err)
			continue
		}

		clientStr := fmt.Sprintf("ID: %s, Name: %s, Phone: %s, Country: %s\n", clientID, clientName, clientPhone, clientCountry)
		clients = append(clients, clientStr)

	}

	if err := rows.Err(); err != nil {
		return make([]string, 0), err
	}
	return clients, nil
}
func insertRoom(roomId string, isfree string, db *sql.DB) (string, error) {
	_, err := db.Exec("INSERT INTO Rooms (RoomNumber, IsAvailable) VALUES (?, ?)", roomId, isfree)
	if err != nil {
		return "Added Failed", err
	}
	return "Room Successfully Added", nil
}
func editRoom(roomId string, isfree string, db *sql.DB) (string, error) {
	_, err := db.Exec("UPDATE Rooms SET IsAvailable = ? WHERE RoomNumber = ?", isfree, roomId)
	if err != nil {
		return "Edit Failed", err
	}
	return "Room Successfully Edited", nil
}
func deleteRoom(roomId string, db *sql.DB) (string, error) {
	_, err := db.Exec("DELETE FROM Rooms WHERE RoomNumber = ?", roomId)
	if err != nil {
		return "delete Failed", err
	}
	return "Room Successfully deleted", nil
}
func showRoom(db *sql.DB) ([]string, error) {
	var RoomID, isFree string
	rows, err := db.Query("SELECT * FROM Rooms")
	if err != nil {
		fmt.Println(err)
		return make([]string, 0), err
	}
	defer rows.Close()

	var rooms []string

	for rows.Next() {

		if err := rows.Scan(&RoomID, &isFree); err != nil {
			fmt.Println(err)
			continue
		}
		var free string
		if isFree == "1" {
			free = "free"
		} else {
			free = "not free"
		}

		roomStr := fmt.Sprintf("ID: %s, Free: %s\n", RoomID, free)
		rooms = append(rooms, roomStr)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)

		return make([]string, 0), err
	}

	return rooms, nil
}
func insertReservation(reserId string, roomNum string, clientId string, checkIn string, checkOut string, db *sql.DB) (string, error) {
	_, err := db.Exec("INSERT INTO Reservations (ReservationID, RoomNumber, ClientID, CheckInDate, CheckOutDate) VALUES (?, ?, ?, ?, ?)", reserId, roomNum, clientId, checkIn, checkOut)
	if err != nil {
		fmt.Println(err)
		return "Add Failed", err
	}
	return "Reservation Successfully Added", nil
}
func editReservation(reserId string, roomNum string, clientId string, checkIn string, checkOut string, db *sql.DB) (string, error) {
	_, err := db.Exec("UPDATE Reservations SET RoomNumber = ?, ClientID = ?, CheckInDate = ?, CheckOutDate = ? WHERE ReservationID = ?", roomNum, clientId, checkIn, checkOut, reserId)
	if err != nil {
		fmt.Println(err)
		return "Edit Failed", err
	}
	return "Reservation Successfully Edited", nil
}
func deleteReservation(reserId string, roomNum string, db *sql.DB) (string, error) {
	_, err := db.Exec("DELETE FROM Reservations WHERE ReservationID = ? AND RoomNumber = ?", reserId, roomNum)
	if err != nil {
		fmt.Println(err)
		return "Edit Failed", err
	}
	return "Reservation Successfully Edited", nil
}
func showReservation(db *sql.DB) ([]string, error) {
	var ReservationID, RoomNumber, ClientID, CheckInDate, CheckOutDate string
	rows, err := db.Query("SELECT * FROM Reservations")
	if err != nil {
		return make([]string, 0), err
	}
	defer rows.Close()

	var reservations []string

	for rows.Next() {

		if err := rows.Scan(&ReservationID, &RoomNumber, &ClientID, &CheckInDate, &CheckOutDate); err != nil {
			fmt.Println(err)
			continue
		}

		reservationStr := fmt.Sprintf("Reservation ID: %s, Room ID: %s, Client ID: %s,Check-IN: %s,Check-OUT: %s\n", ReservationID, RoomNumber, ClientID, CheckInDate, CheckOutDate)
		reservations = append(reservations, reservationStr)
	}

	if err := rows.Err(); err != nil {
		return make([]string, 0), err
	}
	return reservations, nil
}
