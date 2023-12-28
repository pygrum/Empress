# Empress

Empress is an implant with backdoor functionality, and the first implant to be integrated with
[Monarch](https://github.com/pygrum/monarch) C2 framework.

It was developed alongside Monarch for testing and POCs, and should _NOT_ be used in production. Reasonable
and ethical modifications for authorized engagements are welcome. To contribute, contact @Pygrum on the 
[BloodHound Gang Slack](https://bloodhoundgang.herokuapp.com/).

## Features

Not many really. It's a pretty basic backdoor.

- Custom TLS-secured TCP and HTTP C2 communication protocols
- Beacon mode (HTTP) and session mode (HTTP/TCP)
- Reports compromised system information on callback, which includes:

    - Host OS
    - Host architecture
    - Username
    - Hostname
    - User ID
    - Group ID
    - User's home folder

- Standard filesystem operations (change directory, list files, etc.)
- Standard networking operations (view interfaces, upload and download files)
- Standard process operations (execute commands, kill processes, list processes)
- View environment variables
- View current user information

### Disclaimer

This software is intended for authorized and lawful testing only.
Any unauthorized or illegal activities facilitated by this software are strictly prohibited.
The developers are not liable for any misuse or illegal actions performed with this software.
Users must comply with all applicable laws and ethical standards when using this software.
The developers disclaim responsibility for any damages or legal consequences resulting from its misuse.
By using this software, you agree to use it responsibly and strictly for lawful purposes.