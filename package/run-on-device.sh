#!/bin/bash

extractFieldsAndExecute() {
    # go to any website so I can get the metadata (payphoneMAC, payphoneID, payphoneTime)
    output=$(curl --interface "$interface" http://google.com -m 10) || { echo "$(date) [ERROR] inital curl timeout"; return 1; }
    urlAddress=$(echo "$(date) $output" | cut -d\" -f2)
    # need to use -L to get cookies (cookies are used to authenticate)
    curl -c /tmp/cookies.txt -s -L --interface $interface http://google.com -m 10 || { echo "$(date) [ERROR] initial curl timeout"; return 1; }
    payphoneMAC=$(echo $urlAddress | grep -oP '\?mac=\K[^&]*')
    # a little trick to urldecode. Replace %3A into :
    payphoneMAC=$(echo "${payphoneMAC//%3A/:}")
    payphoneID=$(echo $urlAddress | grep -oP 'a=\K[^&]*')
    payphoneTime=$(echo $urlAddress | grep -oP 'b=\K[^&]*')

    echo "$(date) [INFO] urlAddress is $urlAddress"
    if [ -z "$payphoneMAC" ] || [ -z "$payphoneID" ] || [ -z "$payphoneTime" ]; then
        echo "$(date) [INFO] $payphoneMAC"
        echo "$(date) [INFO] $payphoneID"
        echo "$(date) [INFO] $payphoneTime"
        echo "$(date) [ERROR] One or more payphone details are missing or empty."
        return 1
    fi
    echo "$(date) [INFO] Succesfully got urlAddress"
    # this will connect us to the internet (using the cookies for authentication)
    curl --interface "$interface" -m 20 -H "X-Requested-With: XMLHTTPRequest" -b /tmp/cookies.txt \
    https://apac.network-auth.com/splash/NAxIVbNc.5.167/grant?continue_url= || { echo "$(date) [ERROR] curl timeout"; return 1; }

    # check if connected to internet
    echo "$(date) [INFO] Successfully connected to internet"
    echo "$(date) [INFO] $payphoneMAC"
    echo "$(date) [INFO] $payphoneID"
    echo "$(date) [INFO] $payphoneTime"
    # run executable
    /usr/local/bin/client-indicum $payphoneMAC $payphoneID $payphoneTime

    echo "$(date) [INFO] Successfully ran client-indicum"
    # change our MAC address (so we will have to sign in again when we re-see payphone)
    sudo ifconfig "$interface" down && sudo macchanger -r "$interface" && sudo ifconfig "$interface" up

    echo "$(date) [INFO] successfully changed mac address";
}

while [ true ]; do
    interface="wlan0"
    expectedSSID="Free Telstra Wi-Fi"
    currentSSID=$(iwgetid -r "$interface")
    # initial internet test check
    if ping -c 1 8.8.8.8 -W 5 -I "$interface" &> /dev/null; then
        echo "$(date) [INFO] ping success"
	if [ "$currentSSID" = "$expectedSSID" ]; then
        	extractFieldsAndExecute
	fi
    else
        # if failed, check if interface is connected to wifi but just doesn't have internet access
        if [ "$currentSSID" = "$expectedSSID" ]; then
            # if so, reset MAC address, this usually fixes it
            sudo ifconfig "$interface" down && sudo macchanger -r "$interface" && sudo ifconfig "$interface" up
        fi
        echo "$(date) [INFO] currentSSID is $currentSSID"
    fi

    sleep 5

done