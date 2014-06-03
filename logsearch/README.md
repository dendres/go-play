

Analysis to complete
====================
* ????


pathologies of big data
http://queue.acm.org/detail.cfm?id=156387

Survey of the field:
http://www.kuehnel.org/bachelor.pdf
* rsyslog RELP transport???


Goals
=====

* long-term offline storage of normalized events
* case-independent exact-term search
* file-based indexes
* file-based results cache
* sharding = separate instances of the whole service.. probably per-environment or per-service

Interface
=========

* locate(downcased_term_list, time_range) returns term -> token -> day_stamp -> link

```
{
"timeout" : [
    "TimeOut" : {
        1404345600 : [
            link,
            link,
        ]
    },
    "Timeout" : {
        1403395200 : [
            link,
            link,
        ]
    }
]
}
```

* get([]links) returns link -> event

Index Files
===========

* for a given time period of long term storage, index files are assumed to be built once and never modified.
* index: term -> token
* index: token -> link

Token Substitution
==================

* a table of the 256 most frequent tokens
* created from the token index
* token(0..255) returns the token byte slice
* xxxx([]byte) returns 0...255

Event Store and Linking
=======================

* Store events for a given day in one big file
* have a link syntax that allows for minimal disk seeks to retrieve a given event
* find a way to allow the link syntax to roughly tell at least the hour of the day when the event took place?????
* must be able to detect corruption of any byte in the file
    checksum of the whole file after written
    compression checksum of each event


* could do CDB where time + checksum -> token substituted, compressed event
* listB?
* pick kv??


S3 partial download
===================

* http://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectGET.html
* http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.35

