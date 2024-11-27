import torch
print(torch.cuda.is_available())  # True なら GPU が認識されています
print(torch.cuda.device_count())  # 利用可能な GPU の数
print(torch.cuda.get_device_name(0))  # 利用可能な GPU の名前
