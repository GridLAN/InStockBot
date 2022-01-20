# InStockBot Template
InStockBot notifies a selected Discord channel when a specific product is back in stock.

This is primarily for me to purchase more Ubiquiti networking equipment. 

Get InStockBot in your discord [here](https://discord.com/api/oauth2/authorize?client_id=933558277571244082&permissions=534723951680&scope=bot%20applications.commands).

## Development:
1. Compile and run the project.

    ```
    TOKEN=abc123 go run main.go
    ```

2. Alternatively, build and run the project inside of a container.

    ```
    docker build -t instockbot . && docker run -d --env TOKEN='abc123' instockbot
    ```