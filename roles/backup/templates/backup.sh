#!/usr/bin/fish

set PIDFILE /run/backup.pid

if test ! -e $PIDFILE
    echo 999999 > $PIDFILE
end

if ps -A | sed "s/^[ ]*//g" | cut -d" " -f 1 | grep -e \^(cat $PIDFILE)\$ > /dev/null
    echo "Backup process is already running"
    exit 5
end

echo %self > $PIDFILE

{% for key, value in backup.restic_env.items() %}
set -x {{ key }} "{{ value }}"
{% endfor %}

# Initialization of the repository, don't worry about fatal error if the repo is created :-)
/usr/local/bin/restic init

# Backup all drives
{% for drive in drives %}
/usr/local/bin/restic backup /mnt/{{ drive.name }}
{% endfor %}

# Clear old data
/usr/local/bin/restic forget --prune --keep-daily {{ backup.keep_daily }} --keep-weekly {{ backup.keep_weekly }} --keep-monthly {{ backup.keep_monthly }}
