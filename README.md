# Binance WebSocket Average Price Project

## Overview

This project demonstrates how to calculate the average price of the Binance BTC/USD order book using a WebSocket connection and broadcast the result to connected clients. The main components are: the **Backend (Golang)** that handles Binance WebSocket communication and calculation, and the **Frontend (Svelte)** that displays the received prices to the users in real-time.

## Backend

The backend is built using **Golang**. It connects to Binance's WebSocket to retrieve realtime order book data for BTC/USD, calculates the average price of the bids and asks, and broadcasts this average price to all connected WebSocket clients.

### Running the Backend

To run the backend server:

1. Navigate to the `backend` directory.
2. Install the dependencies:

```
go mod download
```

3. Run the following command:

```
go run main.go
```

The backend server will start listening on `ws://localhost:8080/ws` for WebSocket connections from the frontend.

### Running the Tests

To run the unit tests:

```
go test
```

## Frontend

The frontend is built using **Svelte** with TypeScript, and it connects to the backend WebSocket server to display real-time average prices on the web page.

### Running the Frontend

1. Navigate to the `frontend` directory.
2. Install the dependencies:

```
npm install
```

3. Run the development server:

```
npm run dev
```

The frontend server will be available at `http://localhost:5173`.
