
Goals
=====

* store 7 years of normalized events at ~ 1TB of events per day.
* find events containing terms
* shard with separate instances of the whole service: per: organzation, environment, location etc...


Analysis
========
* Survey of the field http://www.kuehnel.org/bachelor.pdf
* pathologies of big data http://queue.acm.org/detail.cfm?id=156387
* comparison of golang embedded b+ indexes ????
* comparison of golang embedded LSM indexes ????


Terminology
===========

term:
* case insensitive unicode string without whitespace or punctuation
* the human input to the search process

token:
* key_name + value_token

key_name:
* case sensitive and preserving
* strip non-printable characters
* strip whitespace

value_token:
* case sensitive and preserving
* strip non-printable characters
* strip whitespace
* contains some punctuation
* limited in length

day_stamp:
* seconds since epoch at the start of the day containing the event 00:00:00.0000000
* 4 bytes, uint32 each rolls over in 2106
* converted to base 10 millis for javascript

event_id:
* daystamp
* offset in that day's event store


7 year Indexes
==============

* term -> []token: stored as space separated strings
* token -> []hour: 6byte hour = 4byte epoch + 2byte count

hours in 7 years = 61320


Hour Event Store
===============

* token -> []event_id  fixed-width x-byte event_id's ??????
* token table:null terminated, frequency sorted list of tokens
* pile checksum

Binary pile of compressed events:

* operations: append(), get(offset), get_all()
* can detect all corruption and mark any corrupt events in search results
* event_id is offset in file

Search Process
==============

step 1: exclude irrelevant hours

List hours in which all tokens were found at least once in any message during the hour.

* text input gets a good list of terms
* submit terms and return tokens
* select and submit tokens
* return hours containing ALL tokens
* histogram by week or month

step 2: detailed search over a small time range

* pull the per-hour data from object storage
* load events into indexes capable of compound queries
* query the per-hour indexes


Hour Prune Process
=================

* pick hour to prune and get hour_stamp
* for each combo in the combo -> []event_id index

```
find token in token -> []hour index
    remove hour
    if empty, delete token:
       find all terms containing token????
       remove token from each
       if last entry, delete term
```

* delete per-hour folder/file??? from S3 or object storage


7 Year Index Backup and Restore
===============================

* boltdb snapshots or something?


S3 partial download
===================

* http://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectGET.html
* http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.35



Old Attempts
============

Try 1: single data structure

* token -> []event_id
* return event_id's that appear in all sets with a single lookup

problems:

* single massive, constantly updating token -> []event_id table
* pruning by date requires table scan and rewrite


Try 2: map reduce

```
for each day:
  token -> []day_event_id
merge(all day_event_id sets)
7 year query requires 2555 index lookups and set operations
```
