"# ovpn_user" 

`此程序是用于批量创建openvpn-as用户,创建成功后,通过企业微信发送信息`
##使用方法
####可执行文件应放到openvpn-as所在服务器,且服务器内存在/usr/local/openvpn_as/scripts/sacli命令,若此命令路径不一样,请更改ovpn包下的方法`
####在可执行文件同级目录下创建`config.ini`文件,此文件用于配置redis及企业微信
####在可执行文件同级目录下创建`users.txt`文件,此文件用于配置需要创建的openvpn用户及其group,用户名和组名之间用","分割,如果只写用户名,则默认分配到default_group组中