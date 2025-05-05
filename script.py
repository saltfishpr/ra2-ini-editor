#!/usr/bin/env -S uv run --script --env-file .env
# /// script
# requires-python = ">=3.13"
# dependencies = ["click", "google-genai", "google-cloud-translate", "pydantic", "rich"]
# ///

import configparser
import json
import logging
import os
import sys

import click
from google import genai
from google.cloud import translate_v3
from rich.progress import Progress

logging.basicConfig(
    format="%(asctime)s - %(name)s - %(levelname)s - %(filename)s:%(lineno)d - %(message)s",  # noqa: E501
    level=logging.INFO,
    handlers=[logging.StreamHandler()],
)

logger = logging.getLogger(__name__)


PROJECT_ID = os.environ.get("GOOGLE_CLOUD_PROJECT")


def detect_language(text: str) -> str:
    client = translate_v3.TranslationServiceClient()
    parent = f"projects/{PROJECT_ID}/locations/global"

    # Detail on supported types can be found here:
    # https://cloud.google.com/translate/docs/supported-formats
    response = client.detect_language(
        content=text,
        parent=parent,
        mime_type="text/plain",
    )

    for language in response.languages:
        logger.debug(
            f"Language code: {language.language_code}, Confidence: {language.confidence}"
        )

    return response.languages[0].language_code


def translate_text(
    text: str,
    source_language_code: str = "en",
    target_language_code: str = "zh",
) -> str:
    client = translate_v3.TranslationServiceClient()
    parent = f"projects/{PROJECT_ID}/locations/global"

    # Translate text from English to chosen language
    # Supported mime types: # https://cloud.google.com/translate/docs/supported-formats
    response = client.translate_text(
        contents=[text],
        source_language_code=source_language_code,
        target_language_code=target_language_code,
        parent=parent,
        mime_type="text/plain",
    )

    # Display the translation for each input text provided
    for translation in response.translations:
        logger.debug(f"Translated text: {translation.translated_text}")

    return response.translations[0].translated_text


def read_csf(file_path: str) -> dict:
    def intb(b: bytes) -> int:
        return int.from_bytes(b, "little")

    with open(file_path, "rb") as f:
        csf = f.read()
    # cut CSF content to a list by " LBL"
    csf_list = csf.split(b" LBL")[1:]
    csf_dict = dict()
    # cut every label to key-value by " RTS"/"WRTS"
    for i in range(len(csf_list)):
        if b" RTS" in csf_list[i]:
            csf_list[i] = csf_list[i].split(b" RTS")
        else:
            csf_list[i] = csf_list[i].split(b"WRTS")
        csf_key = csf_list[i][0][8 : 8 + intb(csf_list[i][0][4:7])].decode("ASCII")
        csf_val_b = csf_list[i][1][4 : 4 + 2 * intb(csf_list[i][1][0:3])]
        csf_val = bytes([0xFF - b for b in csf_val_b]).decode("UTF-16-LE")
        csf_val = csf_val.replace("\n", "\\n")
        # store key-value as an dictionary
        csf_dict[csf_key] = csf_val
    return csf_dict


@click.group()
def cli():
    pass


@cli.command()
@click.argument("csf_path")
def csf_to_ini(csf_path: str) -> None:
    output_path = csf_path[:-4] + ".ini"  # .csf 替换为 .ini
    logger.info(f"csf_path: {csf_path}, output_path: {output_path}")

    csf_dict = read_csf(csf_path)

    config = configparser.ConfigParser(interpolation=None)
    config.optionxform = str
    config["zh-TW"] = csf_dict
    with open(output_path, "w", encoding="utf-8") as configfile:
        config.write(configfile)


@cli.command()
@click.argument("desc")
def gen_name(desc: str) -> str:
    client = genai.Client(api_key=os.environ.get("GEMINI_API_KEY"))
    response = client.models.generate_content(
        model="gemini-2.0-flash",
        contents=f"根据给定的配置描述，生成简短的配置名称。\n你只需给出配置名称，**不要**做其他回答。\n配置描述：{desc}",
    )
    return response.candidates[0].content.parts[0].text


def try_translate_text(
    text: str,
    source_language_code: str = "en",
    target_language_code: str = "zh",
) -> str:
    if text is None or text == "":
        return ""
    try:
        return translate_text(text, source_language_code, target_language_code)
    except Exception:
        logger.exception(f"translate {text} to {target_language_code}")
        return ""


def get_translated_keys(output_data: dict) -> list[str]:
    return [flag["key"] for flag in output_data.get("flags", []) if flag.get("desc")]


@cli.command()
@click.argument("input_filename")
@click.option(
    "--lang", "target_language_code", default="zh", help="目标语言代码 (默认: zh)"
)
def translate_schema(input_filename: str, target_language_code: str):
    if not input_filename.endswith(".json"):
        raise ValueError("input_filename must end with .json")

    output_filename = input_filename.replace(".json", f".{target_language_code}.json")

    translated_keys = []
    output_data = None

    if os.path.exists(output_filename):
        with open(output_filename, "r", encoding="utf-8") as out_file:
            output_data = json.load(out_file)
            translated_keys = get_translated_keys(output_data)

    with open(input_filename, "r", encoding="utf-8") as infile:
        data = json.load(infile)

    flags = data.get("flags", [])
    with Progress() as progress:
        task = progress.add_task("Translating...", total=len(flags))
        for flag in flags:
            if flag.get("key") in translated_keys:
                progress.advance(task)
                continue
            flag["desc"] = try_translate_text(
                flag.get("desc", ""), target_language_code=target_language_code
            )
            progress.advance(task)

    with open(output_filename, "w", encoding="utf-8") as outfile:
        json.dump(data, outfile, ensure_ascii=False, indent=4)


if __name__ == "__main__":
    sys.exit(cli())
