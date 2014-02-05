
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


Services
========

tailer:

* runs on every server that produces log and metric events
* log file tail: process lines starting at the end of the file. go back max 24 hours. journal by timestamp
* add any missing fields (hostname, process/service etc.., timestamp, some kind of message id )
* remove any duplicate fields
* arbitrarily buffer a few messages and gzip
* kafka topic: "t-<environment_name>", partition: rand_int % partition_count=15
* send to kafka over optional encrypted channel
* choose a new topic and reconnect every M minutes


kafka:

* queues messages from harvester on disk for D days
* cluster of 3 to 5 nodes (starting with the smallest nodes that do not crash)
* non-raid ephemeral disks


logio:

* runs on each log.io server
* for each kafka message, extract the events and send them to the local log.io instance


porter:

* for each kafka message, form an elasticsearch bulk import
* cluster of 3 to 5 nodes
* manage D days worth of logstash-style (but smaller) "hour" elasticsearch indexes


buck:

* bucket messages onto local ephemeral filesystem then compress and archive to s3
* a SPOF server that can be offline for up to 2 days with minimal impact
* after 48 hours process the oldest buckets
* sort, bz2, send to s3


archie:

* generate elasticsearch indexes on demand from archive
* takes (environment_name, start_time, end_time)
* generates a tmp index name based on the input paramaters and return it to the caller
* makes a list of 5min time blocks to retrieve
* downloads blocks from s3
* form elasticsearch bulk import messages


Stage 1
=======

* all_syslog -> file -> tailer -> kafka -> logio
* with optional NaCL transport encryption


Stage 2
=======

* porter -> kibana "recent"


Stage 3
=======

* buck and archie


Stage 4
=======

* expand harvester input plugins
