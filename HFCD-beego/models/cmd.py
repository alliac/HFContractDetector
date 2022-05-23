import subprocess
import sys
import os


def check_output(cmd):
    process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
    output = process.communicate()[0]
    if process.returncode != 0:
        raise subprocess.CalledProcessError(process.returncode, cmd, output=output)
    return output

try:
    # a = os.popen(sys.argv[1]).read()
    output = check_output(sys.argv[1])
except Exception as e:
    print('Exception=================>: ', str(e))