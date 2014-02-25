

Analysis to complete
====================

* duralog: How durable does that logged event have to be?
* sortalog: Why sort logged events by time? and Longer wait, better sort!
* headless: Masterless is more stable and just as fast.
* shrinkalog: pre-compress events for long term storage


pathologies of big data
http://queue.acm.org/detail.cfm?id=156387

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


Event Delivery by Time
======================
* separate protocol layer delivers events and allows query from consumers
* content is a contract between producers and consumers and is NOT touched by the event delivery system



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



