import _lib
import re
import time

def GetBlocks(datadir):
    _lib.StartTest("Load blocks chain")
    res = _lib.ExecuteNode(['printchain','-datadir',datadir, '-view', "short"])
    _lib.FatalAssertSubstr(res,"Hash: ","Blockchain display returned wrong data or no any blocks")
    
    regex = ur"Hash: ([a-z0-9A-Z]+)"

    blocks = re.findall(regex, res)
    
    return blocks

def WaitBlocks(datadir, explen):
    blocks = []
    i = 0
    while True:
        blocks = GetBlocks(datadir)
        
        if len(blocks) >= explen or i >= 5:
            break
        time.sleep(1)
        i = i + 1
        
    return blocks