# https://ithelp.ithome.com.tw/questions/10209610
# https://stackoverflow.com/questions/7505988/importing-from-a-relative-path-in-python
# https://www.facebook.com/groups/pythontw/posts/10162487432003438/

# python -u "/home/caesar/CodeDev/experiment/py_import/tree/traversal.py"
# or
# python -m traversal

import sys,os

print(os.getcwd())
# /home/caesar/CodeDev/experiment/py_import/tree

print(os.getcwd().rsplit('/',1))
# ['/home/caesar/CodeDev/experiment/py_import', 'tree']

print()

# fail case
#
# sys.path.insert(0,'../')
# from ..tool import TestNode
#
# Traceback (most recent call last):
#   File "/home/caesar/CodeDev/experiment/py_import/tree/traversal.py", line 7, in <module>
#     from ..tool import TestNode
# ImportError: attempted relative import with no known parent package

# success case
# sys.path.insert(0,os.getcwd().rsplit('/',1)[0]) 意思等於 sys.path.insert(0,'../')
# sys.path.insert(0,os.getcwd().rsplit('/',1)[0])
sys.path.insert(0,'../')
from tool import TestNode

print(sys.path)
# ['../',
# '/home/caesar/CodeDev/experiment/py_import/tree',
# '/home/caesar/.pyenv/versions/3.10.1/lib/python310.zip',
# '/home/caesar/.pyenv/versions/3.10.1/lib/python3.10',
# '/home/caesar/.pyenv/versions/3.10.1/lib/python3.10/lib-dynload',
# '/home/caesar/.pyenv/versions/3.10.1/lib/python3.10/site-packages']

if __name__ == '__main__':
    print('\n',TestNode(10))