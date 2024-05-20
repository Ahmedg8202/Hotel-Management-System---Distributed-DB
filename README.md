# Hotel Management System

## Overview

The Hotel Management System is a client-server application designed to facilitate the management of hotel operations, including client management, room management, and reservations. The system consists of a master server and multiple slave clients. The master server handles the database operations and communication with the slave clients, while the slave clients provide a user interface for hotel administrators to interact with the system.

## Features

- **Client Management**: Add, edit, delete, and display clients.
- **Room Management**: Add, edit, delete, and display rooms.
- **Reservation Management**: Add, edit, delete, and display reservations.
- **User Authentication**: Admin login for secure access to the system.

## Technologies Used

- Go (Golang) for server and client-side programming.
- MySQL for database management.
- TCP for client-server communication.

## Architecture

- **Master Server**: Manages the database and handles requests from slave clients.
- **Slave Clients**: Provides a command-line interface for admins to manage hotel operations and communicates with the master server.

## Setup Instructions

### Prerequisites

- Go (1.16 or later)
- MySQL database

### Database Setup

1. Install MySQL and create a database named `hotel`.
2. Create the necessary tables:

```sql
CREATE TABLE admin (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE clients (
    client_id VARCHAR(255) PRIMARY KEY,
    client_name VARCHAR(255) NOT NULL,
    client_phone VARCHAR(255) NOT NULL,
    client_country VARCHAR(255) NOT NULL
);

CREATE TABLE rooms (
    room_id VARCHAR(255) PRIMARY KEY,
    is_free BOOLEAN NOT NULL
);

CREATE TABLE reservations (
    reservation_id VARCHAR(255) PRIMARY KEY,
    room_id VARCHAR(255) NOT NULL,
    client_id VARCHAR(255) NOT NULL,
    check_in_date DATE NOT NULL,
    check_out_date DATE NOT NULL,
    FOREIGN KEY (room_id) REFERENCES rooms(room_id),
    FOREIGN KEY (client_id) REFERENCES clients(client_id)
);
```

### Configuration

1. Update the MySQL connection details in the `main.go` file of the master server:

```go
cfg := mysql.Config{
    User:                 "root",
    Passwd:               "root",
    Net:                  "tcp",
    Addr:                 "127.0.0.1:3306",
    DBName:               "hotel",
    AllowNativePasswords: true,
}
```

### Running the Application

1. **Start the Master Server**:
    ```sh
    cd master
    go run main.go
    ```

2. **Start a Slave Client**:
    ```sh
    cd slave
    go run main.go
    ```

## Usage

1. **Login**:
    - On the slave client, enter your admin username and password to log in.
2. **Client Management**:
    - Press `1` to manage clients.
    - Options:
        - Press `1` to add a client.
        - Press `2` to edit an existing client.
        - Press `3` to delete a client.
        - Press `4` to show all clients.
        - Press `5` to go back to the main menu.
3. **Room Management**:
    - Press `2` to manage rooms.
    - Options:
        - Press `1` to add a room.
        - Press `2` to edit an existing room.
        - Press `3` to delete a room.
        - Press `4` to show all rooms.
        - Press `5` to go back to the main menu.
4. **Reservation Management**:
    - Press `3` to manage reservations.
    - Options:
        - Press `1` to add a reservation.
        - Press `2` to edit an existing reservation.
        - Press `3` to delete a reservation.
        - Press `4` to show all reservations.
        - Press `5` to go back to the main menu.
5. **Logout**:
    - Press `4` to log out and exit the application.

## Contribution

Contributions are welcome! Please submit a pull request or open an issue to discuss any changes or improvements.

## Contact

For any questions or support, please contact the project maintainer:

- Name: [Ahmed Gamal]
- Email: [agamal8202@gmail.com]
