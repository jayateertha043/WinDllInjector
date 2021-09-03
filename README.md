<h1 align="center">WinDllInjector</h1>

>Simple Dll injector tool to inject dll into a running Windows process.


## REQUIREMENTS AND INSTALLATION

Installation:
```
go get github.com/jayateertha043/WinDllInjector
```

Run WinDllInjector:

```
.\WinDllInjector -h
```

## USAGE:

```
 _    _  _        ______  _  _  _____           _              _
| |  | |(_)       |  _  \| || ||_   _|         (_)            | |
| |  | | _  _ __  | | | || || |  | |   _ __     _   ___   ___ | |_   ___   _ __
| |/\| || || '_ \ | | | || || |  | |  | '_ \   | | / _ \ / __|| __| / _ \ | '__|
\  /\  /| || | | || |/ / | || | _| |_ | | | |  | ||  __/| (__ | |_ | (_) || |
 \/  \/ |_||_| |_||___/  |_||_| \___/ |_| |_|  | | \___| \___| \__| \___/ |_|
                                              _/ |
                                             |__/

Author:Jayateertha G
Version:1.0.0


Usage of WinDllInjector.exe:
  -dll string
        Absolute Path to DLL
  -pid uint
        Process ID where DLL should be injected
```
The tool takes 2 parameters absolute path to dll and process id to which dll should be injected:

```
.\WinDllInjector -dll="D:\test\test.dll" -pid=1234  
```
## Author

👤 **Jayateertha G**

* Twitter: [@jayateerthaG](https://twitter.com/jayateerthaG)
* Github: [@jayateertha043](https://github.com/jayateertha043)

