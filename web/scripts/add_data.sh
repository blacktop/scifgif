#!/bin/bash

curl -XPUT 'localhost:9200/shakespeare?pretty' -H 'Content-Type: application/json' -d'
{
 "mappings" : {
  "_default_" : {
   "properties" : {
    "speaker" : {"type": "keyword" },
    "play_name" : {"type": "keyword" },
    "line_id" : { "type" : "integer" },
    "speech_number" : { "type" : "integer" }
   }
  }
 }
}
'

curl -H 'Content-Type: application/x-ndjson' -XPOST 'localhost:9200/shakespeare/_bulk?pretty' --data-binary @shakespeare.json

curl -XGET 'localhost:9200/_cat/indices?v&pretty'
