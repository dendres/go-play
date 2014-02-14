
Event Collection and Processing Goals
=====================================

* reduce cost of searching logs
* compensate for a lack of reliability in high performance structured data sets


Anti-Goals
==========

* provide a single interface for all log/event query needs
* provide years worth of data: indexed and available for query with millisecond latency
* rely on memory hungry indexes as the long term storage mechanism for event data


Interfaces
==========

Provide 3 interfaces to access log messages:

* log.io "starting now": lossy. collect events only after the browser app is loaded
* kibana "recent": index data from M min ago through the last D days
* kibana "history": tmp_index = load_index(start_time, end_time). tmp_index expires after E hours/days.


Assumptions
===========
* "domain separation" means separate full copies of the log processing system and interfaces
* log processing "domain" can be defined arbitrarily to cover any number of servers
* for each "domain", "history" can only be retrieved by time.
* ONLY the broker needs to be clustered and available


Data Pattern
============

Event: one event to be logged (in-memory representation after parsing)

* point uint64:   10^-8 since unix epoch loops in year 7819
* shn   string:   short host name
* app   string:   overloaded = app_name || process_name || process_id
* line  string:   log line where time, shn, and app have been removed
* words []string: list of word-like tokens from the log line
* stems []string: list of stemmed search terms from the log line

Message: message/tim/ee/ unsync'd buffers used to build a message. pruned after 2 message ack. server sends ack after writing Bucket

* events.csv: base64(10^-8 since ti/me),shn,app,restore_original(line)... or can gob be re-opened and appended to???
* incrementally build the token list????
* incrementally build a b-tree of stemmed search terms



Message wire format: capnp

* list token strings
* list of search terms
* list Events with encoded(line) = []byte{bitset + terms}

Bucket: buckets/ti/me/<files> time and term searchable data, pruned after 7 days

* events.csv: base64(10^-8 since ti/me),shn,app,restore_original(line)
* incrementally build the frequency tree using the received token index
* incrementally build a b-tree of stemmed search terms

are you sure there is no way to store the frequency tree of tokens such that:
 - they get deduplicated
 - once a node ID is created, it never changes... yes. this is not possible bec


Static: static/ti/me/<files> created after mtime(bucket/ti/me) older than 48 hours

* fast searchable file of lowercased and stemmed search terms??? maybe cdb???
* list of token strings
* list of Events

Search: search/te/rm/sg/ohere

* list of "Static" entries containing this term and a count of how many times





log.io input format?


Elasticasearch bulk insert format?


Process and Message Pattern
===========================

* client: tail, parse, serialize, small in-memory buffer with fast timeout. send to multiple servers
* server: receive, bucket and dedup on disk (let OS manage memory and decide when to fsync). small disk queue for low-latency consumers
* SPOF "recent" process reads all buckets one bucket ago, and sends to 30 day elasticsearch. prunes elasticsearch
* SPOF "history" process waits 48 hours, collects from all receivers and archives a bucket to S3
* SPOF process prunes remotes after 30 - 45 days
* non-clustered elasticsearch behind loadbalancer with health check.
* webui to load history into separate elasticsearch


Services
========

tailer:

* runs on every server that produces log and metric events
* buffering happens ONLY in the tailed files and in memory.
* tail files by buffered inotify and record state on disk (fsync) after each block is sent on the network
* state/path/to/file. mtime = stat. 16 character, zero padded base10 offset in file + 64 characters of last line read from the file
* up to the entire content of the file may be processed if the file is rotated
* low latency between inotify and in-memory buffer should prevent loss of the last messages before rotate.
* add any missing fields (hostname, process/service etc.., timestamp, some kind of message id )
* remove any duplicate fields
* keep a long-standing compressed and encrypted messaging channel to multiple servers.
* sleep and retry on server disconnect or error in compression, encryption, or tcp.
* clean disconnect and reconnect periodically


buck: (holds the data)

* per-tailer: receive, decrypt, decompress, write to buffered pub/sub channel
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

* all_syslog_file -> tailer -> buck -> disk buckets


Stage 2
=======

* porter -> recent_es <- kibana "recent"


Stage 3
=======

* add log.io
* add archie interface
* add optional NaCL transport encryption



