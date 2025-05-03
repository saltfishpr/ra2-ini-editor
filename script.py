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
from pydantic import BaseModel
from rich.progress import (
    BarColumn,
    Progress,
    SpinnerColumn,
    TaskProgressColumn,
    TextColumn,
)

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


def first_index(text: str, sub: str) -> int:
    try:
        return text.index(sub)
    except ValueError:
        return -1


class Property(BaseModel):
    key: str
    name: str
    desc: dict[str, str]


@cli.command()
@click.argument("ini_path")
def gen_schema(ini_path: str) -> None:
    config = configparser.ConfigParser(interpolation=None)
    config.optionxform = str
    config.read(ini_path)

    output_path = ini_path[:-4] + ".json"  # .ini 替换为 .json
    logger.info(f"ini_path: {ini_path}, output_path: {output_path}")

    schema = {}
    if os.path.exists(output_path):
        with open(output_path, "r", encoding="utf-8") as f:
            schema = json.load(f)

    for section in config.sections():
        properties = []
        if section in schema:
            properties = [Property.model_validate(prop) for prop in schema[section]]

        items = list(config.items(section))
        with Progress(
            SpinnerColumn(),
            TextColumn("[progress.description]{task.description}"),
            BarColumn(),
            TaskProgressColumn(),
        ) as progress:
            task = progress.add_task(
                f"[cyan]Processing section {section}...", total=len(items)
            )
            for key, value in items:
                if any(prop.key == key for prop in properties):
                    progress.advance(task)
                    continue

                idx = first_index(value, ";")
                if idx == -1:
                    desc = value
                else:
                    desc = value[idx + 1 :]
                desc = desc.strip()
                try:
                    descZh = translate_text(desc)
                except Exception as e:
                    logger.error(f"Error translating text: {desc}, Error: {e}")
                i18nDesc = {
                    "en": desc,
                    "zh": descZh or "",
                }
                prop = Property(key=key, name=key, desc=i18nDesc)
                properties.append(prop)
                progress.advance(task)
        schema[section] = [prop.model_dump() for prop in properties]

    with open(output_path, "w", encoding="utf-8") as f:
        json.dump(schema, f, indent=4)


if __name__ == "__main__":
    sys.exit(cli())
