
Masterless is more stable and just as fast.
===========================================

A log event processing system with an elected master does not significantly improve delivery efficiency or reliability over a masterless system where producers choose a random processing server and processing servers choose random replication partners. allow 2 server failures in a 3 to 6 server system. replication factor 2. message exists on 3 random servers. any 2 can die. compare this to the kafka case with replication factor 2 in 3,4,5,6 server systems.



requirement:
 - throughput/efficiency should approach kafka for 10 second batching when sorted time output is irrelevant
 - for ~15min batching, with various hold down timers, throughput must also match kafka, but provide sorted output over that time interval
 - hold_down timer from 0 seconds to 48 hours
 - batch size from 10 to 1024 seconds
XXXX so this implies that either the file structure should converge to kafka's for small hold_down and batch_size
     OR... that the cost of the file structure is always the same as kafka's independent of hold_down and batch_size
   must take a much closer look at kafka's disk handling to determine!!!!

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
