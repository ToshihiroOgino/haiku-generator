from bs4 import BeautifulSoup
from resource import load_from_file, save_to_file
import re

_SEASONS = [
    "春",
    "夏",
    "秋",
    "冬",
    "新年",
]


def __extract_kigo_list(html_str: str):
    soup = BeautifulSoup(html_str, "html.parser")
    kigo_dict = dict()
    REGEXP = re.compile(r"kigo_work_list\.php\?kigo_cd=(\d+)$")
    for kigo in soup.find_all(
        "a", href=lambda href: href and "kigo_work_list.php?kigo_cd=" in href
    ):
        cd = int(REGEXP.match(kigo["href"]).group(1))
        kigo_dict.update({kigo.text: cd})
    return kigo_dict


def extract_kigo(verbose=False):
    data = dict()
    for season in _SEASONS:
        file_name = f"季語_{season}.html"
        if verbose:
            print(f"Loading {file_name}...")
        html_str = load_from_file(file_name)
        if verbose:
            print("Extracting kigo list...")
        kigo_dict = __extract_kigo_list(html_str)
        if verbose:
            print(f"Found {len(kigo_dict)} kigo.")
        data.update({season: kigo_dict})
    return data


async def __extract_haiku_list(season, kigo_cd):
    file_name = f"{season}/kigo_{kigo_cd}.html"
    html_str = load_from_file(file_name)
    soup = BeautifulSoup(html_str, "html.parser")
    haiku_list = list()
    for haiku in soup.find_all(
        "a", href=lambda href: href and "work_detail.php?cd=" in href
    ):
        haiku_list.append(haiku.text)
    return haiku_list


async def extract_haiku(verbose=False):
    if verbose:
        print("Extracting kigo list...")
    data = extract_kigo()
    if verbose:
        print("Extracting haiku list...")
    fut = dict()
    for season, kigo_dict in data.items():
        for _, cd in kigo_dict.items():
            fut.update({cd: __extract_haiku_list(season=season, kigo_cd=cd)})
    for season, kigo_dict in data.items():
        for kigo, cd in kigo_dict.items():
            data[season][kigo] = await fut[cd]
    return data


if __name__ == "__main__":
    import asyncio
    import json

    data = asyncio.run(extract_haiku(verbose=True))
    data_str = json.dumps(data, indent=2, ensure_ascii=False)
    save_to_file("haiku.json", data_str)
