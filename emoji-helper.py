import urllib.request
import json
import sys

emojis = "๐๐ฅ๐ง๐ค๐ฆช๐ญ๐ฆ๐ฏ๐ฅ๐ฉ๐ฟ๐ท๐ง๐ฝ๐ฅ๐ซ๐ฌ๐ง๐๐๐๐ฎ๐ฅ๐๐ถ๐๐ฆ๐ฅพ๐๐๐๐๐ต๐งจ๐งฒ๐ฎ๐"

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