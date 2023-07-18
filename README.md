# FediFriday Winlink Net Form

WARNING: This is experimental, could use some more work, exercise caution.

## Usage: For the participants

You want the contents of [form](form) directory (specifically, [form.zip](/Mihara/ffwn/releases/latest/download/ffwn.zip)) to be installed into your Winlink client to create a message from a Winlink form. The advantage of using the form is that the resulting message is guaranteed to be machine-readable, which saves hair on the net controller. If this is impossible (like when checking in with APRS) you will have to manually format the message:

```
To: <net controller's callsign>
Subject: FFWN

<callsign>,<firstname>,<city>,<state/province/locale>,<country>,<mastodon username>,<VHF/HF/APRS/Telnet/Other>
<Your answer to this week's freeform question.>
```

But if you can, please use the form. Refer to numerous tutorials around the net and YouTube on how to create a message using a form.

### Installing for RMS Express

Copy the contents of [form](form) directory to `<drive>:/<Winlink installation directory>/Global Folders/Templates` -- by default that's `C:\RMS Express\Global Folders\Templates`

### Installing for Pat

Copy the contents of [form](form) directory to `~/.local/share/pat/Standard_Forms/FFWN`

## For the net controller

This repository includes a program which will automatically process form replies into a CSV file for you to publish or do statomancy on. The program works on XML message files exported from RMS Express, or directly on your Pat inbox. Both usages will require some familiarity with the concept of command line.

See `ffwn-checkout -h` for more detailed help, but the basic gist is like this:

### Working with Winlink Express or WoAD

WoAD mail export format is the same as RMS Winlink, so files exported from WoAD can be processed the same way, which is one of the reasons we're using an export file at all. *(The other reason is that the structure of RMS Express mailboxes is not documented, while export files are at least known to be stable.)*

Using [Termux](https://termux.dev/) and the executable built for Raspberry, you can process checkin messages on Android, without involving a PC at all.

1. Select all the check-in messages in your inbox.
2. Export them as an XML file.
3. Invoke ffwn-checkout like so:

```bash
ffwn-checkout.exe rms <file.xml>
```

Where `file.xml` is the file you exported.

This will produce an `output.csv` file containing all data from the messages. If any message contains data that can't be parsed, the program will print its ID and continue, so you can then deal with the offending message manually.

### Working with Pat

Assuming the messages are still in your inbox, and the inbbox contains no messages from *prior* checkins, you just run it like this:

```bash
ffwn-checkout pat <your callsign>
```

Otherwise you need to clean up the inbox so that it does not contain messages with the subject `FFWN` that should not be processed. (All other messages will be ignored.)

Just like when processing an XML file export, messages that cannot be parsed will report their message IDs so you can deal with them manually.

## License

This program and form are released under the terms of [WTFPL2](https://en.wikipedia.org/wiki/WTFPL).

