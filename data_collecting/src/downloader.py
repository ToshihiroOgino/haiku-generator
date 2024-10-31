import requests
from resource import save_to_file


def __get_html(url: str) -> str:
    response = requests.get(url)
    response.encoding = response.apparent_encoding
    return response.text


_SEASON_IDS = {
    "春": 1,
    "夏": 2,
    "秋": 3,
    "冬": 4,
    "新年": 5,
}


def __download_kigo_list():
    for season, season_id in _SEASON_IDS.items():
        url = f"https://haiku-data.jp/kigo_list.php?season_cd={season_id}"
        html = __get_html(url)
        save_to_file(f"季語_{season}.html", html)


__download_kigo_list()
