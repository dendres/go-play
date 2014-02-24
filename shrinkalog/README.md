
pre-compress events for long term storage
=========================================

Pre-compressing large sets of log events saves a non-trivial amount of space and can be done efficiently.
Applying word-scale pre-compression to a set of N log events reduces the storage space by ???N??some???formula??? and is highly applicable to to the long term storage of large numbers of events.

Given set of N log events in a single file with size S, S can be reduced by ???N??some???formula??? by making a frequency sorted list of most frequently occurring "words" in the event's utf-8 strings and then replacing the words with an easily detectable, fixed size token representing the offset in the list. As a side effect, the token list is a kind of searchable index that can answer the question, did this exact "word" occurr within a given time range. The token list itself can then be further processed via downcase, stemming, and even 3gram or soundex to allow natural language sorting.


abstract (4 sentences, 100 readers)
  here is a problem
    it's an interesting problem
    it's an unsolved problem
  here is my idea/claim/contribution
    your paper should have only 1 clear, sharp, specific, idea
    if you have lots of ideas, write lots of papers.
  my idea works (details, data)



XXX have to make sure utf-8 fields do not have to be discovered
separate utf-8 from varint encoded numbers?????
avoid tokenizing large sets of numbers??? why?? overhead is small. search works better, and for long term compression single instance words can be filtered out??????

each substitution must be:
bitset for every utf-8 string
vs.
non-utf8 byte + varint
???????????????????????


Serialized Event:
* version: 1 byte fixed
* point  : 8 bytes fixed
* crc    : 4 bytes fixed
* shn    : 1 byte length + data
* app    : 1 byte length + data
* marks  : 2 byte length + data
* tokens : 2 byte length + data
* line   : 2 byte length + data


