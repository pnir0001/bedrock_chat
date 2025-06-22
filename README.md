# bedrock_chat
bedrock chat test repository

## description
This is the test repository for bedrock chat.

## usage


```
go mod tidy
go run main.go

curl --location 'http://localhost:8080/chat' \
--header 'Content-Type: application/json' \
--data '{
    "message": "テストです。簡単な応答を返却してください。あなたの自己紹介を簡単に行ってください。model名も添えてください。"
}'
```
