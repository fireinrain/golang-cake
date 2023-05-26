# openai account
# xxxxxxx@gmail.com
# shellcheck disable=SC2164
cd /root/docker-compose/chatgpt-web

refresh_token=""

# shellcheck disable=SC2016
response=$(curl -X POST -H "Content-Type: application/json" -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36" \
-d '{
    "redirect_uri": "com.openai.chat://auth0.openai.com/ios/com.openai.chat/callback",
    "grant_type": "refresh_token",
    "client_id": "pdlLIX2Y72MIl2rhLhTE9VV9bN905kBh",
    "refresh_token": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}' https://auth0.openai.com/oauth/token)

access_token=$(echo "$response" | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')

echo $access_token

query=$(echo "$access_token")
#echo $query

compose_file_path="docker-compose-chatgpt.yml"

sed -i "s/OPENAI_ACCESS_TOKEN: .*/OPENAI_ACCESS_TOKEN: $query/" "$compose_file_path" && \
docker-compose -f "$compose_file_path" down -v && docker-compose -f "$compose_file_path" up -d

exit 0




