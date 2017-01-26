# slack-count-history

counts the number of messages per channel

## Install

```
$ go get github.com/ktateish/slack-count-history
```

## Usage

```
$ export SLACK_API_TOKEN=xxx-set-your-token-here
$ slack-count-history
Connected to your-team as your-user
There are 8 channels
(1/8) aaabbb:  3
(2/8) bbbbbb:  10
(3/8) bbbccc:  55
(4/8) dccdd:  3
(5/8) general:  12
(6/8) random:  2
(7/8) tototo:  6
(8/8) xxxxxx:  3

    94 TOTAL
    55 #bbbccc
    12 #general
    10 #bbbbbb
     6 #tototo
     3 #aaabbb
     3 #dccdd
     3 #xxxxxx
     2 #random
```

* You can use -i flag to specify interval for api call in seconds

## Author

Katsuyuki Tateishi <kt@wheel.jp>
