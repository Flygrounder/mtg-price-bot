from json import dumps

import requests
from lxml import etree

URL_TEMPLATE = "https://magic.wizards.com/{}/articles/archive/card-image-gallery/{}"
OUTPUT_FILE_TEMPLATE = "{}.json"


def get_card_names(language, set_name):
    spoiler_url = URL_TEMPLATE.format(language, set_name)
    response = requests.get(spoiler_url)
    dom = etree.HTML(response.content.decode())
    card_names = dom.xpath('//div[@class="resizing-cig"]//p/text()')
    return [str(name).strip() for name in card_names]


def match_names(keys, values):
    return dict(zip(keys, values))


set_name = input("Введите сет: ")

russian_names = get_card_names("ru", set_name)
english_names = get_card_names("en", set_name)

match = match_names(russian_names, english_names)
print(match)

with open(OUTPUT_FILE_TEMPLATE.format(set_name), 'w') as output:
    output.write(dumps(match))
