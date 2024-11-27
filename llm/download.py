from transformers import AutoModelForCausalLM, AutoTokenizer

# モデルをローカルに保存
MODEL_NAME = "elyza/ELYZA-japanese-Llama-2-7b-fast"
model = AutoModelForCausalLM.from_pretrained(MODEL_NAME)
tokenizer = AutoTokenizer.from_pretrained(MODEL_NAME)

# 保存するパス
LOCAL_DIR = "./elyza_model"
model.save_pretrained(LOCAL_DIR)
tokenizer.save_pretrained(LOCAL_DIR)
