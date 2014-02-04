
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


Services
========

harvester:

* runs on ever server that produces log and metric events
* input plugins accept various event serializations and normalize events enough to determine how to forward them
* log file tail: process lines starting at the end of the file. go back max 24 hours. journal by timestamp
* add any missing fields (hostname, process/service etc.., timestamp, some kind of message id )
* remove any duplicate fields
* arbitrarily buffer a few messages and gzip
* kafka topic: "h-<environment_name>", partition: hostname.to_i(36) % partition_count=15
* send to kafka over optional encrypted channel

sorter:

* manually assign partitions

* divide timestamp buckets up evenly across sorter servers?????
* each message received, if it's remote, send it to the other sorter, if local, write to local bucket

https://cwiki.apache.org/confluence/display/KAFKA/FAQ#FAQ-HowdoIchoosethenumberofpartitionsforatopic?

XXXXXXXXXXXXXXXXXXX handle bucket reassignment on failover?????

* pick the oldest bucket older than 48 hours
* sort messages, bz2, and send to s3



* if a sorter dies, either manually reassign partitions, or build another sorter within 72 hours
* a sorter must be able to tell if a disk has died so it can replay?????

sorter picks a 5min bucket to process:
* find the oldest bucket older than 48 hours

* pichostname.to_i(36) % partition_count




* rsync from all sorters


* if the server and disks die, the consumer can replay the topic from the beginning
* buckets that are 48 hours old: sort messages, bz2, and send to s3
* buckets that are 4-7 days old: delete



consumer group: each message published to a topic is delivered to a single consumer instance in the consumer group

multiple instances in a consumer group = queue balanced across instances
one instance per consumer group = broadcast to all instances




es_recent: XXX Better name?
* watch topic for the assigned environment "h-<environment_name>"
* pub/sub receive all messages on the topic
* process each message into an elasticsearch bulk import on the "recent" elasticsearch cluster
* drop logstash-style-day indexes after X days (configurable)

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
