
def fun():
    mp = {
        '内置函数':0,
        '字段声明漏洞':0,
        '多重继承':0,
        '程序并发性漏洞':0,
        '映射结构迭代':0,
        '未初始化存储指针':0,
        '全局变量漏洞':0,
        '写后读漏洞':0,
        '未使用的隐私数据机制':0,
        '范围查询风险':0,
        '隐私数据安全风险':0,
        '系统命令执行漏洞':0,
        '外部库调用漏洞':0,
        'Web服务漏洞':0,
        '外部文件访问':0,
        '随机数生成漏洞':0,
        '系统时间戳漏洞':0,
        '未处理的错误':0,
        '未加密的敏感数据':0,
        '函数未检查输入参数':0,
        '注释标题不足以检查实现和使用情况':0,
        '存在无限循环':0,

    }
    with open('./models/TestReport/test_report.txt', 'r', encoding = 'utf8') as f:
        lines = f.readlines()
        for line in lines:
            for key in mp:
                if key in line:
                    mp[key] += 1
                    break
    type_count = 0
    total_count = 0                
    for key in mp:
        if mp[key] != 0:
            print(f'{key}({mp[key]})')
            type_count += 1
            total_count += mp[key]
    print(f'*Result: 总计{type_count}种漏洞类型，{total_count}个漏洞')


fun()