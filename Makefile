install:
	sudo apt install -y ansible
	ansible-playbook -s -U root -i "localhost," -c local -e @/etc/raspirack/config.yml playbook.yml

mount:
	ansible-playbook -s -U root -i "localhost," -c local -e @/etc/raspirack/config.yml -t mount playbook.yml

umount:
	ansible-playbook -s -U root -i "localhost," -c local -e @/etc/raspirack/config.yml -t umount playbook.yml

remote:
	ansible-playbook -u pi -i "192.168.1.2," -e @config.yml playbook.yml
