install:
	sudo apt install -y ansible
	ansible-playbook -s -U root -i "localhost," -c local -e @/etc/raspirack/config.yml playbook.yml

mount:
	echo "Encryption password: " && read ENCRYPTION_PASSWORD && ansible-playbook -s -U root -i "localhost," -c local -e @/etc/raspirack/config.yml -e ENCRYPTION_PASSWORD=$$ENCRYPTION_PASSWORD -t mount playbook.yml

umount:
	ansible-playbook -s -U root -i "localhost," -c local -e @/etc/raspirack/config.yml -t umount playbook.yml

remote:
	ansible-playbook -u pi -i "192.168.1.2," -e @config.yml playbook.yml

sync:
	rsync -av --exclude .git/ --exclude config.yml --exclude .history/ --exclude .vscode/ ./ pi@192.168.1.2:/home/pi/raspirack/

build:
	go build -o api src/main/*.go