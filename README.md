# Hotshot Query Example

This is a simple Go project that fetches and prints namespace transactions from the Hotshot query service. The application periodically polls for the latest block height, then retrieves and pretty-prints the namespace transactions in your terminal. It skips fallback responses (e.g. those containing `"FetchBlock"`).

## Overview

- **Block Height Endpoint:**  
  Retrieves the latest block height from `/status/block-height`.

- **Namespace Transactions:**  
  Uses the configured HotShot URL along with the block height and namespace to build the full endpoint:
  ```
  <hotshot_url>/availability/block/<blockHeight>/namespace/<namespace>
  ```
  The returned JSON is filtered and pretty-printed to display only the relevant transaction data. Fallback responses (such as those containing `"FetchBlock"`) are omitted.

## Requirements

- Go 1.23.2 or higher

## Getting Started

1. **Clone the Repository**

   ```bash
   git clone <repository-url>
   cd hackathon-example
   ```

2. **Run the Application**

   ```bash
   go run main.go
   ```

## Configuration

The application uses a configuration file located at `config/config.json`. This file allows you to customize the following parameters:

- **Namespace:** Change the namespace identifier.
- **Query Service Endpoint:** Set your Hotshot query service URLs.
- **Namespace (Chain ID):** Change the namespace identifier. **Note:** In this project, the namespace is the same as the chain id used by the network. This means that when you configure the namespace, you are effectively specifying the chain identifier.
- **Polling Interval:** Adjust the time interval for how frequently the project polls the API.

Below is an example configuration file:

```json
{
  "hotshot_url": "https://query.decaf.testnet.espresso.network/v0",
  "namespace": 20250115,
  "polling_interval": 10
}
```

## Example Logs

Example logs are available in the `example-logs` directory. A typical log entry will look like:

```json
Hotshot Namespace Transactions:
{
  "proof": null,
  "transactions": []
}
```

## Extending the Project

You can easily extend this project to interact with additional endpoints (such as block headers or explorer transactions) by following the existing structure. Simply update the configuration and add new functions to handle the additional API calls.

## License

This project is provided as-is with no warranty.