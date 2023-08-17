# 执行编译
cd /root/dorycli/
git pull
date && time go build
ls -alh

# # 复制所有到new1-dory
# scp -r dorycli root@new1-dory:/usr/bin/

# # 复制所有到new2-dory
# scp -r dorycli root@new2-dory:/usr/bin/
