import re
import z3
import random
import subprocess
import os
import time
import chardet
from Corpus import corpus_fun

# 链码符号执行类 针对每个合约创建一个实例
class CCSE():
    # 实例初始化 参数：待测合约路径,生成M组测试文件,一个测试文件执行N次
    def __init__(self, loc=None, M=5, N=2, file_flag=None):
        try:
            # 提取原始合约内容 检测编码
            raw = open(loc, 'rb').read()
            result = chardet.detect(raw)
            encoding = result['encoding']
            fp = open(loc, 'r', encoding=encoding)
            self.chaincode = fp.read()
            fp.close()
        except Exception as e:
            raise Exception('待测合约路径错误', str(e))

# 通过符号执行获取init接口的参数组合  返回值：[{'conds': [约束条件], 'params': {解集}}, ...]
    def get_init_args(self):
        func_init = self.get_func_content('Init')
        if func_init and ('GetStringArgs()' in func_init or 'GetFunctionAndParameters()' in func_init):
            args = {}
            # 提取GetStringArgs()的返回值 即调用init接口的参数 为[]string类型
            if 'GetStringArgs()' in func_init:
                args[self.last_word(func_init,
                                    self.last_word(func_init, func_init.find('GetStringArgs()'))[1])[0]] = '[]string'
            # 提取GetFunctionAndParameters()的第二个返回值 即调用init接口的参数  为[]string类型
            else:
                args[self.last_word(func_init,
                                    self.last_word(func_init, func_init.find('GetFunctionAndParameters()'))[
                                        1])[0]] = '[]string'
            #print('函数名 Init')
            #print('参数列表', args)
            res=[{'fname':'Init'}]
            res.append(self.symbolic_execute(func_init, args))
            return res
        else:
            return None

# 通过符号执行获取invoke接口的参数组合  返回值：[{'conds': [约束条件], 'params': {解集}}, ...]
    def get_invoke_args(self):
        func_invoke = self.get_func_content('Invoke')
        # 提取GetFunctionAndParameters()的第一个返回值“通过invoke接口调用的函数” 为string类型
        if func_invoke and 'GetFunctionAndParameters' in func_invoke:
            args = {}
            args[self.last_word(func_invoke, self.last_word(func_invoke, self.last_word(func_invoke, func_invoke.find(
                'GetFunctionAndParameters()'))[1])[1])[0]] = 'string'
            #print('函数名 Invoke')
            #print('参数列表', args)
            res=[{'fname':'Invoke'}]
            res.append(self.symbolic_execute(func_invoke, args))
            return res
        else:
            return None
