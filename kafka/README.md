
Event Collection and Processing Goals
=====================================

* reduce cost of indexed/structured data sets (lucene indexes and whisper)
* compensate for lack of reliability in high performance structured data sets
* provide reliable replay of long-term log and metric source data
* minimize delay for log consumers


Anti-Goals
==========

* years worth of data indexed and available for query with millisecond latency
* rely on memory hungry indexes as the long term storage mechanism for log data


Interfaces
==========

Provide 3 interfaces to access log messages:

* log.io starting now: lossy, low-latency path from message source to browser. collect events only after the browser app is loaded
* kibana recent: index data from 5min ago through the last X days (determined by adjustable pruning policy)
* kibana history: tmp_index = load_in_index(start_time, end_time). tmp_index expires after X hours (adjustable pruning policy).


Assumptions
===========
* multiple instances of these interfaces will be created to separate production from non-production systems
* for a given pipeline instance, data will only be retrieved by timestamp. message attributes will not be available outside the compressed block.


Services
========

harvester:

* runs on ever server that produces log and metric events
* input plugins accept various event serializations and normalize events enough to determine how to forward them
* add any missing fields (hostname, process/service etc.., timestamp, some kind of message id )
* remove any duplicate fields
* send to log.io server based on filter criteria
* buffer 30 seconds
* gzip
* kafka topic: "h-<environment_name>"
* kafka partition: based on timestamp... hash timestamp or mod timestamp or something????
* send to kafka


compressor:

* determines the time interval to work in
* generates the output timestamp (number of 5min intervals since unix epoc) base36
* kafka consumer collects all messages for a given 5 minute period (10 messages from each host)
* sorts them into a single message containing all events for the 5min period
* applies bz2 compression
* sends the block back to another kafka topic "c-<environment_name>"
* sends the block to s3 as

es_consumer:

* watch topic for the assigned environment "h-<environment_name>"
* pub/sub receive all messages on the topic
* process each message into an elasticsearch bulk import


indexer:

* kafka and s3 client
* takes (environment_name, start_time, end_time)
* generates a tmp index name based on the input paramaters and return it to the caller
* makes a list of 5min time blocks to retrieve
* downloads blocks from kafka and/or s3
* stream messages out of the compressed blocks
* form elasticsearch bulk import messages
* import data into elasticsearch, and then delete each block


log.io receiver:

* cluster receives and deduplicates messages UDP from clients
* accept pub-sub log.io connections
* broker messages through to log.io with minimal latency



Stage 1
=======

harvester service and kafka:

* syslog -> kafka -> es_consumer -> kibana recent index
* add optional NaCL transport encryption to client and setup corresponding proxy service on kafka server


Stage 2
=======

* harvester -> log.io


Stage 3
=======

* compressor and indexer services


Stage 4
=======

* expand harvester input plugins
