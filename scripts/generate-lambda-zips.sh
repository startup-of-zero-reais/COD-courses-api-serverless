

#!bin/bash

TERRAFORM_LAMBDAS=terraform/lambda/functions

if [ ! -d "$TERRAFORM_LAMBDAS" ]; then
  echo "O diretório Terraform lambda functions não existe"
  mkdir -p $TERRAFORM_LAMBDAS
fi

function isEmpty {
    local LAMBDA=$1
    local IS_EMPTY=$(find ./src/"$LAMBDA" -maxdepth 0 -empty)

    if [ "$IS_EMPTY" != "" ]; then
        echo "true"
    else
        echo "false"
    fi    
}

function buildGoLambda {
    local LAMBDA=$1

    printf "Buildando Go lambda para %s\n" "$LAMBDA"
    cd ./src/"$LAMBDA"

    if [ -f "$LAMBDA" ]; then
        rm "$LAMBDA"
    fi

    GOOS=linux CGO_ENABLED=0 go build -o "$LAMBDA" .
    cd ../../

    return 0
}

function generateZip {
    local LAMBDA=$1

    if [ -f "./src/$LAMBDA/$LAMBDA" ]; then
        printf "Gerando zip de %s\n" "$LAMBDA"

        if [ -f "$TERRAFORM_LAMBDAS/$LAMBDA.zip" ]; then
            rm "$TERRAFORM_LAMBDAS/$LAMBDA.zip"
        fi

        cd src/"$LAMBDA"
        zip -r ../../"$TERRAFORM_LAMBDAS/$LAMBDA.zip" "$LAMBDA"
        cd ../../
    else
        printf "Pulando %s\n" "$LAMBDA"
    fi

    return 0
}

LAMBDAS_SOURCE_DIR=$(ls -d src/*/ | awk -F '/' '{print $2}')

for LAMBDA in $LAMBDAS_SOURCE_DIR
do
    if [ "$(isEmpty $LAMBDA)" == "true" ]; then
        echo "Lambda $LAMBDA está vazio"
        continue
    fi

    buildGoLambda $LAMBDA
    sleep 1
    generateZip $LAMBDA
done