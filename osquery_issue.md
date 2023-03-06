# Bug report

### What operating system and version are you using?

version = 20.04.2 LTS (Focal Fossa)
   build = 
platform = ubuntu

(Also reproduces on Ubuntu 22.04)

### What version of osquery are you using?

version = 5.7.0

(Tested on 5.3.0 and also reproduces)


### What steps did you take to reproduce the issue?

For the project I am working on, I use osqueryd in the following scenario:

* I have an osquery table extension which creates 3 osquery tables and continuously populates them
* Osqueryd is deployed with 3 scheduled queries, one for each table, the scheduled interval for each query is 10 seconds
* I have a logging extension which posts all generated logs to a web server.

The setup can be simulated using the code and instructions at https://github.com/mirceatanase1994/osquery-bpf-extension-issu which contains:
* the source code and pre-built binaries for the extensions necessary for reproducing the issue, which are based on the examples at https://github.com/osquery/osquery-go
* the osquery.conf, osquery.flags file necessary for reproducing the issue
* helper scripts for building the extensions and preparing the environment

After the setup is complete, osqueryctl is ran by:

` sudo osqueryctl start`


### What did you expect to see?

The extensions run and events are posted indefinetely.

### What did you see instead?


While neither the osqueryd or the extensions processes crash, at some point the table extension goes away, after what seems to be a random period of time between 0 and 48 hours.

Concretely, the following log appears:
`Mar 05 19:05:45 lurkerserver osqueryd[120207]: I0305 19:05:45.501915 120345 extensions.cpp:348] Extension UUID 28863 has gone away`

After that, each scheduled query generates a log line like this:
`Mar 05 19:05:46 lurkerserver osqueryd[120207]: E0305 19:05:46.793545 120363 scheduler.cpp:128] Error executing scheduled query foobar: vtable constructor failed: foobar`

### Notes

After reproducing the issue multiple times, I observed that the `Extension UUID XXXXX has gone away` log always appears after the `The BPF system state tracker has been successfully restarted` log. This lead me to the conclusion that the issue is somehow related to the bpf events generation, which was confirmed by the fact that it did not reproduce when the `enable_bpf_events` osquery flag was set to false.

While I was not able to find a way to reproduce the issue faster (sometimes it took 24-48 hours for reproduction), the reproduction rate seems to be 100%.

The issue also reproduces if running the osqueryd process directly, not only when using osqueryctl.

The issue does not seem to be related to the resource consumption of the extension (I tried writing extensions which use more memory/CPU and the osquery behaviour was different). It also does not seem to be related to the volume of the data which is generated/logged.

The extension process does not crash, and the `/var/osquery/osquery.em.XXXXX` sockets it opens are still open after osquery declares that "the extension has gone away"

The only log which is related to the issue is:
`Mar 05 19:05:45 lurkerserver osqueryd[120207]: I0305 19:05:45.501915 120345 extensions.cpp:348] Extension UUID 28863 has gone away`
Except for the `bpf system state tracker has been reset` log, no other logs seem to be related, although the verbose flag is set to true.
Here is a journalctl log dump for the osqueryd service:

