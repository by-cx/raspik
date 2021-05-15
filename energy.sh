#!/usr/bin/fish

set csrf_token (http --session=logged_energomonitor https://app.energomonitor.cz/login | grep csrf | grep -o -e "value='.*'" | sed "s/value=//" | sed "s/'//g")
echo http --session=logged_energomonitor POST https://app.energomonitor.cz/login username=cx password=lFqcvwRYe8uN csrfmiddlewaretoken=$csrf_token
http --session=logged_energomonitor POST https://app.energomonitor.cz/login username=cx password=lFqcvwRYe8uN csrfmiddlewaretoken=$csrf_token