# 符号执行  参数：函数主体和参数列表  返回值：[{'conds': [约束条件], 'params': {解集}}, ...]
    def symbolic_execute(self, func_content=None, args=None):

        #print('---符号执行开始---')
        # 过滤“_”参数
        for key in args.keys():
            if key == '' or key == '_':
                del args[key]
                break
        # 求解参数
        z3_args = {}
        # 参数类型
        for key in args.keys():
            # 整数
            if args[key] in ['uint8', 'uint16', 'uint32', 'uint64', 'int8', 'int16', 'int32', 'int64', 'int']:
                z3_args[key] = z3.Int(key)
            # 实数 未测试
            elif args[key] == 'real':
                z3_args[key] = z3.Real(key)

        # 用于存储所有条件
        cond = []

        # 第一个if条件的位置 添加条件
        next_cond_loc = func_content.find('if ')
        while next_cond_loc != -1:
            next_cond = self.get_if_cond(func_content, next_cond_loc + 1)
            next_cond_loc = func_content.find('if ', next_cond_loc + 2)
            # 不考虑判断值是否为nil的条件
            if 'nil' in next_cond:
                continue
            # && || 拆开条件加入
            if '&&' in next_cond:
                cond.append(next_cond.split('&&')[0].strip())
                cond.append(next_cond.split('&&')[1].strip())
            elif '||' in next_cond:
                cond.append(next_cond.split('||')[0].strip())
                cond.append(next_cond.split('||')[1].strip())
            else:
                cond.append(next_cond)
        # print(cond)

        # 如果是invoke接口 不再进行条件组合 直接等于条件中的各函数名
        if 'Invoke(' in func_content:
            # 如果cond中没有if条件 继续考虑switch条件
            if not cond and 'switch' in func_content and 'case' in func_content:
                switch_cond = self.next_word(func_content, func_content.find('switch'))[0].strip()
                next_cond_loc = func_content.find('case ')
                while next_cond_loc != -1:
                    next_cond = self.next_word(func_content, next_cond_loc)[0].strip()[:-1]
                    next_cond_loc = func_content.find('case ', next_cond_loc + 4)
                    cond.append(switch_cond + ' == ' + next_cond)
            #*******************
            #print(cond)
            return cond
            # 有约束条件则直接取等值简化处理，条件为空则直接返回None
            if cond:
                res = []
                for c in range(len(cond)):
                    cp = {}
                    cp['conds'] = [cond[c]]
                    if '>=' in cond[c] or '<=' in cond[c]:
                        cond[c] = re.sub('<=|>=', '==', cond[c])
                    elif '>' in cond[c] or '<' in cond[c]:
                        cond[c] = re.sub('<|>', '==', cond[c])
                    elif '!=' in cond[c]:
                        cond[c] = re.sub('!=', '==', cond[c])
                    cp['params'] = {cond[c].split('==')[0].strip(): cond[c].split('==')[1].strip()}
                    res.append(cp)
                # 对第一个函数名取反 得到不和任何函数名相等的参数
                cp = {}
                cp['conds'] = [cond[0].split('==')[0].strip() + ' != any function']
                cp['params'] = {
                    cond[0].split('==')[0].strip(): self.random_char(len(cond[0].split('==')[1].strip()) - 2)
                }
                res.append(cp)
                # [{'conds': ['fn == "set"'], 'params': {'fn': '"set"'}}, {'conds': ['fn != any function'], 'params': {'fn': '"juz"'}}]
                #print('Invoke约束条件和解集：\n' + str(res) + '\n---符号执行结束---\n')
                return res
            else:
                return None

        # 所有能求解的条件的参数
        params_can_solve = list(args.keys())
        # 有约束条件则求解集，条件为空则直接返回None
        if cond:
            str_cond = []
            int_cond = []
            # 区分是数值条件还是字符串条件
            for c in cond:
                # 字符串、数组不用z3求解
                if 'len' in c or '"' in c or '[' in c:
                    # 替换条件中的函数内变量
                    need = True
                    if '==' in c:
                        para = c.split('==')[0].strip()
                    elif '!=' in c:
                        para = c.split('!=')[0].strip()
                    elif '<' in c:
                        para = c.split('<')[0].strip()
                    elif '<=' in c:
                        para = c.split('<=')[0].strip()
                    elif '>' in c:
                        para = c.split('>')[0].strip()
                    elif '>=' in c:
                        para = c.split('>=')[0].strip()
                    try:
                        for k in args.keys():
                            if k in para:
                                need = False
                    except:
                        pass
                    if need:
                        try:
                            var_loc = func_content.find(para + ' :=')
                            if var_loc != -1:
                                var_val = self.next_word(func_content, var_loc + len(para + ' :='))[0]
                                c = c.replace(para, var_val)
                                params_can_solve.append(para)
                            var_loc = func_content.find(para + ':=')
                            if var_loc != -1:
                                var_val = self.next_word(func_content, var_loc + len(para + ':='))[0]
                                c = c.replace(para, var_val)
                                params_can_solve.append(para)
                            var_loc = func_content.find('var ' + para + ' =')
                            if var_loc != -1:
                                var_val = self.next_word(func_content, var_loc + len('var ' + para + ' ='))[0]
                                c = c.replace(para, var_val)
                                params_can_solve.append(para)
                            var_loc = func_content.find('var ' + para + '=')
                            if var_loc != -1:
                                var_val = self.next_word(func_content, var_loc + len('var ' + para + '='))[0]
                                c = c.replace(para, var_val)
                                params_can_solve.append(para)
                        except:
                            print('替换参数错误')
                    str_cond.append(c)

                # 整数、实数用z3求解
                else:
                    # 找到所有的变量声明 替换条件中的变量
                    var_loc = func_content.find(":=")
                    while var_loc != -1:
                        var_name = self.last_word(func_content, var_loc)[0]
                        var_val = self.next_word(func_content, var_loc + 2)[0]
                        c = c.replace(var_name.strip(), '(' + var_val + ')')
                        var_loc = func_content.find(":=", var_loc + 2)

                    # 替换条件 $$$防止变量名和z3_args冲突
                    for key in args.keys():
                        c = c.replace("z3_args", "$$$")
                        c = c.replace(key, "z3_args['" + key + "']")
                        c = c.replace("$$$", "z3_args")
                    int_cond.append(c)
            # print(str_cond)
            # print(int_cond)

            # 筛除一些无法求解的数值条件（比如条件里的值涉及到函数调用）
            for i in int_cond.copy():
                try:
                    z3.Not(eval(i))
                except:
                    int_cond.remove(i)
            # 把字符串长度的大于小于条件转化成等于
            for i in str_cond.copy():
                if '>=' in i or '<=' in i:
                    i = re.sub('<=|>=', '==', i)
                elif '>' in i or '<' in i:
                    i = re.sub('<|>', '==', i)

            # 添加第一组全部未取反的整数、实数型条件
            int_conds = [int_cond]
            # 依次将每一个int_cond条件取反，加到条件组int_conds中
            for i in range(len(int_cond)):
                temp_conds = int_conds.copy()
                for c in temp_conds:
                    c_copy = c.copy()
                    c_copy[i] = z3.Not(eval(c[i]))
                    int_conds.append(c_copy)

            # 添加第一组全部未取反的字符串、数组型条件
            str_conds = [str_cond]
            # 依次将每一个str_cond条件取反，加到条件组str_conds中
            for i in range(len(str_cond)):
                temp_conds = str_conds.copy()
                for c in temp_conds:
                    c_copy = c.copy()
                    c_copy[i] = self.reverse_cond(c[i])
                    str_conds.append(c_copy)
            # print(str_conds)
            # print(int_conds)

            # 解析每一个条件组，替换参数执行
            res = []
            for str_c in str_conds:
                for int_c in int_conds:
                    #print('约束条件：\n', str_c + int_c)
                    conds_and_params = {}
                    conds_and_params['conds'] = str_c + int_c

                    # z3符号执行求解int条件
                    solver = z3.Solver()
                    # 检查每一个条件，是字符串则eval转化，不是字符串则直接add
                    solver.reset()
                    for c in int_c:
                        if isinstance(c, str):
                            solver.add(eval(c))
                        else:
                            solver.add(c)
                    # 如果有解
                    if str(solver.check()) == 'sat':
                        ans = solver.model()
                        # print('解集：\n', ans)

                        # 如果int条件有解 继续求解str条件 替换参数
                        new_params = dict.fromkeys(args.keys())
                        # 解析条件组里的每一个条件
                        for c in str_c:
                            try:
                                # 如果条件中没有能求解的参数
                                can_solve = False
                                for pa in params_can_solve:
                                    if pa in c:
                                        can_solve = True
                                        break
                                if can_solve:
                                    # 如果是相等的条件
                                    if '==' in c:
                                        para = c.split('==')[0].strip()
                                        val = c.split('==')[1].strip()
                                        # 如果是len相关条件
                                        if 'len' in para:
                                            arg = self.find_para(para, 'len')
                                            new_params[arg] = []
                                            for i in range(int(val)):
                                                new_params[arg].append('""')
                                        # 如果是数组值相关条件
                                        elif '[' in para:
                                            arg, index = self.find_para(para, '[]')
                                            if arg in new_params.keys() and new_params[arg]:
                                                try:
                                                    new_params[arg][index] = val
                                                except:
                                                    for j in range(len(new_params[arg]), index + 1):
                                                        new_params[arg].append('""')
                                                    new_params[arg][index] = val
                                            else:
                                                new_params[arg] = []
                                                for j in range(index + 1):
                                                    new_params[arg].append('""')
                                                new_params[arg][index] = val
                                        # 如果是字符串条件
                                        else:
                                            new_params[para] = val
                                    # 如果是不等的条件
                                    elif '!=' in c:
                                        para = c.split('!=')[0].strip()
                                        val = c.split('!=')[1].strip()
                                        # 如果是len相关条件
                                        if 'len' in para:
                                            arg = self.find_para(para, 'len')
                                            new_params[arg] = []
                                            for i in range(int(val) + 1):
                                                new_params[arg].append('""')
                                        # 如果是数组值相关条件
                                        elif '[' in para:
                                            arg, index = self.find_para(para, '[]')
                                            if arg in new_params.keys() and new_params[arg]:
                                                try:
                                                    if new_params[arg][index] == '""':
                                                        if val != '""':
                                                            new_params[arg][index] = self.random_char(len(val) - 2)
                                                        else:
                                                            new_params[arg][index] = self.random_char(
                                                                random.randint(2, 5))
                                                except:
                                                    for j in range(len(new_params[arg]), index + 1):
                                                        new_params[arg].append('""')
                                                    if val != '""':
                                                        new_params[arg][index] = self.random_char(len(val) - 2)
                                                    else:
                                                        new_params[arg][index] = self.random_char(random.randint(2, 5))
                                            else:
                                                new_params[arg] = []
                                                for j in range(index + 1):
                                                    new_params[arg].append('""')
                                                if val != '""':
                                                    new_params[arg][index] = self.random_char(len(val) - 2)
                                                else:
                                                    new_params[arg][index] = self.random_char(random.randint(2, 5))
                                        # 如果是字符串条件
                                        else:
                                            # 防止前面已经有相等的条件把字符串限定
                                            if new_params[para] == '""' or new_params[para] is None:
                                                if val != '""':
                                                    new_params[para] = self.random_char(len(val) - 2)
                                                else:
                                                    new_params[para] = self.random_char(random.randint(2, 5))
                            except Exception as e:
                                print(c, '约束求解错误',str(e))
                        for p in new_params.keys():
                            # 值为None一般代表是int参数
                            if not new_params[p]:
                                try:
                                    # print(p,ans[z3_args[p]])
                                    new_params[p] = ans[z3_args[p]]
                                except:
                                    # 也可能是条件求解过程中出错 生成随机参数
                                    #print(p)
                                    if 'int' in args[p]:
                                        new_params[p] = self.random_value('int')
                                    elif args[p] == 'string':
                                        new_params[p] = self.random_value('string')
                                    elif args[p] == '[]string':
                                        new_params[p] = self.random_value('[]string')
                                    elif args[p] == 'real':
                                        new_params[p] = self.random_value('real')
                        # 最后把解集中所有的空字符串换成固定字符串或随机字符串
                        for p in new_params.keys():
                            if str(type(new_params[p])) == "<class 'list'>":
                                for index in range(len(new_params[p])):
                                    if new_params[p][index] == '""':
                                        new_params[p][index] = '"abc"'
                                        # new_params[p][index] = self.random_char(random.randint(2, 6))
                            elif str(type(new_params[p])) == "<class 'str'>":
                                if new_params[p] == '""':
                                    new_params[p] = '"abc"'
                                    # new_params[p] = self.random_char(random.randint(2, 6))
                        # 弥补z3的一个bug 0参数加1
                        for key in new_params.keys():
                            if str(type(new_params[key])) == "<class 'z3.z3.IntNumRef'>":
                                if new_params[key] == 0:
                                    new_params[key] = 1
                                else:
                                    # 加法
                                    if int(str(new_params[key])) > 0:
                                        new_params[key] = int(str(new_params[key])) - 1
                                    # 减法
                                    else:
                                        new_params[key] = int(str(new_params[key])) + 1
                        #print('解集：\n', new_params)
                        conds_and_params['params'] = new_params
                        res.append(conds_and_params)
                    else:
                        print('无解')
                        conds_and_params['params'] = None
            #print('---符号执行结束---\n')
            return res
        else:
            return None
    def get_func_content(self, func_name):
        # 在整个链码中匹配函数头
        res = re.search('func (\(.*\) )?' + func_name, self.chaincode)
        if res:
            # 提取匹配结果的起始位置
            func_loc = res.span()[0]
            start = func_loc
            end = -1
            count = 0
            for i in range(func_loc + 7 + len(func_name), len(self.chaincode)):
                if self.chaincode[i] == '{':
                    # if count == 0:    只提取函数主体，start初值为-1，最后return的条件为start != -1 and end != -1
                    #     start = i
                    count += 1
                elif self.chaincode[i] == '}':
                    count -= 1
                    if count == 0:
                        end = i + 1
                        break

            if end != -1:
                return self.chaincode[start:end]
            else:
                return None
        else:
            return None
