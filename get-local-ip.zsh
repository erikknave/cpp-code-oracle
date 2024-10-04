#!/bin/zsh

# Check if Wi-Fi is connected and get the IP address
WIFI_IP=$(ipconfig getifaddr en0)

# Check if Ethernet is connected and get the IP address
ETH_IP=$(ipconfig getifaddr en1)

# Determine which interface is active and display the IP address
if [ -n "$WIFI_IP" ]; then
    echo "Your Wi-Fi IP address is: $WIFI_IP"
elif [ -n "$ETH_IP" ]; then
    echo "Your Ethernet IP address is: $ETH_IP"
else
    echo "No active network connection found."
fi
