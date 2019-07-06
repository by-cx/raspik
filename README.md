# RaspiRack

Goal of this project is to deliver home made NAS from scratch. The hardware part is based on Raspberry Pi 4 and 3D printed case very similar to a server's racks. The software part is based on Ansible and vanilla Raspbian. It's still in alpha stage of development and I am quite sure it remains there for some time.

I am happy user of Synology NAS and my DS216play has been in use for a few year now. It's pretty neat piece of hardware but it's not perfect. Encryption is slow and can't be used for photos managed by Synology's PhotoStation. Also the interface became slow after the years and a lot of updates that I started using only fraction of its capabilities. They are Samba shares, backup and SynologyDrive. Time of replacing this device is coming so I decided to explore possibilities of building it myself. Raspberry Pi 4 is out. Its encryption capabilities are not the best but also not the worst and it supports USB 3 devices. The only missing piece is the case and software that would make managing my data easy again.

## The Rack

I designed modular case for Raspberry Pi that is not just for one purpose but it could help you with your other projects too. It has 8 slots where you can put different kind of hardware like:

* Power supply
* Raspberry Pi
* 2.5" HDD/SSD
* Switch
* Your personal devices

The main idea behind this is I wanted to have a place where I can put two Raspberry Pis and two 2.5" SSD. In my case it's all feeded with 14 A 5 V power supply and the network is managed by 5 port gigabit switch. Everything is enclosed and it's easy to attach a fan to cool down the hot new Raspberry Pi 4. Every device in this box has it's own blade or a shelf if you want. Every blade just slides in and it can be removed any time. That means it's not hard to change order of your devices in RaspiRack.

There is still a lot of work in the hardware part, especially in power supply and switch part, because they take a lot of space but if you want it's good enough to use it. My Raspberry is sitting inside right now without any issue.

## The main features

This project is designed by my needs but feel free to send an idea or pull request. What ever you done to help this project, please keep in mind it's focused mainly to hobbyist and home users. The main features I don't want to break are those:

### Security

All data you send to RaspiRack should be encrypted. Your files like photos, documents or simple computer backups shouldn't be available publicly if somebody steals your Raspberry Pi. The project can't protect you from everything. It's still possible to run different kinds of attacks to your data but holding the hard drive physically shouldn't guarantee access to the data inside it.

### Backup

Implementing RAID with Raspberry Pi for could be overhead. Hard drives are not failing every day so this project focuses on doing proper backups instead. It uses [Restic](https://restic.net/) that supports many backend and allows doing backups every hour without significant load for your internet connection. I am honestly telling you that this project won't protect you from losing recently uploaded files just before hardware failure happens, but it makes sure you won't loose all your data in that case.

### Simplicity

The structure of shares is split into two zones. The first one is private space for all users and shared space for those users. That means every user has his own space where he can upload his private data no body else can see. The shared space is every other Samba share that is created additionally in the configuration. Those shares are available for all users.
