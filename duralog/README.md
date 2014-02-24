
How durable does that logged event have to be?
==============================================

While it is possible to consider durability a property of the event processing and/or consuming services, the business need for event delivery can be more directly tied to the cost of event processing when applications specify how durable they would like each event to be.


Counter Examples??

Tests???

statement must be too weak... make something testable!!!!


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


Serialized Event:
* version: 1 byte  fixed
* crc    : 4 bytes fixed
* durab  : 1 byte  fixed
* point  : 8 bytes fixed
* shn    : 1 byte length + data
* app    : 1 byte length + data
* marks  : 2 byte length + data
* tokens : 2 byte length + data
* line   : 2 byte length + data














abstract (4 sentences, 100 readers)
  here is a problem
    it's an interesting problem
    it's an unsolved problem
  here is my idea/claim/contribution
    your paper should have only 1 clear, sharp, specific, idea
    if you have lots of ideas, write lots of papers.
  my idea works (details, data)

introduction (1 page 100 readers)
  describe the problem with an example
  bullet list of the claims with forward references to the evidence

the problem (1 page 10 readers)
  explain the intuition first as if speaking at the whiteboard... examples are good.

my idea( 2 pages 10 readers)
  "the main idea of this paper is..."
  "in this section we present the main contributions of the paper."

the details( 5 pages, 3 readers)
  choose the most direct route to the idea
    don't spend time explaining an approach that did not work.. maybe in related work???
      don't force the reader down all the dead ends you had to go through to find the idea/solution
        your sweat, blood, and tears are not interesting to the reader
    but do note why very obvious solutions won't work to keep the reader from being distracted.
  give forward references to related work

related work(1-2 pages, 10 readers)
  the rest of the paper forms a lense (assumptions and definitions) for viewing the related work
  giving credit to others does NOT diminish the credit you get from your paper
  warmly acknowledge pwople who have helped you
  acknowledge weaknesses in your approach. better on x, but worse on y. explain the tradeoffs

conclusions and further work(0.5 pages)

