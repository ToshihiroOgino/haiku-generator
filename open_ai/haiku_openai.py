# -*- coding: utf-8 -*-

import openai
import os
import sys
import time

# APIキーを環境変数から取得
openai.api_key = os.getenv("API_KEY")


def generate_haiku(prompt, messages, retries=3, retry_delay=20):
    for attempt in range(retries):
        try:
            response = openai.ChatCompletion.create(
                model="gpt-4o-mini",
                messages=messages + [{"role": "user", "content": prompt}],
                temperature=0.7,
            )
            return response.choices[0].message["content"].strip()
        except openai.error.RateLimitError:
            if attempt < retries - 1:
                print(f"Rate limit reached. Retrying in {retry_delay} seconds...")
                time.sleep(retry_delay)  # リトライの前に待機
            else:
                print("Max retries reached. Please try again later.")
                raise


def prompt_user(message):
    user_input = input(f"{message} \n ")
    if user_input.lower() == "終了" or user_input.lower() == "end":
        print("生成を終了します．")
        sys.exit()
    return user_input


def read_file(file_path):
    """指定されたテキストファイルを読み込む関数"""
    with open(file_path, "r", encoding="utf-8") as file:
        return file.read()


def split_text_into_chunks(text, chunk_size=1000):
    """テキストを指定したサイズでチャンクに分割する"""
    return [text[i : i + chunk_size] for i in range(0, len(text), chunk_size)]


def save_haiku_to_file(conversation, filename):
    with open(filename, "a", encoding="utf-8") as file:
        file.write(conversation)


# メイン関数
def main():
    # チャットの履歴を初期化
    messages = [{"role": "system", "content": "あなたは、有名な俳人です。"}]

    # ファイルのパスを指定
    file_path = "haiku_random.txt"  # 俳句データ（全28143句） からランダムに選んだ1000句の季語と俳句のペア

    # テキストファイルの内容を読み込む
    file_content = read_file(file_path)

    # テキストをチャンクに分割
    chunked_content = split_text_into_chunks(file_content, chunk_size=1000)

    while True:
        kigo = prompt_user("季語を入力してください")
        # チャンクごとにファイル内容をプロンプトに組み込む
        for chunk in chunked_content:
            haiku = generate_haiku(
                f"""ユーザは季語：{kigo}を入力しました．
                                       575の俳句を10つ生成してください．
                                       参考データ：{chunk}
                                    """,
                messages,
            )
            haiku += "\n"

        print(haiku)
        save_haiku_to_file(haiku, "generated_haiku.txt")
        print("保存されました")


if __name__ == "__main__":
    main()
