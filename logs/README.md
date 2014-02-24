
Reliable Event Log Handling
===========================

abstract (4 sentences, 100 readers)
  here is a problem
    it's an interesting problem
    it's an unsolved problem
  here is my idea/claim/contribution
    your paper should have only 1 clear, sharp, specific, idea
    if you have lots of ideas, write lots of papers.
  my idea works (details, data)


paper1:
Tunable durability allows event log processing systems to reduce excessive replication and disk sync.


paper4:
how inconsistent is the delivery order of the log sending application???
 - inconsistencies are possible due to server and service restart
 - but how frequent are they really
 - and how big of a deal is the out of order queue


paper3:
unsorted consistent queue   vs  partially sorted inconsistent queue

find the name of the tree I'm using: http://en.wikipedia.org/wiki/List_of_data_structures
* split on fixed leaf size
* fixed number of new leaf nodes



paper2:
A log event processing system with an elected master does not significantly improve delivery efficiency or reliability over a masterless, stateless system.

allow 2 server failures in a 3 to 6 server system. replication factor 2. message exists on 3 random servers. any 2 can die. compare this to the paxos case with replication factor 2 in 3,4,5,6 server systems

???????? not sure this stands on it's own.... might need to be made more specific by the need for sorting??????

... what kafka does most of the time in a 3 server system, is send every event to the other 2 servers.
    a stateless server can do this without paxos.
... no client can guarantee that events are sent in order they occurr without waiting for 48 hours or so
    so kafka can't queue the events in order
    some event consumers don't care about order
    but some do. so a sort is required... and sorting is expensive.
    so the event receiver should ensure that it's disk write contributes to the sorting effort
??????? too long. write this better!!!
  counter... but I need to scale to 4, 5, or 6 servers
    no problem, the replication count remains fixed. add more servers. choose random servers to replicate to.

XXX is a reliable and efficient way to collect, store, and use a year's worth of event logs.



Moderate durability has already been achieved by existing TCP syslog implementations, logstash, and heka, but none of these solutions address the problem of:
* failures in downstream index storage
* replay of of a time period in detail
* efficient long term storage

the first thing a reliable message delivery system must do is make sure that it's not going to loose the event.
give a detailed account of how kafka ensures message durability with master election, and replication
calculate how much time it takes for the other kafka nodes to recover from failure and show that the most effective way to recover from intermediate storage node failure is to just let new messages come in.




efficient:
Reliability can be set on a per message basis Allowing XXX to determine how to make efficient use of it's limited resources.


It is more reliable, but less efficient at event delivery than rsyslog, syslog-ng, logstash, and heka because events are delivered over multiple time intervals


pathologies of big data
http://queue.acm.org/detail.cfm?id=1563874


compare to syslog's event format:

syslog usess a parsed ASCII format designed for interoperability. see rfc5424 section 6.

the main problems with this format are:
* messages of unlimited size make time and space guarantees much more difficult to achieve. cap messages at 64KB.
* variable width 9 field date specification with ascii separator fields requires X operations to parse.
  Instead a single uint64 requires Y operations to parse.
  events occuring before unix epoch should be considered out of scope by an event log processing system.
* priority = facility,severity is difficult to use
  facility should be completely replaced by a variable length "application name" up to 256 bytes
* structured data contains parsed ascii separator characters
  instead, structured data should contain utf-8 strings accessed by offset:
  1_byte_key_length,1_byte_value_length,key,value


vs:


Serialized Event:
* version: 1 byte fixed
* point  : 8 bytes fixed
* crc    : 4 bytes fixed
* shn    : 1 byte length + data
* app    : 1 byte length + data
* marks  : 2 byte length + data
* tokens : 2 byte length + data
* line   : 2 byte length + data




talk about flume and map reduce

talk about fluentd


http://wiki.apache.org/lucene-java/LuceneImplementations
make sure you understand exactly how lucene operates













Anti-Goals
==========

* provide a single interface for all log/event query needs
* provide years worth of data: indexed and available for query with millisecond latency
* rely on memory hungry indexes as the long term storage mechanism for event data


