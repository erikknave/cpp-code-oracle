#!/bin/zsh

# Unload the ZeroTier service
echo "Unloading ZeroTier service..."
sudo launchctl unload /Library/LaunchDaemons/com.zerotier.one.plist

# Check if the unload was successful
if [[ $? -eq 0 ]]; then
    echo "ZeroTier service unloaded successfully."
else
    echo "Failed to unload ZeroTier service. Exiting."
    exit 1
fi

# Load the ZeroTier service
echo "Loading ZeroTier service..."
sudo launchctl load /Library/LaunchDaemons/com.zerotier.one.plist

# Check if the load was successful
if [[ $? -eq 0 ]]; then
    echo "ZeroTier service loaded successfully."
else
    echo "Failed to load ZeroTier service. Exiting."
    exit 1
fi

echo "ZeroTier service has been restarted successfully."
