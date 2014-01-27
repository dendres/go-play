
Event Collection and Processing
===============================

Provide 3 interfaces to access log messages:

* log.io: provide a lossy, low-latency path from event source to the browser. collect events only after the browser app is loaded
* kibana: index data from 5min ago through the last 12 hours
* kibana: historical data: tmp_index = load_in_index(topic, start_time, end_time). tmp_index expires after 48 hours.

client app:

* input plugins accept various event serializations and normalize events
* generate guid fields for each message
* serialize and buffer
* send immediately to log.io
* serialize, compress, buffer, and ensure centralized delivery to kafka within 5min

app1 log.io server:

* cluster receives and deduplicates messages from clients
* accept pub-sub log.io connections


Goals
=====

* provide reliable replay of long-term log and metric source data
* reduce load on indexed/structured data sets (lucene indexes and whisper)


Anti-Goals
==========

* years worth of data indexed and available for query with millisecond latency
* rely on memory hungry indexes as the long term storage mechanism for log data


Stage 1
=======

syslog -> long term storage buckets -> kafka -> kibana 12 hour dropped index

briefly profile cpu and memory to see if the harvester is too expensive


Stage 2
=======

add NaCL transport encryption to client and setup corresponding proxy service on kafka server


Stage 3
=======

add:

syslog -> log.io lossy low latency

possibly through kafka, but probably directly to a log.io cluster of some kind


Stage 4
=======

incorporate metrics

Stage 5
=======

add additional harvester input plugins
