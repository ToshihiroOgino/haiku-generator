import torch
from torch.utils.data import DataLoader
from transformers import AutoTokenizer, AutoModelForCausalLM, BitsAndBytesConfig
from torch.optim.lr_scheduler import ReduceLROnPlateau
from peft import get_peft_model, LoraConfig, TaskType
from datasets import load_dataset

# ローカルモデルのパス
LOCAL_MODEL_PATH = "./elyza_model"

# ハイパーパラメータ設定
BATCH_SIZE = 16
NUM_EPOCHS = 10
LEARNING_RATE = 1e-4
PATIENCE = 2  # Early stopping の許容エポック数

# モデルとトークナイザの読み込み
tokenizer = AutoTokenizer.from_pretrained(LOCAL_MODEL_PATH)
model = AutoModelForCausalLM.from_pretrained(LOCAL_MODEL_PATH, quantization_config=BitsAndBytesConfig(load_in_8bit=True), device_map="auto")

# LoRA の設定
lora_config = LoraConfig(
    task_type=TaskType.CAUSAL_LM,
    r=8,
    lora_alpha=32,
    target_modules=["q_proj", "v_proj"],
    lora_dropout=0.05,
    bias="none"
)
model = get_peft_model(model, lora_config)

# データセットのロード
dataset = load_dataset("text", data_files={"train": "haiku_dataset.txt"})  # テキスト形式のデータセット
train_data = dataset["train"]

# 最適化設定
optimizer = torch.optim.AdamW(model.parameters(), lr=LEARNING_RATE, weight_decay=0.01)

# 学習率スケジューラー（Lossが停滞した場合に学習率を減少）
scheduler = ReduceLROnPlateau(optimizer, mode='min', factor=0.1, patience=10, verbose=True)

# トークン化関数
def tokenize_function(examples):
    return tokenizer(
        examples["text"], 
        padding="max_length",
        truncation=True,
        max_length=128  # 必要に応じて調整
    )

tokenized_dataset = train_data.map(tokenize_function, batched=True)

# データローダーの作成
train_dataloader = torch.utils.data.DataLoader(
    tokenized_dataset,
    batch_size=BATCH_SIZE,
    shuffle=True
)

# Early Stoppingの初期化
best_loss = float("inf")
patience_counter = 0

# 学習設定
optimizer = torch.optim.AdamW(model.parameters(), lr=LEARNING_RATE)

device = torch.device("cuda")
model.to(device)
model.train()

# 学習ループ
for epoch in range(NUM_EPOCHS):
    model.train()  # モデルを訓練モードに
    epoch_loss = 0
    for step, batch in enumerate(train_dataloader):
        # 入力データの準備
        inputs = {
            "input_ids": torch.stack(batch["input_ids"]).to(device),
            "attention_mask": torch.stack(batch["attention_mask"]).to(device),
            "labels": torch.stack(batch["input_ids"]).to(device),
        }

        # モデルの順伝播と損失計算
        outputs = model(**inputs)
        loss = outputs.loss
        loss.backward()  # 勾配計算
        optimizer.step()  # パラメータ更新
        optimizer.zero_grad()  # 勾配のリセット
        
        # 1ステップごとにLossを表示
        print(f"Epoch {epoch}, Step {step}, Loss: {loss.item()}")
            
        # エポックごとの損失の集計
        epoch_loss += loss.item()

    epoch_loss /= len(train_dataloader)  # エポックごとの平均 Loss
    print(f"Epoch {epoch}, Loss: {epoch_loss}")

    # Loss が改善していない場合は Early Stopping
    if epoch_loss < best_loss:
        best_loss = epoch_loss
        patience_counter = 0
    else:
        patience_counter += 1
    
    # Early Stopping の判断
    if patience_counter >= PATIENCE:
        print("Early stopping triggered")
        break

    # 学習率スケジューラーで学習率を更新
    scheduler.step(epoch_loss)

    # モデルの保存
    model.save_pretrained("lora_finetuned_model_epoch"+str(epoch))
    tokenizer.save_pretrained("lora_finetuned_model_epoch"+str(epoch))
