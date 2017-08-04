NOTES
=====

xkcd API
--------

-	https://github.com/nishanths/go-xkcd

Giphy API
---------

-	https://github.com/peterhellberg/giphy
-	https://github.com/sanzaru/go-giphy

ascii-emoji
-----------

-	https://github.com/dysfunc/ascii-emoji
-	http://emojicons.com/popular

expansion-packs
---------------

- https://github.com/medcl/elasticsearch-migration
- https://github.com/hoffoo/elasticsearch-dump
- https://github.com/olivere/elastic/blob/release-branch.v5/recipes/bulk_insert/bulk_insert.go

Web Service
-----------

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
