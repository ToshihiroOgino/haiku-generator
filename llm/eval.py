from transformers import AutoTokenizer, AutoModelForCausalLM
from peft import PeftModel
import torch

# ベースモデルとLoRAモデルのパス
BASE_MODEL_PATH = "elyza_model"
LORA_MODEL_PATH = "lora_finetuned_model_epoch3"

# デバイスの設定
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

# トークナイザーとモデルのロード
tokenizer = AutoTokenizer.from_pretrained(BASE_MODEL_PATH)  # トークナイザーはベースモデルに依存
tokenizer.padding_side = 'left'  # 左側にパディングを設定

base_model = AutoModelForCausalLM.from_pretrained(BASE_MODEL_PATH)
model = PeftModel.from_pretrained(base_model, LORA_MODEL_PATH)

# モデルをデバイスに転送
model.to(device)
model.eval()

# 入力例
input_ = [
    "夏 | 五月 | ",
    "秋 | 八月 | ",
    "冬 | 十二月 | ",
    "新年 | 去年今年 | ",
    "春 | 啓蟄 | ",
    "冬 | 大寒 | ",
    "夏 | 早乙女 | ",
    "夏 | 汗 | ",
    "夏 | 炎天 | ",
    "秋 | 秋の暮 | ",
    "春 | 花 | ",
    "秋 | 花野 | ",
    "春 | 苗代 | ",
    "秋 | 菊 | ",
    "秋 | 萩 | ",
    "春 | 藤 | ",
    "夏 | 虹 | ",
    "夏 | 蟻 | ",
    "秋 | 露 | ",
    "秋 | 黄落 | ",
]
input_texts = [(i.split(" | ")[0], i.split(" | ")[1]) for i in input_]

# バッチサイズの設定
batch_size = 2
max_length = 50

# バッチ処理
for i in range(0, len(input_texts), batch_size):
    # 入力データをバッチに分割
    batch_data = input_texts[i:i + batch_size]

    # バッチごとにプロンプトを作成
    prompts = []
    for season, keyword in batch_data:
        prompt = f"季節: {season}, 季語: {keyword}. この情報を元に俳句を生成してください:"
        prompts.append(prompt)

    # トークナイズしてテンソルを作成
    inputs = tokenizer(prompts, return_tensors="pt", padding=True, truncation=True).to(device)

    # モデルで生成
    with torch.no_grad():
        outputs = model.generate(
            **inputs,
            max_length=max_length,
            num_return_sequences=1,
            temperature=0.8,
            top_p=0.95,
            repetition_penalty=1.2,  # 繰り返しを抑制
            pad_token_id=tokenizer.pad_token_id,
            eos_token_id=tokenizer.eos_token_id,
        )

    # 結果を表示
    for idx, output in enumerate(outputs):
        season, keyword = batch_data[idx]
        generated_text = tokenizer.decode(output, skip_special_tokens=True)
        # プロンプトを取り除く
        haiku = generated_text.replace(prompts[idx], "").strip()
        print(f"季節: {season}, 季語: {keyword}")
        print(f"生成された俳句: {haiku}")
        print("-" * 50)
        
        r = open("result.txt", "a" ,encoding="utf-8")
        r.write(f"季節: {season}, 季語: {keyword}\n")
        r.write(f"生成された俳句: {haiku}\n")
        r.write("-" * 50 + "\n")
        r.close()
