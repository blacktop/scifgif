import io
import json
from itertools import izip_longest

output = []
index = 1


def clean(word):
    return word.lower().rstrip().replace(' / ', " ")


with open("emoji.txt", 'r') as f:
    for emoji, keywords, blank in izip_longest(* [f] * 3):
        output.append(dict(id=str(index), emoji=clean(emoji), keywords=clean(keywords)))
        index += 1

with io.open('emoji.json', 'w', encoding='utf8') as json_file:
    data = json.dumps(output, ensure_ascii=False, encoding='utf8', indent=True)
    json_file.write(unicode(data))
