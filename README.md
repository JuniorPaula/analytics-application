# Messaging Analytics Application

The Messaging Analytics Application is designed to calculate the average response time of operators using a specific CRM messaging platform in "`real-time`". Its main objective is to calculate the average response time for each operator in their ongoing chats and also determine the number of chats each operator is currently handling. The application is built using the **Go programming language** due to its ability to handle high data processing loads. It utilizes **MySQL** to store the metric information, **RabbitMQ** for program resilience and data integrity during calculations. The logic of the application consists of a job that runs every minute to check the chats of each operator, perform the necessary calculations, and send the data to a RabbitMQ queue. On the other end, a consumer retrieves this data and saves it to the database. Additionally, in a separate thread, another job checks if the operator has closed the chat, and if so, deletes the corresponding data from the database. It is **important** to note that, in order to run the program, besides the aforementioned dependencies, it is **crucial** to have the **access token** for the specific **CRM API** to collect chat information.

## Dependencies

- Go
- MySQL
- RabbitMQ

Make sure you have Go, MySQL, and RabbitMQ properly installed and configured in your environment before running the application.

## Running the Application

1. Clone this repository to your local machine.
2. Navigate to the project's root directory.
3. Install the necessary Go dependencies using the command `go mod download` ou `go mod tidy`.
4. Configure the MySQL database connection details in the application's configuration file.
5. Ensure RabbitMQ is running and accessible to the application.
6. Obtain the access token for the CRM API and configure it in the application.
7. Run the command `go run ./cmd/main.go` to start the application.

## Usage

The application will run as a background process, executing the necessary calculations and storing the data in the database at regular intervals. The calculated metrics can be accessed and analyzed as per your requirements.

## Contributions

Contributions are welcome! If you have any suggestions, improvements, or bug fixes, please submit a pull request. Your contributions will be greatly appreciated.

## License

This project is licensed under the terms of the [MIT License](https://opensource.org/licenses/MIT). Feel free to use and modify it according to your needs.

## Contact

If you have any questions or need further assistance, please feel free to reach out via email at [luke.junnior@icloud.com](mailto:luke.junnior@icloud.com).

Thank you for using the CRM Messaging Analytics Application!