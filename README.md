# Raspik, a Raspberry Pi based NAS project (former RaspiRack)

So I spent some time on designing a new thing that's gonna replace RaspiRack in my home. I call it Raspik (or Raspík in Czech). I trashed the idea of having universal box and I went for a proper NAS solution for all of my stuff. This project transforms with it. The new box is easier to build, looks better, it can support two 3.5" HDD and includes proper cooling for all components. After more than one year on RaspiRack I can say that Raspberry Pi based NAS is a pretty solid solution and this is only next step in my journey to build the perfect home storage.

I am also user of Synology DS216play and this project is gonna replace this device for most of my stuff. I found my self using only fraction of its features and also missing one pretty important one. I ended up using only Samba shares on this baby so simply anything can replace it now but what I really miss there is a proper encryption. Synology uses filesystem based encryption [eCryptFS](https://www.ecryptfs.org/) which I don't see as the best way how to secure my data. It's also very slow on this machine so not everything can be encrypted without sacrificing large portion of the performance.

Goal of this project is simple - replace my commercial NAS with a comparable DIY solution where I can balance pros and cons my own way. I offer this project to open source community as it is. You can use it as you wish and I will be more than happy if you send a pull request with an awesome idea.

