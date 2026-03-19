#!/usr/bin/env python3
"""Generate a 32-byte TOKEN_ENC_KEY_B64 file for local development."""

import argparse
import base64
import os
import secrets
import stat
from pathlib import Path

DEFAULT_PATH = Path("./secrets/token_enc_key.b64")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Generate TOKEN_ENC_KEY_B64 file")
    parser.add_argument("--path", type=Path, default=DEFAULT_PATH, help="Output key file path")
    parser.add_argument("--force", action="store_true", help="Overwrite existing file")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    out_path = args.path

    out_path.parent.mkdir(parents=True, exist_ok=True)

    if out_path.exists() and not args.force:
        print(f"Refusing to overwrite existing file: {out_path}. Use --force to replace it.")
        return 1

    raw_key = secrets.token_bytes(32)
    encoded = base64.b64encode(raw_key).decode("ascii")
    out_path.write_text(encoded + "\n", encoding="utf-8")

    os.chmod(out_path, stat.S_IRUSR | stat.S_IWUSR)
    print(f"Wrote TOKEN_ENC_KEY_B64 to {out_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
