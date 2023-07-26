# Chatbot StockAPI 🤖

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/ThaisGLeite/chatbot-stockapi/issues)

Chatbot StockAPI is a **highly interactive real-time application** designed to provide users with **live updates on stock market data**. Leveraging an integration with the [Stooq](https://stooq.com) API for its stock market information, this chatbot ensures users are provided with the most up-to-date financial data.

![Chatbot Stock API Architecture](./images/chatbot-stockapi-architecture.png)

The application leverages **NATS** as a messaging system to enable seamless communication between components and **Redis** for efficient data caching and storage.

## 📚 Table of Contents

1. [About The Project](#about-the-project)
2. [Project Structure](#project-structure)
3. [Getting Started](#getting-started)
4. [Using the `stop-start.sh` Script](#using-the-stop-startsh-script)
5. [Services](#services)
6. [Contribution](#contribution)
7. [License](#license)

## 💡 About The Project

The primary objective of Chatbot StockAPI is to enable users to query for live stock market data through a convenient chat interface. The application is structured around two main components:

- The **chatroom service**, which manages user interactions, including registration, login, and chatroom activities.
- The **bot service**, which interfaces with the external Stock API to fetch live stock market data and communicates this information to the chatroom service via NATS messaging.

## 📖 Project Structure

The application is structured into several key directories:

- `botService`: Contains the logic for the bot service that communicates with the Stock API.
- `cmd`: Includes the main application entry points and utility scripts like `stop-start.sh`.
- `Docker`: Contains Dockerfiles for the different services and docker-compose for starting the services.
- `handle`: Includes the handlers for the different routes in the application.
- `model`: Contains the data structures used across the application.
- `natsclient`: Code related to NATS client implementation.
- `redis`: Code related to Redis client implementation.
- `static`: Holds static files like HTML, CSS, and JavaScript for the frontend of the chatbot.
- `ws`: Includes the implementation of the Websocket server used for real-time communication.

## 🚀 Getting Started

To get the application up and running, follow these steps:

1. **Clone the repository**: `git clone https://github.com/ThaisGLeite/chatbot-stockapi.git`
2. **Navigate into the project directory**: `cd chatbot-stockapi`
3. **Run the `stop-start.sh` script**: `./cmd/stop-start.sh`

The Chatbot application will now be accessible at `http://localhost:8080`, and the BotService at `http://localhost:3000`.

## 🖥️ Using the `stop-start.sh` Script

The `stop-start.sh` script is a utility script provided in this repository to help manage the Docker containers for this project. Here is a brief rundown of its functionality:

1. It navigates to the directory containing the `docker-compose.yml` file.
2. It stops any running containers associated with the `docker-compose.yml` file using `docker-compose down`. This is useful if you have previously started containers and want to ensure a clean slate.
3. It builds the services and recreates the containers from scratch using `docker-compose up --build --force-recreate`. This ensures that any changes you have made to your Docker files will be included when the containers are restarted.

You can run this script from the `cmd` directory with the command `./stop-start.sh`. Ensure the script has the appropriate permissions to be executed (`chmod +x stop-start.sh` if needed).

## 🔧 Services

1. **BotService**: This is a Go service that communicates with the Stock API. It subscribes to a NATS subject "stock_codes", and upon receiving messages, it fetches the corresponding stock data and publishes it back to NATS.
2. **Chatroom**: This is the main chat application. It handles user registration and login, JWT token generation and validation, creation of chatrooms, and sending/receiving messages. It listens to stock updates published by the BotService via NATS and pushes them to connected clients over WebSockets.
3. **NATS**: An open-source messaging system used for communication between the BotService and Chatroom.
4. **Redis**: An open-source, in-memory data structure store used as a cache and a message broker. In this application, it stores user data and session information, as well as chatroom details and messages.

## 🤝 Contribution

We welcome contributions from the community. Feel free to fork the project, make changes, and submit a pull request. We ask that you respect the existing coding style and commit message format for consistency.

## 📜 License

This project is licensed under the terms of the MIT license.
