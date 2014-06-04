
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

combo:
* sorted token combination
* space separated
* 1 to 4 tokens

token:
* case sensitive and preserving
* key_name:value_token

key_name:
* strip non-printable characters
* strip whitespace

value_token:
* case sensitive and preserving
* contains punctuation
* no whitespace
* only printable characters
* extracted from fields by white space only

day_stamp:
* seconds since epoch at the start of the day containing the event 00:00:00.0000000

event_id:
* daystamp
* offset in that day's event store


7 year Indexes
==============

These should be read-optimized b+ trees:

* term -> token
* token -> []combo
* combo -> []day_stamp


Day Event Store
===============

* combo -> []event_id is a write-optimized LSM or CDB index
* token table:null terminated, frequency sorted list of tokens
* pile checksum

Binary pile of compressed events:

* operations: append(), get(offset), get_all()
* can detect all corruption and mark any corrupt events in search results
* event_id is offset in file

Search Process
==============

* autocomplete string to terms conversion in browser
* submit terms and return tokens
* select and submit 1-4 tokens, get back list of combos
* submit combo, return list of days
* submit day,combo, return event_id's
* submit event_id, return event


Day Prune Process
=================

* pick day to prune and get day_stamp
* for each combo in the combo -> []event_id index

```
find combo in combo -> []day_stamp index
    remove day_stamp
    if empty, delete combo
    else write new value for combo
    if empty, delete combo from term -> []combo index:
        extract terms from combo
        find all terms
        remove combo from any terms that contain the combo
        write new value for term
        if term is empty, delete term
```

* delete per-day folder from S3 or object storage


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
