#!/bin/bash
export GPG_TTY=$(tty)              # prompt for password

read -r -p "Have you already setup a gpgkey that you would like to use for a pass store?(Y/n)" has_gpg
if [ "$has_gpg" = "y" ] || [ "$has_gpg" = "Y" ]; then
	:
else
	gpg --gen-key
	gpg-connect-agent reloadagent /bye # restart gpg agent
fi

echo "Initializing password store utility"
pass init onomy
