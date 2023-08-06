# Simple CRUD API :computer: :pencil:

This project provides a simple CRUD (Create, Read, Update, Delete) API for managing users. It uses MongoDB as the database backend and relies on Vault for secret management. The API exposes endpoints for interacting with user data, each one designed to perform a specific operation. :rocket:

## Table of Contents :scroll:

1. [Features](#features)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
4. [Running the Project](#running-the-project)
5. [API Endpoints](#api-endpoints)
6. [Rate Limiting](#rate-limiting)
7. [Security](#security)
8. [Contributing](#contributing)
9. [License](#license)

## Features :sparkles:

- **Create User**: Add a new user to the system. :heavy_plus_sign:
- **Read User**: Retrieve details of a single user or all users. :mag:
- **Update User**: Modify details of an existing user. :pencil2:
- **Delete User**: Remove a user from the system. :x:
- **Validation**: Validate user input before saving it to the database. :white_check_mark:
- **Rate Limiting**: Limit the number of requests from a single client. :hourglass_flowing_sand:
- **Secure Headers**: Security-enhanced HTTP headers. :lock:

## Prerequisites :memo:

To run the project, you need the following:

- MongoDB (as a data store) :card_file_box:
- Vault (for storing sensitive data like MongoDB credentials) :key:
- Go (for building and running the application) :gear:

## Installation :inbox_tray:

Clone the repository:

1. `git clone https://github.com/username/simplecrud.git`
2. `cd simplecrud`

Build the project:

1. `go build`

## Running the Project :running:

Set up your MongoDB and Vault instances, then run:

1. `./simplecrud`

The server will start at the default port 8080.

## API Endpoints :link:

Get All Users

- `GET /users`

Get User by ID

- `GET /users/:id`

Create User

- `POST /users`

Update User

- `PUT /users/:id`

Delete User

- `DELETE /users/:id`

## Rate Limiting :hourglass:

The API employs rate limiting to restrict clients to 1 request per second.

## Security :shield:

The application uses security headers to mitigate common web vulnerabilities.

## Contributing :handshake:

Feel free to fork the project and submit a pull request with your changes. Make sure to follow the code style used throughout the project.

## License :page_with_curl:

This project is licensed under the MIT License - see the LICENSE.md file for details.
