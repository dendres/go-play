
Why sort logged events by time?
===============================

The importance of sorting events in time is determined by the event consumer. Consumers that perform monitoring event correlation and alerting are generally more sensitive to out-of-order arrival and less sensitive to latency. Consumers that aim to display events for debugging are often willing to tolerate out-of-order arrival to achieve low-latency delivery. Consumers that archive events offline can wait a long time to ensure that events have stopped coming in before packing them up and deleting the local copy.


Longer wait, better sort!

Server and service crashes and restarts lead to events arriving out of order. If a distributed, event-receiving service is tasked with maintaining a sorted list of events, then the longer an event consumer waits before querying for a given time period, the more ordered the events will be when they arrive.


how to best maintain the sorted events for high throughput and reliability



abstract (4 sentences, 100 readers)
  here is a problem
    it's an interesting problem
    it's an unsolved problem
  here is my idea/claim/contribution
    your paper should have only 1 clear, sharp, specific, idea
    if you have lots of ideas, write lots of papers.
  my idea works (details, data)


kafka design docs and
pathologies of big data
http://queue.acm.org/detail.cfm?id=1563874


talk about flume and map reduce

talk about fluentd

http://wiki.apache.org/lucene-java/LuceneImplementations
make sure you understand exactly how lucene operates




requirement:
 - throughput/efficiency should approach kafka for 10 second batching when sorted time output is irrelevant
 - for ~15min batching, with various hold down timers, throughput must also match kafka, but provide sorted output over that time interval
 - hold_down timer from 0 seconds to 48 hours
 - batch size from 10 to 1024 seconds
XXXX so this implies that either the file structure should converge to kafka's for small hold_down and batch_size
     OR... that the cost of the file structure is always the same as kafka's independent of hold_down and batch_size
   must take a much closer look at kafka's disk handling to determine!!!!


how important is sorting events by time?

how inconsistent is the delivery order of the log sending application???
 - inconsistencies are possible due to server and service restart
 - but how frequent are they really
 - and how big of a deal is the out of order queue
Can through match Kafka but also give the benefit of sorted time output?
and for outputs that don't care about time


find the name of the tree I'm using: http://en.wikipedia.org/wiki/List_of_data_structures
* split on fixed leaf size
* fixed number of new leaf nodes