# 返回从index位置开始的上一个单词及位置
    def last_word(self, str, index):
        start = -1
        end = -1
        for i in range(index, 0, -1):
            if str[i] in [' ', '\n', '\t', '(', ')', ',']:
                i -= 1
                while str[i] in [' ', '\n', '\t', '(', ')', ',']:
                    i -= 1
                end = i + 1
                break
        for i in range(end - 1, 0, -1):
            if str[i] in [' ', '\n', '\t', '(', ')', ',']:
                start = i + 1
                break
        if start != -1 and end != -1:
            return str[start:end], start
        else:
            return None
    # 返回从index位置开始的下一个单词及位置
    def next_word(self, str, index):
        start = -1
        end = -1
        for i in range(index, len(str)):
            if str[i] in [' ', '\n', '\t', '(', ')', ',']:
                i += 1
                while str[i] in [' ', '\n', '\t', '(', ')', ',']:
                    i += 1
                start = i
                break
        for i in range(start + 1, len(str)):
            if str[i] in [' ', '\n', '\t', '(', ')', ',']:
                end = i
                break
        if start != -1 and end != -1:
            return str[start:end], start
        else:
            return None
 # 获取if后面的条件 参数：函数主体内容，if的位置 返回：条件内容
    def get_if_cond(self, func_content, if_loc):
        start = -1
        end = -1
        for i in range(if_loc, len(func_content)):
            if start == -1 and func_content[i] == ' ':
                start = i + 1
            elif func_content[i] == '{':
                end = i
                break
        # print(start,end,self.prog[start:end].strip())
        if start != -1 and end != -1:
            return func_content[start:end].strip()
        else:
            return None
