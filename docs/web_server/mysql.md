如果出现Error 1292: Incorrect datetime value: '2019-04-16T18:09:02.146062768+08:00' for column 'token_expiry_time' at row 1
则去修改mysql配置文件

    my.ini中查找sql-mode，    
    我的MySQL版本为5.7.9，默认为：
    sql-mode="STRICT_ALL_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ZERO_DATE,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER"
    将红色标注的NO_ZERO_DATE,NO_ZERO_IN_DATE,删掉保存重启mysql即可；
    
    如果版本低的话默认可能是：
    默认为sql-mode="STRICT_TRANS_TABLES,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION"，
    将红色标注的STRICT_TRANS_TABLES,删掉保存重启mysql即可；