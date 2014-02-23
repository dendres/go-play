
Incident Response
=================

What are the symptoms?  What are the causes?  Triage/Escalate/Patch

* got a customer complaint about service X. get events from service X in that time range in that cluster
* service X got an error from service Y. get events from service Y in that time range in that cluster

Debugging
=========

* I'm going to run a test against a load balanced and tiered application. let me watch all events on the cluster as the test unfolds.


Alerts
======

* Alert when event[key] == exact_value
* alert when event[key] =~ /regex/
* alert when exact_word in event[message]

Metrics
=======

* count event when 


