import json

# 入力データ
r=open("haiku.json", "r", encoding="utf-8")
raw_data = json.loads(r.read())
r.close()

# フラットな構造に変換
formatted_data = []
for season, kigo_dict in raw_data.items():
    for kigo, haikus in kigo_dict.items():
        for haiku in haikus:
            formatted_data.append({"season": season, "kigo": kigo, "haiku": haiku})

# 結果を保存
with open("formatted_haiku_data.json", "w", encoding="utf-8") as f:
    json.dump(formatted_data, f, ensure_ascii=False, indent=2)

# JSON ファイルをテキスト形式に変換
with open("formatted_haiku_data.json", "r", encoding="utf-8") as f:
    data = json.load(f)

with open("haiku_dataset.txt", "w", encoding="utf-8") as f:
    for entry in data:
        line = f"{entry['season']} | {entry['kigo']} | {entry['haiku']}\n"
        f.write(line)
