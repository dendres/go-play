
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

Did a set of conditions occurr within a given time frame?
for example, combinations of any of the following:
* event[key] == exact_value
* event[key] =~ /regex/
* exact_word found in event[message]

Metrics
=======

* calculate event rate during a time interval


