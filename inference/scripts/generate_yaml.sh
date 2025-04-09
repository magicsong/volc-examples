#!/bin/bash

# 设置默认值
ENV_FILE="${ENV_FILE:-./.env}"
TEMPLATE_FILE="${TEMPLATE_FILE:-./all-in-one.yaml}"
OUTPUT_FILE="${OUTPUT_FILE:-./final-all-in-one.yaml}"

# 加载 .env 文件
if [ -f "$ENV_FILE" ]; then
  export $(grep -v '^#' $ENV_FILE | xargs)
else
  echo ".env file not found at $ENV_FILE"
  exit 1
fi

set -e
# 替换模板中的变量
sed -e "s|\${NAMESPACE}|$NAMESPACE|g" \
    -e "s|\${IMAGE}|$IMAGE|g" \
    -e "s|\${GPU}|$GPU|g" \
    -e "s|\${BUCKET}|$BUCKET|g" \
    -e "s|\${MODEL_PATH_TOS}|$MODEL_PATH_TOS|g" \
    -e "s|\${URL}|$URL|g" \
    $TEMPLATE_FILE > $OUTPUT_FILE

echo "Final YAML generated at: $OUTPUT_FILE"
