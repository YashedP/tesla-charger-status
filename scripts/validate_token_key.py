#!/usr/bin/env python3
"""Validate TOKEN_ENC_KEY_B64 file format and length."""

import argparse
import base64
import binascii
from pathlib import Path

DEFAULT_PATH = Path("./secrets/token_enc_key.b64")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Validate TOKEN_ENC_KEY_B64 file")
    parser.add_argument("--path", type=Path, default=DEFAULT_PATH, help="Key file path")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    key_path = args.path

    if not key_path.exists():
        print(f"Key file not found: {key_path}")
        return 1

    content = key_path.read_text(encoding="utf-8").strip()
    if not content:
        print("Key file is empty")
        return 1

    try:
        raw = base64.b64decode(content, validate=True)
    except binascii.Error as exc:
        print(f"Invalid base64 key: {exc}")
        return 1

    if len(raw) != 32:
        print(f"Decoded key length must be 32 bytes, got {len(raw)}")
        return 1

    print(f"Valid TOKEN_ENC_KEY_B64 file: {key_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
