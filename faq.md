---
layout: splash
title: Frequently Asked Questions
---

# Frequently Asked Questions

## General

### What is Restreamer?
Restreamer is a software which allows for free real-time video streaming, e.g. transferring video signals of webcams or other video-inputs directly to your website. You can be your own Restreamer streaming provider.

### How do I realise video streaming without Restreamer?  
The classic way is the data transfer with a streaming provider who collects monthly fees. Some services are supported by advertising, others additionally sell you their own hardware. Some companies then sell your live stream in content syndication campaigns, too. As a result you lose the control of your video content. The added value of the so-called "content syndication" is not measurable for you in most cases. Restreamer allows you to transfer H.264 video in TV quality to your own website without these middlemen. You are free to choose your hardware for yourself, the use is free of charge and you remain in possession of your videodata!  

### For whom is it suitable?
Restreamer is suitable for all kinds of live video. Private and professional webcams for events or tourism marketing, in-house TV in hotels, restaurants, industrial and safety monitoring or surveillance are just a few applications. Only the number of viewers is limited. For a lot of viewers and professional applications you need background technology, you can easily connect to with Restreamer.  

### What is different from the "others"?
Restreamer does not charge monthly fees and the streaming is completely free, because it uses your own internet connection. You will be your own video streaming provider. Restreamer is Open Source. It is free to use for both private and commercial use.  

### How old is Restreamer?  
Restreamer was released 30.12.2015 on GitHub

### Is Restreamer available for my operating system?  
Restreamer can be used on any type of operating system. Windows, MacOS, Linux, BSD and Solaris are supported.

### Restreamer Notice
Restreamer provides the technology to transmit video data in real time. We request every user to respect the privacy laws of your country and the privacy requirements of every human being!  


## Basics

### What kind of hardware is supported by Restreamer?
Restreamer has been tested with single board computers like Raspberry Pi1, Pi2, and ODROID. Windows, MacOS and Linux works with Docker. 

### Which camera hardware I should buy?
You can find out about camera requirements for Restreamer in our [Requirements-Guide](wiki/requirements.html).  

### How much data upload is needed to produce a proper video live stream?  
Find a rule of thumb for required upload and data volume of your video streams in the wiki chapter about [video compression](wiki/video-compression.html).  

### Does Restreamer work with UMTS/G3 or LTE/G4?  
Yes - but consider your bandwidth carefully. An HD live stream can quickly consume over 1 GB of data PER DAY(!). If you host yourself, we recommend that you use a DSL or better VDSL connection with flatrate.  

### Where can I find a suitable network camera from?
We prepared a [suiteable selection of network cameras](wiki/buy-hardware-index.html) for you on Amazon.

### How do I set a dynamic DNS?  
For more information about Dynamic DNS take a look [in the Wiki](wiki/dynamic-dns.html).

### Is streaming on Apple iOS devices and Android possible?  
Yes - Restreamer supports iPod Touch, iPhone, iPad and Android devices, too.

### Is a professional use of Restreamer possible?  
Yes - You should definitely get in touch with us and let us advise you on what steps you have to take to realise your project.  

### Which browsers are supported?  
You will find a list of supported browser [here](features.html#supported-browser)

### What about simultaneous running with third party tools like WeeWX or Zoneminder?  
Simultaneous running third party tools is possible and tested but not docomented yet...

### Where can I get Help for my project?  
There are several ways to get help. [More here](docs/learn-more-features.html#support)

### What about audio?  
Restreamer is supporting video streaming with audio

### Icecast  
To be honest we do not know yet. If you have tested this with Icecast please tell us your results :) 

### Auto resume 
If you have set up a stream the way you like, for example, from OBS to restreamer, then to youtube, what happens when you stop the stream in OBS? 

The streams in Restreamer that are currently connected will try to reconnect on loss of connection. Restreamer stores the wanted state ("connected" or "disconnected") on disk such that it can reconstruct them after a reboot.

In the example, you are pulling the stream from the local RTMP. This is the stream you send to Restreamer from OBS. If you stop OBS, Restreamer gets disconnected from the RTMP and will try to reconnect every few seconds. As soon as you start OBS again, it will be able to connect to the local RTMP.

The same happens for the stream that you send to Youtube. This one also pulls the stream from the local RTMP and forwards it to Youtube. This will also try to reconnect if the incoming stream stops.

You do not have to have the browser open for this process.

## Support the project
The easiest way to support us is to tell everyone how cool Restreamer is. Post it on Facebook, Twitter, Linkedin, or any social network you prefer.

### Donation
Please check out our [Donation](donate.html) page.

### Software developer  
As a software developer, feel free to get in touch with us :)

### Help with any other things.
If you like Restreamer and you want to help us even if you are not a software developer, get in contact with us. We have a lot of things to do around the Restreamer project. E.g we are not native english speakers and could need some help writing and text corrections.
