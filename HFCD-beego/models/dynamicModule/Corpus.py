import json

def corpus_fun(res):
    functionNames = list(res.keys())  # 所有函数名

    id = 0

    '''
    递归获取所有可能的组合
    index:当前计算的keys下标
    '''

    def dfs(index: int, ans: dict):
        nonlocal id  # 这样才能修改闭包内的外部变量
        if index == len(functionNames):  # 所有函数都找完了，保存为json
            jsonData = json.dumps(ans, indent=4)
            fileName = f'{id}.txt'
            with open(fileName, 'w') as f:
                f.write(jsonData)
            id += 1
            return
        functionName = functionNames[index]
        if res[functionName] is None:
            ans[functionName] = None
            dfs(index + 1, ans.copy())
        else:
            for i in range(len(res[functionName][1])):  # 遍历所有可能的参数
                ans[functionName] = res[functionName][1][i]['params']
                dfs(index + 1, ans.copy())

    # dfs(0, {})

    '''
    只获取包含所有参数的最少的文件
    '''

    def get_min_files():
        nonlocal id  # 这样才能修改闭包内的外部变量
        length = []  # 各个函数的参数个数
        for i in range(len(functionNames)):
            if res[functionNames[i]] is None or res[functionNames[i]][1] is None:
                length.append(0)
            else:
                length.append(len(res[functionNames[i]][1]))

        maxLength = max(length)

        for i in range(maxLength):  # 一共生成maxLength个文件
            ans = {}
            for j in range(len(functionNames)):
                if res[functionNames[j]] is None:
                    ans[functionNames[j]] = {"":""}
                elif res[functionNames[j]][1] is None:
                    ans[functionNames[j]] = {"":""}
                else:
                    k = i % len(res[functionNames[j]][1])
                    if 'args' not in res[functionNames[j]][1][k]['params']:
                        ans[functionNames[j]] = {"":""}
                        return
                    args = res[functionNames[j]][1][k]['params']['args']
                    #print(res[functionNames[j]][1][k]['params']['args'])
                    format_args={}
                    for  arg_index in range(len(args)):
                        format_args[str(arg_index)] = args[arg_index][1:-1]
                    #format_args = json.dumps(format_args)
                    #format_args = json.loads(format_args)
                    #print(format_args)
                    ans[functionNames[j]] = format_args#res[functionNames[j]][1][k]['params']
            jsonData = json.dumps(ans, indent=4)
            # print(jsonData)
            fileName = f'./models/dynamicModule/Corpus/{id}.txt'
            with open(fileName, 'w') as f:
                f.write(jsonData)
            id += 1

    get_min_files()


#fun()
