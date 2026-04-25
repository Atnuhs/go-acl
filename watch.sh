#!/usr/bin/env bash

BUNDLE="./submit/main.go"
BIN="/tmp/ac_watch_$$"
LAST_HASH=""

cleanup() {
    rm -f "$BIN"
    exit 0
}
trap cleanup INT TERM

try_compile() {
    local hash
    hash=$(md5sum main.go 2>/dev/null | cut -d' ' -f1)
    [[ "$hash" == "$LAST_HASH" ]] && return 0

    printf "\033[33m=== bundling & compiling... ===\033[0m\n"
    local err
    if err=$(go-bundler -dir . -with-metrics -with-sustainability-metrics > "$BUNDLE" 2>&1 \
              && go build -o "$BIN" "$BUNDLE" 2>&1); then
        LAST_HASH="$hash"
        printf "\033[32m=== ready ===\033[0m\n\n"
        return 0
    else
        printf "\033[31m=== compile error ===\033[0m\n%s\n\n" "$err"
        rm -f "$BIN"
        return 1
    fi
}

try_compile

while true; do
    # コンパイル失敗中は1秒ポーリングして再コンパイル待ち
    if [[ ! -x "$BIN" ]]; then
        sleep 1
        try_compile
        continue
    fi

    # 入力前にファイル変更チェック（instant）
    try_compile

    printf '%s\n' "--- input ---"
    output=$("$BIN") || true

    printf '%s\n' "" "*--- ANS ---*"
    printf '%s\n' "$output"
    printf '\n'
done

cleanup
