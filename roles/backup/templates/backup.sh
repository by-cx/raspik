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

# Backup homes
{% if general.homes_backup %}
/usr/local/bin/restic backup /mnt/{{ drives[general.homes_drive].name }}/homes/
{% endif %}

# Backup all drives
{% for share in shares %}
{% if share.backup %}/usr/local/bin/restic backup /mnt/{{ drives[share.drive].name }}/{{ share.name }}{% endif %}
{% endfor %}

# Clear old data
/usr/local/bin/restic forget --prune --keep-daily {{ backup.keep_daily }} --keep-weekly {{ backup.keep_weekly }} --keep-monthly {{ backup.keep_monthly }}
