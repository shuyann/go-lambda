### go lambda function example

```
cd xxx
GOOS=linux go build .
zip handler.zip ./xxx
aws lambda create-function \
  --region region \
  --function-name xxx \
  --memory 128 \
  --role arn:aws:iam::account-id:role/execution_role \
  --runtime go1.x \
  --zip-file fileb://path-to-your-zip-file/handler.zip \
  --handler xxx
```
