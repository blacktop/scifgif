# NOTES

## xkcd API

* https://github.com/nishanths/go-xkcd

## Giphy API

* https://github.com/peterhellberg/giphy
* https://github.com/sanzaru/go-giphy

## ascii-emoji

* https://github.com/dysfunc/ascii-emoji
* http://emojicons.com/popular

## expansion-packs

* https://github.com/medcl/elasticsearch-migration
* https://github.com/hoffoo/elasticsearch-dump
* https://github.com/olivere/elastic/blob/release-branch.v5/recipes/bulk_insert/bulk_insert.go

## Web Service

```
GET /xkcd/random
GET /xkcd/{search}
GET /giphy/random
GET /giphy/{search}
```

```
curl -XGET "http://elasticsearch:9200/_search" -H 'Content-Type: application/json' -d'
{
  "query": {
    "function_score": {
      "query": {
        "match": {
          "text": "sad"
        }
      },
      "boost": "5",
      "random_score": {},
      "boost_mode": "multiply"
    }
  }
}'
```

## Kibana Query

```
http://localhost:5601/app/kibana#/discover?_g=()&_a=(columns:!(emoji),filters:!(('$state':(store:appState),meta:(alias:!n,disabled:!f,index:AV7VwO5jCz-NMQd5a_vF,key:_type,negate:!f,type:phrase,value:ascii),query:(match:(_type:(query:ascii,type:phrase))))),index:AV7VyG7A8u-ErzWhyq_K,interval:auto,query:(query_string:(analyze_wildcard:!t,query:'double+flipping')),sort:!(_score,desc))
```

## Web UI

* https://codepen.io/simonswiss/pen/PNeJmy
* https://picsum.photos/
* https://codepen.io/philcheng/pen/YWyYwG

## ElasticDump 2 JSON

* https://github.com/medcl/esm
