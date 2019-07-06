#!/usr/bin/fish

{% for key, value in backup.restic_env %}
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
