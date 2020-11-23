# instructions and info for autisticvegan's submission to skydb hackathon

csharp project - drop the code in and use .NET 5.0
unity/maya - nothing really here
speedtest - go build
stresstest - go build

INFO:
Windows Defender might detect the .exe you build as malware, especially if you connect to publically available proxies when stress-testing.

IMPORTANT NOTE:
proxies must be in the format with the protocol and port number
Also, if you supply a blank one, then it will just use the local machine and not use a proxy for that iteration

csharpproject - proj.exe DATA_KEY stuff_to_upload.txt
This will upload the contents of stuff_to_upload.txt to the skynet registry under the key DATA_KEY




******************************************

# stresser - 
stress.exe siasky.net d 6 proxylist.txt VACKuQGhq6HN15CEmzRIXi5PDz9KGczdEFrnC_RcFWC4sg

 explanation:
 stress - the exe
 siasky.net - the portal to use
 d - d for download
 6 - the count of simultaneous connections to open per entry in proxylist
 proxylist.txt - the list of proxies to use
 VACKuQGhq6HN15CEmzRIXi5PDz9KGczdEFrnC_RcFWC4sg - the file (skylink) to download
 
 note: this uses https://github.com/autisticvegan/go-skynet where support for proxy use was put into the skynet SDK
 
 ******************************************
 
 # speedtest - 
 speedtest.exe (no params for this one, can do bigger or smaller tests in code)
 explanation:
 scrapes siastats for a list of portals, in case of fallback it uses portals.txt
 uploads and then downloads a 1 kb file 3 times per portal, taking the median time
 saves results to a file
 
 *******************************************
 
 video:
 https://www.youtube.com/watch?v=vVTls-LtSiE&feature=youtu.be
 https://siasky.net/AACnCh-9Q0blzfDMr4GJveTvlFNO8ulh3eIGI2Tu4pNAkQ