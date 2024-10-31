from bs4 import BeautifulSoup
from resource import load_from_file
import re

_SEASONS = [
    "春",
    "夏",
    "秋",
    "冬",
    "新年",
]


def extract_kigo_list(html_str: str):
    soup = BeautifulSoup(html_str, "html.parser")
    kigo_dict = dict()
    REGEXP = re.compile(r"kigo_work_list\.php\?kigo_cd=(\d+)$")
    for kigo in soup.find_all(
        "a", href=lambda href: href and "kigo_work_list.php?kigo_cd=" in href
    ):
        cd = int(REGEXP.match(kigo["href"]).group(1))
        kigo_dict.update({kigo.text: cd})
    return kigo_dict


if __name__ == "__main__":
    for season in _SEASONS:
        file_name = f"季語_{season}.html"
        print(f"Loading {file_name}...")
        html_str = load_from_file(file_name)
        print("Extracting kigo list...")
        kigo_dict = extract_kigo_list(html_str)
        print(f"Found {len(kigo_dict)} kigo.")