# 随机产生n位字母的字符串
    def random_char(self, length):
        string = ''
        for i in range(length):
            string += chr(random.randint(97, 97 + 25))
        return '"' + string + '"'
# 通过符号执行获取某函数的参数组合 参数：函数名  返回值：[{'conds': [约束条件], 'params': {解集}}, ...]
    def get_func_args(self, func_name):
        func_content = self.get_func_content(func_name)
        if func_content:
            # print(func_content)
            # 从函数头中提取参数列表
            args = self.get_params_list(func_content, func_name)
            #print('函数名', func_name)
            #print('参数列表', args)
            res=[{'fname':func_name}]
            res.append(self.symbolic_execute(func_content, args))
            return res
        else:
            return None
# 提取函数的参数列表  参数：函数主体，函数名  返回值：
    def get_params_list(self, func_content, func_name):
        try:
            func_head = func_content.split('\n')[0]
            start = func_head.find(func_name) + len(func_name) + 1
            end = -1
            for i in range(start, len(func_head)):
                if func_head[i] == ')':
                    end = i
                    break
            if end != -1:
                res = {}
                params_list = func_head[start:end].split(',')
                for p in params_list:
                    p = p.strip().split()
                    # 只保留int、string、[]string类型
                    if p[1] in ['uint8', 'uint16', 'uint32', 'uint64', 'int8', 'int16', 'int32', 'int64', 'int',
                                'string', '[]string']:
                        res[p[0]] = p[1]
                # print(res)
                return res
            else:
                return None
        except:
            return None
 # 取反字符串条件
    def reverse_cond(self, str):
        if '!=' in str:
            return str.replace('!=', '==')
        elif '==' in str:
            return str.replace('==', '!=')
 # 生成随机值 参数：数据类型
    def random_value(self, num_type):
        if num_type == 'int':
            return str(random.randint(1, 100))
        if num_type == 'string':
            return self.random_char(random.randint(2, 6))
        if num_type == '[]string':
            return [self.random_char(random.randint(2, 6))] * random.randint(0, 5)
        if num_type == 'real':
            return str(random.random())
# 找到len、[]等字符串条件中的参数
    def find_para(self, cond, way=None):
        # 如果是len相关条件
        if way == 'len':
            return cond[cond.find('len') + 4:-1]
        # 如果是数组值相关条件
        elif way == '[]':
            temp = cond.split('[')
            return (temp[0], int(temp[1][:-1]))
        # 如果是字符串条件
        else:
            return cond

#产生初始语料库
    def get_corpus(self):
        res = {}
        init_args = self.get_init_args()
        #print(init_args)
        res['Init'] = init_args
        invoke_args = self.get_invoke_args()
        #print(invoke_args)
        for i in range(len(invoke_args[1])):
                #print(invoke_args[1][i][7:-1])
                args=self.get_func_args(invoke_args[1][i][7:-1])
                #print(args)
                res[invoke_args[1][i][7:-1]] = args
        #fileName = f'./DynamicModule/SymRes.txt'
        #with open(fileName, 'w') as f:
            #f.write(str(res))
        corpus_fun(res)
if __name__ == '__main__':
    ccse = CCSE('./models/chaincode.go', 2, 2)
    ccse.get_corpus()
