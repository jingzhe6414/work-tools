#!/bin/bash
#$1 传需要上传的文件/目录
#$2 传minio中的目标目录

function minio() {
host=minio.example
s3_key=minio_admin
s3_secret=minio_secret


resource="$3/$1"  #自己定义传到minio的目录
content_type="application/octet-stream"
date=`date -R`
_signature="PUT\n\n${content_type}\n${date}\n${resource}"
signature=`echo -en ${_signature} | openssl sha1 -hmac ${s3_secret} -binary | base64`

curl -X PUT -T "$1" \
          -H "Host: ${host}" \
          -H "Date: ${date}" \
          -H "Content-Type: ${content_type}" \
          -H "Authorization: AWS ${s3_key}:${signature}" \
          http://${host}${resource}
#判定有些问题，返回值非200也会成功
if [ $? -eq 0 ];then
  echo "upload $1 success!"
else
  echo "upload $1 error!!!"
fi
}


#多级目录递归
function listFiles()
{

        for file in $1*;
        do
                if [ -f "$file" ]; then
                    echo $file
                    minio $file
                else
                    listFiles "$file/" " $2" $3
                fi
        done
}

listFiles $1 "" $2
