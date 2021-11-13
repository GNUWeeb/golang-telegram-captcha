import urllib.request
import json
import sys

emojis = "🍜🥘🧆🍤🦪🍭🍦🍯🥜🍩🍿🍷🧉🍽🥄🍫🍬🧅🍖🍗🍕🌮🔥🌈🐶🐒🦉🥾💍🌂🚚🚜🛵🧨🧲🔮🎉"

def download_emoji():
    for emoji in emojis:
        code = 'u{:x}'.format(ord(emoji))
        
        emoji_img_url = f"https://github.com/googlefonts/noto-emoji/blob/main/png/72/emoji_{code}.png?raw=true"
        urllib.request.urlretrieve(emoji_img_url, f"./assets/image/emoji/{code}.png")
        print(emoji, " downloaded...")

def emoji_map():
    emoji_map = {}
    for emoji in emojis:
        code = 'u{:x}'.format(ord(emoji))
        emoji_map[code] = emoji
    print(json.dumps(emoji_map, ensure_ascii=False))

if __name__ == '__main__':
    globals()[sys.argv[1]]()