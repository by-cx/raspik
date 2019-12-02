# RaspiRack

Goal of this project is to deliver home made NAS from scratch. The hardware part is based on Raspberry Pi 4 and 3D printed case very similar to server's racks. The software part is based on Ansible and vanilla Raspbian. It's still in alpha stage of development and I am quite sure it remains there for some time.

I am happy user of Synology NAS and my DS216play has been in use for a few year now. It's pretty neat piece of hardware but it's not perfect. Encryption is slow and can't be used for photos managed by Synology's PhotoStation. Also the interface became slow after the years and a lot of updates that I started using only fraction of its capabilities. They are Samba shares, backup and SynologyDrive. Time of replacing this device is coming so I decided to explore possibilities of building it myself. Raspberry Pi 4 is out. Its encryption capabilities are not the best but also not the worst and it supports USB 3 devices. The only missing piece is the case and software that would make managing my data easy again.

You don't need the case to use the software part and vice versa.

## Quick tutorial



## The Rack

I designed modular case for Raspberry Pi that is not just for one purpose but it could help you with your other projects too. It has 8 slots where you can put different kind of hardware like:

* Power supply
* Raspberry Pi
* 2.5" HDD/SSD
* Switch
* Your personal devices

The main idea behind this is I wanted to have a place where I can put two Raspberry Pis and two 2.5" SSD. In my case it's all feeded with 14 A 5 V power supply and the network is managed by 5 port gigabit switch. Everything is enclosed and it's easy to attach a fan to cool down the hot new Raspberry Pi 4. Every device in this box has it's own blade or a shelf if you want. Every blade just slides in and it can be removed any time. That means it's not hard to change order of your devices in RaspiRack.

There is still a lot of work in the hardware part, especially in power supply and switch part, because they take a lot of space but if you want it's good enough to use it. My Raspberry is sitting inside right now without any issue.

Dimensions of the case:

* Whole case: 206 × 128 × 173 mm
* Blade: 1.6 × 110 × 149 mm (space for your equipment is 22 × 110 × 149 mm)
* Blade with sliders: 1.6 × 115 × 149 mm 

Cooling is handled by a fan 92×92 mm. It can be attached at the back of front of the case. Except this front/back part there is one big and one small part. You can use them to cover front/back panels exactly like you want. The 92×92 mm fans are usually 12 V but they work perfectly on 5 V too. Raspberry Pi doesn't need anything extra. Slowly moving air is enough even in during the summer.

Needed screws:

* 7× 3x8mm+ (power source top cover and EU socket)
* 2× 3x6 mm (power supply mount)
* ~20× 3x16mm+ (front and back panel)

The case is not picky. Except the two power supply screws every other screw is minimal lenght.

## Blades

I use multiple Raspberry Pis so having more than one of them in one case is added value I couldn't miss while I was working on the case. RPis are pretty sensitive about power supply and I needed two 2.5" SSD and 433 MHz antenna inside. Handle all of it inside one box made smile on my face :-) Let's check each blade. Standard server racks were great inspiration for me so the this case looks like one of them.

### Power supply blade

I use switching power supply [MEAN WELL LRS-75-5](http://www.mean-well.cz/assets/data/LRS-75-spec.pdf). It's sold here in Czech Republic. It delivers max. 14 A in 5 V. The maximal current is much bigger than I need or I will need. On the other hand if you decide to fill all 8 slots with RPis this power source can handle it. If I am choosing the power supply now I would pick something with 5 V and 12 V output. That allows to supply 3.5" HDD. Unfortunately 3.5" HDD barely fits inside the rack.

** IMPORTANT: Please print both parts of the power supply blade including the top cover. You will be touching the rack and you don't want to touch the live wire. **

This blade requires two slots inside the case.

** IMPORTANT: I experienced instabilities in Raspberry Pi when filtering capacitor wasn't connected at the output of the power supply. Add something like 330uF capacitor at the output if you want to avoid that. **

If you decide, like me, to supply one or more USB devices directly from the power supply and power line into Raspberry Pi won't be disconnected. It's a good idea to add a poly-fuse, cut the 5V line before the USB connector or both,

Also output of the power supply should be covered with 10 or 15 A fuse. I use 10 A. The selected power supply has integrated protection mechanisms but the datasheet doesn't mention short circuit protection specifically. Only Hiccup mode I have no experience with.

### Switch blade

I would say this blade is optional. You can depend on WiFi or having the switch somewhere else. I use Tenda TEG1005D. It's super small, 1 Gbit switch that fits inside perfectly. The switch is glued to the blade.

As the power supply blade this one also requires two slots. I am planning to squeeze them together so they will need only three slots in total but I have plenty of space available for now.

### 2.5" HDD/SSD blade

This blade is designed around [Icy Box IB-AC703-U3](https://www.amazon.de/IB-AC703-U3-SATA-Adapter-Schutzbox-Laufwerk-wei%C3%9F/dp/B01GDZACDK). It works almost perfectly with Raspberry Pi. Unfortunately there is a bug in RPi's kernel that doesn't support UAS. That means you have to force usb-storage module. If you encouter same problem add this piece in */boot/cmdline.txt*:

    usb-storage.quirks=152d:0578:u

The number can be found via *lsusb* and don't forget the *:u* at the end. After reboot RPi won't try to use UAS and the driver will work just fine.

### Raspberry Pi blade

This blade supports any Raspberry Pi type A and B from version 2 to the most recent ones. The whole case is designed for Raspberry Pi so the holes on the side panels follow RPis dimensions and they let you to connect anything what RPi allows you without much trouble.

## The main features

This project is designed to cover my needs but feel free to send me an idea or pull request. What ever you do to help this project, please keep in mind, it's focused mainly to hobbyist and home users.

### Security

All data you send to RaspiRack should be encrypted. Your files like photos, documents or simple computer backups shouldn't be available publicly if somebody steals your Raspberry Pi. The project can't protect you from everything. It's still possible to run different kinds of attacks to your data but holding the hard drive physically shouldn't guarantee access to the data inside it.

### Backup

In home based environment RAID covers situations that won't occur much. For example if one hard drive fails usually you are ok to wait a little bit until new drive arrives so it can be fixed. What you do care is the data. This project is focused to deliver the best backup solution instead of 100 % uptime. It uses [Restic](https://restic.net/) that supports many backends and allows doing backups every hour without significant load for your internet connection.

### Simplicity

I would like to keep the future user interface as simple as possible. I prefer plug&play solutions over the most flexible ones. Sometimes it's going hand by hand, sometimes flexibility takes a price on usability. That's why I decided to use Ansible and Raspbian because the basic product can be delivered in days instead of weeks. Same time it's very easy to configure the whole project and set it up with a single command.
