install:
	sudo apt install -y ansible
	ansible-playbook -s -U root -i "localhost," -c local -e @/etc/raspirack/config.yml playbook.yml

mount:
	echo "Encryption password: " && read ENCRYPTION_PASSWORD && ansible-playbook -s -U root -i "localhost," -c local -e @/etc/raspirack/config.yml -e ENCRYPTION_PASSWORD=$$ENCRYPTION_PASSWORD -t mount playbook.yml

umount:
	ansible-playbook -s -U root -i "localhost," -c local -e @/etc/raspirack/config.yml -t umount playbook.yml

remote:
	ansible-playbook -u pi -i "192.168.1.2," -e @config.yml playbook.yml

sync: build
	rsync -av --exclude .git/ --exclude config.yml --exclude .history/ --exclude .vscode/ ./ pi@192.168.1.2:/home/pi/raspirack/
	ssh pi@192.168.1.2 sudo cp /home/pi/raspirack/api_arm /usr/local/bin/raspirack

build:
	go build -o api src/main/*.go
	env GOOS=linux GOARCH=arm GOARM=5 go build -o api_arm src/main/*.go

deploy: build
	scp api_arm root@raspik.local:/home/pi/raspirack/api_arm.tmp
	ssh root@raspik.local mv /home/pi/raspirack/api_arm.tmp /home/pi/raspirack/api_arm
	scp src/templates/index.html root@raspik.local:/home/pi/raspirack/src/templates/
	ssh root@raspik.local systemctl restart raspirack
