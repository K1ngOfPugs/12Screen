#!/bin/bash

BACKUP_FOLDER="ID_BACKUP_($(date '+%Y-%m-%d_%H-%M-%S'))"

cd /var/mobile
mkdir $BACKUP_FOLDER
cd $BACKUP_FOLDER
mv -p /private/var/mobile/Library/Accounts/ Accounts/
mv -p /private/var/mobile/Library/Prefrences/com.apple.accounts.plist com.apple.accounts.plist
mkdir /private/var/mobile/Library/Accounts/

# Soft reboot to remove Apple ID #
ldrestart