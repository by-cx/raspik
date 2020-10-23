#!/usr/bin/fish

{% for key, value in backup.restic_env.items() %}
set -x {{ key }} "{{ value }}"
{% endfor %}
