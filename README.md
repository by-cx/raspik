# Raspik, a Raspberry Pi based NAS project (former RaspiRack)

I spent some time designing a new thing that's gonna replace RaspiRack in my home. I call it Raspik (or Raspík in Czech). I trashed the idea of having universal box and I went for a proper NAS solution for all of my digital stuff. This project transforms with it and it's now renamed because it can be hardly called a rack. The new box is easier to build, looks better, it can support two 3.5" HDD and includes proper cooling for all components. It has a small OLED display that can show anything you want. After more than one year on RaspiRack I can say that Raspberry Pi based NAS is a pretty solid solution and this is only next step in my journey to build the perfect home storage.

I am also user of Synology DS216play and this project is gonna replace this device for most of my stuff. I found my self using only fraction of its features and also missing one pretty important one. I ended up using only Samba shares on this baby so simply anything can replace it now but what I really miss there is a proper encryption. Synology uses filesystem based encryption [eCryptFS](https://www.ecryptfs.org/) which I don't see as the best way how to secure my data. It's also very slow on this machine so not everything can be encrypted without sacrificing large portion of the performance.

Goal of this project is simple - replace my commercial NAS with a comparable DIY solution where I can balance pros and cons in my own way. I offer this project to open source community as it is. You can use it as you wish and I will be more than happy if you build something based on stuff in this repository or send a pull request with an awesome idea.

![Raspík, the NAS](https://raw.githubusercontent.com/by-cx/raspik/master/photos/home.png)

## Parts needed to build this

The NAS itself, without the drives, can be build under $100. You can skip OLED display, use cheaper fans, skip dust filters, use 1 GB Raspberry Pi (which is more than enough) and use cheap SATA-USB adapters. Let's see the part list:

* Raspberry Pi 4 with any amount of RAM ($35)
* MEAN WELL RS-75-12 12V power source  ($23)
* 12 V to 5 V step down converter  ($4)
* SATA to USB adapters (~$4+ each)
* OLED display (~$4)
* 40mm and 80mm fans (~$10+)
* wire that supports 6A (~$2)
* EURO K241 socket ($1)
* pin headers sockets for the OLED display ($2)
* heat sink for Raspberry's CPU ($1)

If we sum the cheapest variant for two disks it's $90 plus printed parts. You will need around 500g of material to print this which is $10. If you have some parts already or you will find great deals on eBay you can get easily under $100. Unfortunately there is one part you shouldn't consider to pay a little bit more for. It's the SATA to USB adapter. I used ones with ASM1153 chipset which is pretty stable in current Raspbian. I have no experience from the cheap ones from eBay but I think it worths for an attempt. If you find one that's stabile including UASP, please let me know. I will add it here. Careful about adapters without 12 V connector. If there is none you have to solder two wires directly to the connector. It's [not hard though](https://www.youtube.com/watch?v=bS5Wsu1iSsY).

Based on selected SATA-USB adapters you will need a custom HDD covers. The cover is not totally needed because HDDs are screwed to the main body but in my case I wanted to protect the hard drive connectors so I designed HDD cover that prevents the connector to be broken with wrong move. But if you have the box somewhere out of reach no one will care.

Power consumption of this machine is less than 20 W but during boot the consumption can be 40+ W so 75 W is a reasonable choice.

OLED display is probably not the best solution for this kind of use but it has very convenient size. It doesn't have to be [a total catastrophe](https://hackaday.com/2019/04/23/a-year-long-experiment-in-oled-burn-in/) at the end anyway.

Pick the fans based on where your NAS will be located. I have it under my working table so I wanted something more silent and picked up Noctua. I also use dust filters. You have to clean them every few weeks but they keep inside of the box almost out of dust. The fan covers are designed for my Noctua fans which have a small piece of rubber in each corner. I hope it will fit to regular fans too. The source is attached so it's not hard to update it and reprint it. If you have an idea how to make it more universal, please let me know.

To print the main body use PETG, ABS, ASA or anything else that can withstand 55+°C. If the fan fails the heat from the power source and hard drives goes up pretty quickly. All other covers can be made from PLA if you prefer that material like me. The box is designed with Prusa MINI in mind. The body takes 24h to print with 0.4mm nozzle. 14h if you use 0.6mm nozzle. I used MK3S to print the body and MINI for everything else with no issues at all.

## Software

The software part is not done yet and I will fill up the information later. Feel free to use projects like [OpenMediaVault.org](https://www.openmediavault.org/). It's probably better solution for ordinary people compared to what I am trying to achieve here.