```Mar 05 19:05:35 lurkerserver osqueryd[120207]: I0305 19:05:35.429879 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:35 lurkerserver osqueryd[120207]: I0305 19:05:35.429939 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:36 lurkerserver osqueryd[120207]: I0305 19:05:36.787986 120363 scheduler.cpp:120] Executing scheduled query foobar: SELECT * FROM foobar;
Mar 05 19:05:36 lurkerserver osqueryd[120207]: I0305 19:05:36.791887 120363 scheduler.cpp:120] Executing scheduled query foobar2: SELECT * FROM foobar2;
Mar 05 19:05:36 lurkerserver osqueryd[120207]: I0305 19:05:36.793540 120363 scheduler.cpp:186] Found results for query: foobar2
Mar 05 19:05:36 lurkerserver osqueryd[120209]: 2023/03/05 19:05:36 string: {"name":"foobar2","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0>
Mar 05 19:05:36 lurkerserver osqueryd[120209]: Got error for log {"name":"foobar2","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0,"counter">
Mar 05 19:05:36 lurkerserver osqueryd[120209]: 2023/03/05 19:05:36 string: {"name":"foobar2","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0>
Mar 05 19:05:36 lurkerserver osqueryd[120209]: Got error for log {"name":"foobar2","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0,"counter">
Mar 05 19:05:36 lurkerserver osqueryd[120207]: I0305 19:05:36.797371 120363 scheduler.cpp:120] Executing scheduled query foobar3: SELECT * FROM foobar3;
Mar 05 19:05:36 lurkerserver osqueryd[120207]: I0305 19:05:36.801471 120363 scheduler.cpp:186] Found results for query: foobar3
Mar 05 19:05:36 lurkerserver osqueryd[120209]: 2023/03/05 19:05:36 string: {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0>
Mar 05 19:05:36 lurkerserver osqueryd[120209]: Got error for log {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0,"counter">
Mar 05 19:05:36 lurkerserver osqueryd[120209]: 2023/03/05 19:05:36 string: {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0>
Mar 05 19:05:36 lurkerserver osqueryd[120209]: Got error for log {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0,"counter">
Mar 05 19:05:36 lurkerserver osqueryd[120209]: 2023/03/05 19:05:36 string: {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0>
Mar 05 19:05:36 lurkerserver osqueryd[120209]: Got error for log {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0,"counter">
Mar 05 19:05:36 lurkerserver osqueryd[120209]: 2023/03/05 19:05:36 string: {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0>
Mar 05 19:05:36 lurkerserver osqueryd[120209]: Got error for log {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0,"counter">
Mar 05 19:05:36 lurkerserver osqueryd[120209]: 2023/03/05 19:05:36 string: {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0>
Mar 05 19:05:36 lurkerserver osqueryd[120209]: Got error for log {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0,"counter">
Mar 05 19:05:36 lurkerserver osqueryd[120209]: 2023/03/05 19:05:36 string: {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0>
Mar 05 19:05:36 lurkerserver osqueryd[120209]: Got error for log {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0,"counter">
Mar 05 19:05:36 lurkerserver osqueryd[120209]: 2023/03/05 19:05:36 string: {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0>
Mar 05 19:05:36 lurkerserver osqueryd[120209]: Got error for log {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0,"counter">
Mar 05 19:05:37 lurkerserver osqueryd[120209]: 2023/03/05 19:05:37 string: {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0>
Mar 05 19:05:37 lurkerserver osqueryd[120209]: Got error for log {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0,"counter">
Mar 05 19:05:37 lurkerserver osqueryd[120209]: 2023/03/05 19:05:37 string: {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0>
Mar 05 19:05:37 lurkerserver osqueryd[120209]: Got error for log {"name":"foobar3","hostIdentifier":"lurkerserver","calendarTime":"Sun Mar  5 19:05:36 2023 UTC","unixTime":1678043136,"epoch":0,"counter">
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019722 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019778 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019789 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019800 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019809 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019820 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019829 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019840 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019860 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019870 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019887 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019896 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019906 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019915 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019932 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019941 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019968 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019979 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.019995 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.020005 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.020020 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.020030 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.020045 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.020053 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.020068 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.020077 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.020092 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.020100 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.020115 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:39 lurkerserver osqueryd[120207]: I0305 19:05:39.020123 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:41 lurkerserver osqueryd[120207]: I0305 19:05:41.621335 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:41 lurkerserver osqueryd[120207]: I0305 19:05:41.621397 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:44 lurkerserver osqueryd[120207]: I0305 19:05:44.432142 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:44 lurkerserver osqueryd[120207]: I0305 19:05:44.432191 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:45 lurkerserver osqueryd[120207]: I0305 19:05:45.501714 120358 bpfeventpublisher.cpp:337] The BPF system state tracker has been successfully restarted
Mar 05 19:05:45 lurkerserver osqueryd[120207]: I0305 19:05:45.501915 120345 extensions.cpp:348] Extension UUID 28863 has gone away
Mar 05 19:05:46 lurkerserver osqueryd[120207]: I0305 19:05:46.670588 120358 systemstatetracker.cpp:290] Created new process context from procfs for pid 120208 some fields may be not accurate
Mar 05 19:05:46 lurkerserver osqueryd[120207]: I0305 19:05:46.791029 120363 scheduler.cpp:120] Executing scheduled query foobar: SELECT * FROM foobar;
Mar 05 19:05:46 lurkerserver osqueryd[120207]: E0305 19:05:46.793545 120363 scheduler.cpp:128] Error executing scheduled query foobar: vtable constructor failed: foobar
Mar 05 19:05:46 lurkerserver osqueryd[120207]: I0305 19:05:46.793752 120363 scheduler.cpp:120] Executing scheduled query foobar2: SELECT * FROM foobar2;
Mar 05 19:05:46 lurkerserver osqueryd[120207]: E0305 19:05:46.795969 120363 scheduler.cpp:128] Error executing scheduled query foobar2: vtable constructor failed: foobar2
Mar 05 19:05:46 lurkerserver osqueryd[120207]: I0305 19:05:46.796121 120363 scheduler.cpp:120] Executing scheduled query foobar3: SELECT * FROM foobar3;
Mar 05 19:05:46 lurkerserver osqueryd[120207]: E0305 19:05:46.798261 120363 scheduler.cpp:128] Error executing scheduled query foobar3: vtable constructor failed: foobar3
Mar 05 19:05:47 lurkerserver osqueryd[120207]: I0305 19:05:47.512607 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:47 lurkerserver osqueryd[120207]: I0305 19:05:47.512725 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:50 lurkerserver osqueryd[120207]: I0305 19:05:50.513999 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:53 lurkerserver osqueryd[120207]: I0305 19:05:53.514642 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:56 lurkerserver osqueryd[120207]: I0305 19:05:56.511538 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:05:56 lurkerserver osqueryd[120207]: I0305 19:05:56.793340 120363 scheduler.cpp:120] Executing scheduled query foobar: SELECT * FROM foobar;
Mar 05 19:05:56 lurkerserver osqueryd[120207]: E0305 19:05:56.794062 120363 scheduler.cpp:128] Error executing scheduled query foobar: vtable constructor failed: foobar
Mar 05 19:05:56 lurkerserver osqueryd[120207]: I0305 19:05:56.794114 120363 scheduler.cpp:120] Executing scheduled query foobar2: SELECT * FROM foobar2;
Mar 05 19:05:56 lurkerserver osqueryd[120207]: E0305 19:05:56.794622 120363 scheduler.cpp:128] Error executing scheduled query foobar2: vtable constructor failed: foobar2
Mar 05 19:05:56 lurkerserver osqueryd[120207]: I0305 19:05:56.794664 120363 scheduler.cpp:120] Executing scheduled query foobar3: SELECT * FROM foobar3;
Mar 05 19:05:56 lurkerserver osqueryd[120207]: E0305 19:05:56.795167 120363 scheduler.cpp:128] Error executing scheduled query foobar3: vtable constructor failed: foobar3
Mar 05 19:05:59 lurkerserver osqueryd[120207]: I0305 19:05:59.515892 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:02 lurkerserver osqueryd[120207]: I0305 19:06:02.516677 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:05 lurkerserver osqueryd[120207]: I0305 19:06:05.517511 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:06 lurkerserver osqueryd[120207]: I0305 19:06:06.796022 120363 scheduler.cpp:120] Executing scheduled query foobar: SELECT * FROM foobar;
Mar 05 19:06:06 lurkerserver osqueryd[120207]: E0305 19:06:06.796697 120363 scheduler.cpp:128] Error executing scheduled query foobar: vtable constructor failed: foobar
Mar 05 19:06:06 lurkerserver osqueryd[120207]: I0305 19:06:06.796746 120363 scheduler.cpp:120] Executing scheduled query foobar2: SELECT * FROM foobar2;
Mar 05 19:06:06 lurkerserver osqueryd[120207]: E0305 19:06:06.797228 120363 scheduler.cpp:128] Error executing scheduled query foobar2: vtable constructor failed: foobar2
Mar 05 19:06:06 lurkerserver osqueryd[120207]: I0305 19:06:06.797257 120363 scheduler.cpp:120] Executing scheduled query foobar3: SELECT * FROM foobar3;
Mar 05 19:06:06 lurkerserver osqueryd[120207]: E0305 19:06:06.797729 120363 scheduler.cpp:128] Error executing scheduled query foobar3: vtable constructor failed: foobar3
Mar 05 19:06:08 lurkerserver osqueryd[120207]: I0305 19:06:08.518563 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:11 lurkerserver osqueryd[120207]: I0305 19:06:11.511498 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:14 lurkerserver osqueryd[120207]: I0305 19:06:14.520639 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:16 lurkerserver osqueryd[120207]: I0305 19:06:16.798521 120363 scheduler.cpp:120] Executing scheduled query foobar: SELECT * FROM foobar;
Mar 05 19:06:16 lurkerserver osqueryd[120207]: E0305 19:06:16.799201 120363 scheduler.cpp:128] Error executing scheduled query foobar: vtable constructor failed: foobar
Mar 05 19:06:16 lurkerserver osqueryd[120207]: I0305 19:06:16.799250 120363 scheduler.cpp:120] Executing scheduled query foobar2: SELECT * FROM foobar2;
Mar 05 19:06:16 lurkerserver osqueryd[120207]: E0305 19:06:16.799782 120363 scheduler.cpp:128] Error executing scheduled query foobar2: vtable constructor failed: foobar2
Mar 05 19:06:16 lurkerserver osqueryd[120207]: I0305 19:06:16.799816 120363 scheduler.cpp:120] Executing scheduled query foobar3: SELECT * FROM foobar3;
Mar 05 19:06:16 lurkerserver osqueryd[120207]: E0305 19:06:16.800324 120363 scheduler.cpp:128] Error executing scheduled query foobar3: vtable constructor failed: foobar3
Mar 05 19:06:17 lurkerserver osqueryd[120207]: I0305 19:06:17.516090 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:20 lurkerserver osqueryd[120207]: I0305 19:06:20.522521 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:23 lurkerserver osqueryd[120207]: I0305 19:06:23.523260 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:26 lurkerserver osqueryd[120207]: I0305 19:06:26.527529 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:26 lurkerserver osqueryd[120207]: I0305 19:06:26.801062 120363 scheduler.cpp:120] Executing scheduled query foobar: SELECT * FROM foobar;
Mar 05 19:06:26 lurkerserver osqueryd[120207]: E0305 19:06:26.803637 120363 scheduler.cpp:128] Error executing scheduled query foobar: vtable constructor failed: foobar
Mar 05 19:06:26 lurkerserver osqueryd[120207]: I0305 19:06:26.803858 120363 scheduler.cpp:120] Executing scheduled query foobar2: SELECT * FROM foobar2;
Mar 05 19:06:26 lurkerserver osqueryd[120207]: E0305 19:06:26.806042 120363 scheduler.cpp:128] Error executing scheduled query foobar2: vtable constructor failed: foobar2
Mar 05 19:06:26 lurkerserver osqueryd[120207]: I0305 19:06:26.806191 120363 scheduler.cpp:120] Executing scheduled query foobar3: SELECT * FROM foobar3;
Mar 05 19:06:26 lurkerserver osqueryd[120207]: E0305 19:06:26.808384 120363 scheduler.cpp:128] Error executing scheduled query foobar3: vtable constructor failed: foobar3
Mar 05 19:06:29 lurkerserver osqueryd[120207]: I0305 19:06:29.525305 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:32 lurkerserver osqueryd[120207]: I0305 19:06:32.525720 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:35 lurkerserver osqueryd[120207]: I0305 19:06:35.526474 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:36 lurkerserver osqueryd[120207]: I0305 19:06:36.802934 120363 scheduler.cpp:120] Executing scheduled query foobar: SELECT * FROM foobar;
Mar 05 19:06:36 lurkerserver osqueryd[120207]: E0305 19:06:36.803742 120363 scheduler.cpp:128] Error executing scheduled query foobar: vtable constructor failed: foobar
Mar 05 19:06:36 lurkerserver osqueryd[120207]: I0305 19:06:36.803794 120363 scheduler.cpp:120] Executing scheduled query foobar2: SELECT * FROM foobar2;
Mar 05 19:06:36 lurkerserver osqueryd[120207]: E0305 19:06:36.804306 120363 scheduler.cpp:128] Error executing scheduled query foobar2: vtable constructor failed: foobar2
Mar 05 19:06:36 lurkerserver osqueryd[120207]: I0305 19:06:36.804337 120363 scheduler.cpp:120] Executing scheduled query foobar3: SELECT * FROM foobar3;
Mar 05 19:06:36 lurkerserver osqueryd[120207]: E0305 19:06:36.804809 120363 scheduler.cpp:128] Error executing scheduled query foobar3: vtable constructor failed: foobar3
Mar 05 19:06:38 lurkerserver osqueryd[120207]: I0305 19:06:38.527258 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:41 lurkerserver osqueryd[120207]: I0305 19:06:41.521417 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:44 lurkerserver osqueryd[120207]: I0305 19:06:44.523910 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:46 lurkerserver osqueryd[120207]: I0305 19:06:46.805366 120363 scheduler.cpp:120] Executing scheduled query foobar: SELECT * FROM foobar;
Mar 05 19:06:46 lurkerserver osqueryd[120207]: E0305 19:06:46.808229 120363 scheduler.cpp:128] Error executing scheduled query foobar: vtable constructor failed: foobar
Mar 05 19:06:46 lurkerserver osqueryd[120207]: I0305 19:06:46.808429 120363 scheduler.cpp:120] Executing scheduled query foobar2: SELECT * FROM foobar2;
Mar 05 19:06:46 lurkerserver osqueryd[120207]: E0305 19:06:46.810649 120363 scheduler.cpp:128] Error executing scheduled query foobar2: vtable constructor failed: foobar2
Mar 05 19:06:46 lurkerserver osqueryd[120207]: I0305 19:06:46.810804 120363 scheduler.cpp:120] Executing scheduled query foobar3: SELECT * FROM foobar3;
Mar 05 19:06:46 lurkerserver osqueryd[120207]: E0305 19:06:46.813032 120363 scheduler.cpp:128] Error executing scheduled query foobar3: vtable constructor failed: foobar3
Mar 05 19:06:47 lurkerserver osqueryd[120207]: I0305 19:06:47.527920 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:50 lurkerserver osqueryd[120207]: I0305 19:06:50.531297 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:53 lurkerserver osqueryd[120207]: I0305 19:06:53.532088 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:56 lurkerserver osqueryd[120207]: I0305 19:06:56.532752 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:06:56 lurkerserver osqueryd[120207]: I0305 19:06:56.808058 120363 scheduler.cpp:120] Executing scheduled query foobar: SELECT * FROM foobar;
Mar 05 19:06:56 lurkerserver osqueryd[120207]: E0305 19:06:56.808773 120363 scheduler.cpp:128] Error executing scheduled query foobar: vtable constructor failed: foobar
Mar 05 19:06:56 lurkerserver osqueryd[120207]: I0305 19:06:56.808821 120363 scheduler.cpp:120] Executing scheduled query foobar2: SELECT * FROM foobar2;
Mar 05 19:06:56 lurkerserver osqueryd[120207]: E0305 19:06:56.809336 120363 scheduler.cpp:128] Error executing scheduled query foobar2: vtable constructor failed: foobar2
Mar 05 19:06:56 lurkerserver osqueryd[120207]: I0305 19:06:56.809376 120363 scheduler.cpp:120] Executing scheduled query foobar3: SELECT * FROM foobar3;
Mar 05 19:06:56 lurkerserver osqueryd[120207]: E0305 19:06:56.809877 120363 scheduler.cpp:128] Error executing scheduled query foobar3: vtable constructor failed: foobar3
Mar 05 19:06:59 lurkerserver osqueryd[120207]: I0305 19:06:59.533654 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:07:02 lurkerserver osqueryd[120207]: I0305 19:07:02.534493 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:07:05 lurkerserver osqueryd[120207]: I0305 19:07:05.534808 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:07:06 lurkerserver osqueryd[120207]: I0305 19:07:06.810482 120363 scheduler.cpp:120] Executing scheduled query foobar: SELECT * FROM foobar;
Mar 05 19:07:06 lurkerserver osqueryd[120207]: E0305 19:07:06.811842 120363 scheduler.cpp:128] Error executing scheduled query foobar: vtable constructor failed: foobar
Mar 05 19:07:06 lurkerserver osqueryd[120207]: I0305 19:07:06.811900 120363 scheduler.cpp:120] Executing scheduled query foobar2: SELECT * FROM foobar2;
Mar 05 19:07:06 lurkerserver osqueryd[120207]: E0305 19:07:06.812463 120363 scheduler.cpp:128] Error executing scheduled query foobar2: vtable constructor failed: foobar2
Mar 05 19:07:06 lurkerserver osqueryd[120207]: I0305 19:07:06.812495 120363 scheduler.cpp:120] Executing scheduled query foobar3: SELECT * FROM foobar3;
Mar 05 19:07:06 lurkerserver osqueryd[120207]: E0305 19:07:06.812992 120363 scheduler.cpp:128] Error executing scheduled query foobar3: vtable constructor failed: foobar3
Mar 05 19:07:08 lurkerserver osqueryd[120207]: I0305 19:07:08.528203 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:07:11 lurkerserver osqueryd[120207]: I0305 19:07:11.535938 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:07:14 lurkerserver osqueryd[120207]: I0305 19:07:14.537036 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:07:16 lurkerserver osqueryd[120207]: I0305 19:07:16.812369 120363 scheduler.cpp:120] Executing scheduled query foobar: SELECT * FROM foobar;
Mar 05 19:07:16 lurkerserver osqueryd[120207]: E0305 19:07:16.814916 120363 scheduler.cpp:128] Error executing scheduled query foobar: vtable constructor failed: foobar
Mar 05 19:07:16 lurkerserver osqueryd[120207]: I0305 19:07:16.815097 120363 scheduler.cpp:120] Executing scheduled query foobar2: SELECT * FROM foobar2;
Mar 05 19:07:16 lurkerserver osqueryd[120207]: E0305 19:07:16.817327 120363 scheduler.cpp:128] Error executing scheduled query foobar2: vtable constructor failed: foobar2
Mar 05 19:07:16 lurkerserver osqueryd[120207]: I0305 19:07:16.817487 120363 scheduler.cpp:120] Executing scheduled query foobar3: SELECT * FROM foobar3;
Mar 05 19:07:16 lurkerserver osqueryd[120207]: E0305 19:07:16.819655 120363 scheduler.cpp:128] Error executing scheduled query foobar3: vtable constructor failed: foobar3
Mar 05 19:07:17 lurkerserver osqueryd[120207]: I0305 19:07:17.530256 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:07:20 lurkerserver osqueryd[120207]: I0305 19:07:20.539191 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:07:23 lurkerserver osqueryd[120207]: I0305 19:07:23.536531 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:07:26 lurkerserver osqueryd[120207]: I0305 19:07:26.540452 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:07:26 lurkerserver osqueryd[120207]: I0305 19:07:26.814726 120363 scheduler.cpp:120] Executing scheduled query foobar: SELECT * FROM foobar;
Mar 05 19:07:26 lurkerserver osqueryd[120207]: E0305 19:07:26.815440 120363 scheduler.cpp:128] Error executing scheduled query foobar: vtable constructor failed: foobar
Mar 05 19:07:26 lurkerserver osqueryd[120207]: I0305 19:07:26.815500 120363 scheduler.cpp:120] Executing scheduled query foobar2: SELECT * FROM foobar2;
Mar 05 19:07:26 lurkerserver osqueryd[120207]: E0305 19:07:26.816002 120363 scheduler.cpp:128] Error executing scheduled query foobar2: vtable constructor failed: foobar2
Mar 05 19:07:26 lurkerserver osqueryd[120207]: I0305 19:07:26.816044 120363 scheduler.cpp:120] Executing scheduled query foobar3: SELECT * FROM foobar3;
Mar 05 19:07:26 lurkerserver osqueryd[120207]: E0305 19:07:26.816529 120363 scheduler.cpp:128] Error executing scheduled query foobar3: vtable constructor failed: foobar3
Mar 05 19:07:29 lurkerserver osqueryd[120207]: I0305 19:07:29.541482 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:07:32 lurkerserver osqueryd[120207]: I0305 19:07:32.542393 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:07:35 lurkerserver osqueryd[120207]: I0305 19:07:35.543226 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
Mar 05 19:07:36 lurkerserver osqueryd[120207]: I0305 19:07:36.817199 120363 scheduler.cpp:120] Executing scheduled query foobar: SELECT * FROM foobar;
Mar 05 19:07:36 lurkerserver osqueryd[120207]: E0305 19:07:36.817976 120363 scheduler.cpp:128] Error executing scheduled query foobar: vtable constructor failed: foobar
Mar 05 19:07:36 lurkerserver osqueryd[120207]: I0305 19:07:36.818024 120363 scheduler.cpp:120] Executing scheduled query foobar2: SELECT * FROM foobar2;
Mar 05 19:07:36 lurkerserver osqueryd[120207]: E0305 19:07:36.818509 120363 scheduler.cpp:128] Error executing scheduled query foobar2: vtable constructor failed: foobar2
Mar 05 19:07:36 lurkerserver osqueryd[120207]: I0305 19:07:36.818538 120363 scheduler.cpp:120] Executing scheduled query foobar3: SELECT * FROM foobar3;
Mar 05 19:07:36 lurkerserver osqueryd[120207]: E0305 19:07:36.819020 120363 scheduler.cpp:128] Error executing scheduled query foobar3: vtable constructor failed: foobar3
Mar 05 19:07:38 lurkerserver osqueryd[120207]: I0305 19:07:38.544145 120359 socket_events.cpp:310] Malformed syscall event. The saddr field in the AUDIT_SOCKADDR record could not be parsed: "0100"
```