Interfaces
==========

Provide 3 interfaces to access log messages:

* log.io "starting now": lossy. collect events only after the browser app is loaded
* "recent": index data from M min ago through the last D days
* "history": tmp_index = load_index(start_time, end_time). tmp_index expires after E hours/days.


Data Pattern
============

Event:
* version uint8:   encoding version
* point uint64:    count of 10^-8 intervals since unix epoch loops in year 7819
* crc   uint32:    crc32 checksum of "line" as it should appear in the browser
* shn   string:    short host name
* app   string:    overloaded = app_name || process_name || process_id
* marks bitset:    a bitset indicating which characters in "line" should be replaced with tokens
* tokens []string: list of word-like tokens from the log line
* line  string:    log line where time, shn, and app have been removed
* l     string:    the encoded string


Serialized Event:
* version: 1 byte fixed
* point  : 8 bytes fixed
* crc    : 4 bytes fixed
* shn    : 1 byte length + data
* app    : 1 byte length + data
* marks  : 2 byte length + data
* tokens : 2 byte length + data
* line   : 2 byte length + data

Tailer:

* buckets/tim/es.b: flat file: append events: 2 byte fixed Event length + data

Buck:
* sort????


Process and Message Pattern
===========================

* tailer: tail, parse, serialize, and send messages on I time intervals to S servers
* buck: receive events, pub/sub events, answer queries
* porter: query buck and offline data to s3
* web UI for static search by time, 3gram, term
* pub/sub buck interface for nagios and graphite


Services
========

tailer:

* runs on every server that produces log and metric events
* poll log files and track offset
* if offset after EOF, find the most recent archive and read any lines after offset. reset offset to 0 and continue polling
* parse Event: pull out time, hostname, and the name of the sending application
* language analyze Event: separate interesting words from punctuation
* keep open connections to all servers. reconnect on timeout
* for each packet worth of events, pick 2 random servers and send the message
* time bucket events to disk
* when a time bucket boundary passes, sleep random, pick a random server and send the entire file as is
* prune buckets on time and size limit
* when pruning: pick a random server and send the whole bucket, then delete from local disk.

buck:

* per-tailer: receive, write to buffered pub/sub channel
* pub/sub channel listener parses and sends on another pub/sub channel
* per-socket.io: pub/sub consume, filter (configured in browser) and broker to any connected socket.io
* per-disk: pub/sub consume, filter not my disk, append (no fsync) buckets/disks/ti/me one serialized message per line
* per-porter: receive bucket requests, pack up messages from a bucket, return the bucket


porter (pack up data and send to es and s3)

* singleton: one bucket_time ago, send bucket requests to all servers, convert messages to bulk import format, dedup messages, send to "recent" es
* singleton: pick the oldest bucket missing from S3 that is more than 48 hours old (configured history start time and bucket black list), pack up and send to s3
* singleton: prune buckets
* singleton: accept (environment_name, start_time, end_time) requests, send bucket requests to all servers, then send to "history" es

log.io (starting now browser app):

* download app including the list of all servers to try to connect to
* app starts up and connects to all servers socket.io
* sleep and retry (per server) on all server errors
* receive messages. buffer/dedup for S seconds. display in browser.

archie (web server app):

* accept web rpc request for (environment_name, start_time, end_time)
* return tmp index name
* respond to status requests for data retrieval and index build


incident recovery tools:
* cli versions of common tasks to allow spinning up additional horsepower to quickly make new indexes
* export from ES, back into individual messages: port backlog out of ES and into S3


Stage 1
=======

prove technology:
* tail file
* parse log lines
* extract tokens from line
* serialize message
* zlib compress
* tcp connection handling and health
* NaCL encryption
* pub/sub channel
* listB
* pick kv
* ?????????


Stage 2
=======

* all_syslog_file -> tailer -> buck -> disk buckets
* porter -> recent_es <- kibana "recent"


Stage 3
=======

* add log.io
* add archie interface
* add optional NaCL transport encryption



