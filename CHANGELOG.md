# Changelog

## 0.1.2
* Added ability to specify multiple topics by seperating them with a comma
* Added ability to specify custom consumer group name as well as add a randomly generated string to the consumer group when it is automatically generated.  This allows for multiple groups to consumer each their own copy of the data if necessary.
* Added two counters (`bytesConsumed` and `bytesDisaplyed`) to track how many bytes were consumed from the selected topics as well has how many bytes were then filted (by the regex) to be displayed.

## 0.1.1
* Update README and Makefile
* Updated flag help messages
* Updated version package

## 0.1.0
* Initial release