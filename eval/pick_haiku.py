import json
import random

ningen=[]
openai=[]
markov=[]
list_kigo=["五月","八月","十二月","去年今年","啓蟄","大寒","早乙女","汗","炎天","秋の暮","花","花野","苗代","菊","萩","藤","虹","蟻","露","黄落"]

with open(f'./haiku/haiku.json')as f:
    ningen=json.load(f)
ningen_haiku=[]
for season in ningen.values():
    for kigo_key, kigo_value in season.items():
        for i in range(10):
            if kigo_key==list_kigo[i]:
                ningen_haiku.append(kigo_value)


with open(f'./haiku/haiku_openai.json')as f:
    openai=json.load(f)

openai_haiku=[]
for index ,kigo in enumerate(openai['Decent'].values()):
    if index<10:
        openai_haiku.append(kigo)


with open(f'./haiku/generated.json')as f:
    markov=json.load(f)

markov_haiku=[]
for index ,kigo in enumerate(markov['Decent'].values()):
    if index<10:
        markov_haiku.append(kigo)

#乱数
random_numbers = random.choices(range(10), k=5)

output=[]
for i in range(5):
        d_output=dict(haiku=ningen_haiku[i][random_numbers[i]],yomite="ningen")
        output.append(d_output)
        d_output=dict(haiku=openai_haiku[i][random_numbers[i]],yomite="openai")
        output.append(d_output)
        d_output=dict(haiku=markov_haiku[i][random_numbers[i]],yomite="markov")
        output.append(d_output)

random.shuffle(output)

with open(f'./haiku/picked_haiku.json','w',encoding='UTF-8')as f:
    json.dump(output,f,ensure_ascii=False)
