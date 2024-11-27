import json

with open(f"./haiku.json", "r") as f:  # データセットロード
    dataset = json.load(f)
output = []
for season_key, season_value in dataset.items():
    prompt = season_key + "の季語を使った俳句を詠んで"
    for kigo in season_value.values():
        for haiku in kigo:
            messages = []
            d_user = dict(content=prompt, role="user")
            d_assistant = dict(content=haiku, role="assistant")
            messages.append(d_user)
            messages.append(d_assistant)
            d_output = dict(prompt=prompt, messages=messages)
            output.append(d_output)

with open(f"./haiku_data.json", "w", encoding="UTF-8") as f:
    json.dump(output, f, ensure_ascii=False)
