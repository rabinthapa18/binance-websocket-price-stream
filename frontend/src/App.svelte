<script lang="ts">
  import { onMount, onDestroy } from "svelte";

  let connectionStatus: string = "Connecting...";
  let ws: WebSocket | null = null;
  let reconnectInterval: number | null = null;
  let prices: { price: number; time: string }[] = []; // array to store prices with time

  // function to create WebSocket connection
  function connect() {
    ws = new WebSocket("ws://localhost:8080/ws");

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      addPrice(data); // adding price to prices array
    };

    ws.onerror = (error) => {
      console.error("WebSocket Error:", error);
      connectionStatus = "Connection error"; // updating connection status
    };

    ws.onclose = () => {
      connectionStatus = "Connection Lost, reconnecting..."; // updating connection status
      attemptReconnect(); // attempting reconnection
    };

    connectionStatus = "Connected"; // setting initial connection status
  }

  // function to attempt reconnection
  function attemptReconnect() {
    if (reconnectInterval) return; // return if already reconnecting

    reconnectInterval = setInterval(() => {
      console.log("Attempting to reconnect...");
      connect(); // reconnecting
    }, 5000); // reconnecting every 5 seconds
  }

  // function to add price to prices array
  function addPrice(price: number) {
    if (prices.length >= 10) {
      prices.pop(); // remove oldest price
    }
    const time = new Date().toLocaleTimeString();
    prices = [{ price, time }, ...prices]; // add new price with time
  }

  onMount(() => {
    connect(); // creating initial connection
  });

  onDestroy(() => {
    if (ws) {
      ws.close(); // close WebSocket connection on unmount
    }

    if (reconnectInterval) {
      clearInterval(reconnectInterval); // clear reconnect interval
    }
  });
</script>

<main>
  <h1>Average Prices</h1>
  <p>Connection status: {connectionStatus}</p>

  {#if prices.length > 0}
    <!-- table to show prices with time -->
    <table>
      <thead>
        <tr>
          <th>Price</th>
          <th>Time</th>
        </tr>
      </thead>
      <tbody>
        {#each prices as { price, time }}
          <tr>
            <td>{price.toFixed(2)}</td>
            <td>{time}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <p>Loading...</p>
  {/if}
</main>

<style>
  main {
    text-align: center;
    padding: 1em;
    margin: 0 auto;
  }

  h1 {
    color: #c6ac8f;
  }

  p {
    font-size: 1.2em;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 1em;
  }

  th,
  td {
    padding: 8px;
    border: 1px solid #c6ac8f;
    text-align: center;
  }

  th {
    background-color: #5e503f;
  }
</style>
