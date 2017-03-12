#!/bin/sh

set -x

source cf-event/ENV.sh

cf login $CF_TARGET_SKIP_SSL_VALIDATION -a $CF_TARGET_API \
    -u $CF_TARGET_USER \
    -p $CF_TARGET_PASSWORD \
    -o $CF_TARGET_ORG \
    -s $CF_TARGET_SPACE

cf save-target -f copy-dest

APPS_TO_COPY=""
for a in $(echo $CF_APPS); do 
    GUID=$(cf app $a --guid)
    if [[ $? -eq 0 ]]; then
        eco "Checking if app '$a' at target is more recent than event received from same app at source."
        
        TARGET_TS=$(date -d"$(cf curl /v2/apps/$GUID | jq .metadata.created_at | sed 's|"||g')" +%s)
        SOURCE_TS=$(cat cf-event/$a.timestamp)
        echo "Target timestamp: $(date -d @$TARGET_TS)"
        echo "Source timestamp: $(date -d @$SOURCE_TS)"

        if [[ $TARGET_TS -gt $SOURCE_TS ]]; then
            echo "Target app is more recent. It will not be updated with source app contents."
        else
            APPS_TO_COPY="$APPS_TO_COPY$a "        
        fi
    else
        APPS_TO_COPY="$APPS_TO_COPY$a "
    fi
done

if [ -n "$APPS_TO_COPY" ]; then

    cf login $CF_SKIP_SSL_VALIDATION -a $CF_API \
        -u $CF_USER \
        -p $CF_PASSWORD \
        -o $CF_ORG \
        -s $CF_SPACE

    ARGS=""
    if [[ -n "$CF_TRACE" ]]; then
        ARGS="$ARGS --debug"
    fi
    if [[ -n "$COPY_AS_UPS" ]]; then
        ARGS="$ARGS --ups $COPY_AS_UPS"
    fi
    if [[ -n "$DEST_DOMAIN" ]]; then
        ARGS="$ARGS --domain $DEST_DOMAIN"
    fi
    if [[ -n "$DEST_HOST_FORMAT" ]]; then
        ARGS="$ARGS --host-format $DEST_HOST_FORMAT"
    elif [[ $CF_API == $CF_TARGET_API ]]; then
        ARGS="$ARGS --host-format {{.host}}-{{.space}}"
    fi

    cf copy $CF_TARGET_SPACE $CF_TARGET_ORG copy-dest --apps $APPS_TO_COPY $ARGS
fi

set +x
