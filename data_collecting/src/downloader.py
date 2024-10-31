import requests
from resource import save_to_file
from parser import extract_kigo


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


def __download_haiku_from_kigo(kigo_cd, season):
    url = f"https://haiku-data.jp/kigo_work_list.php?kigo_cd={kigo_cd}"
    html = __get_html(url)
    save_to_file(f"{season}/kigo_{kigo_cd}.html", html)


def __download_haiku_page(verbose=False):
    kigo_data = extract_kigo()
    total_count = sum([len(kigo_dict) for _, kigo_dict in kigo_data.items()])
    if verbose:
        print(f"Downloading {total_count} kigo pages...")
    count = 0
    for season, kigo_dict in kigo_data.items():
        for kigo, kigo_cd in kigo_dict.items():
            if verbose:
                print(f"Downloading {kigo}({count}/{total_count})...")
            __download_haiku_from_kigo(kigo_cd, season)
            count += 1


if __name__ == "__main__":
    if input("Do you want to download the data? (y/n): ") != "y":
        exit()
    __download_kigo_list()
    __download_haiku_page()
