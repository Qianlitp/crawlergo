#!/usr/bin/python3
# coding: utf-8

"""
    author: 猪猪侠 https://github.com/ring04h

"""

import logging
import subprocess

logging.basicConfig(level=logging.DEBUG)

# 
# (crontab -l;echo '0 2 * * * /usr/local/bin/python3 /data/script/zombie_clean.py') | crontab -
# 

def is_timeout(etime):
    if '-' in etime:
        day, hour = etime.split('-')
        return True if int(day) >= 1 else False
    else:
        return False


def cmdprocess(cmdline):

    pipe = subprocess.Popen(cmdline, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    output, stderr = pipe.communicate()
    return_code = pipe.returncode
    stderr = stderr.decode(errors='replace')
    output = output.decode(errors='replace')
    return output, stderr, return_code



def main():

    cmdline = "ps -ef | grep crawlergo | grep -v grep | awk '{print $2}'"
    output, stderr, return_code = cmdprocess(cmdline)
    
    if return_code != 0:
        return

    zombie_pids = output.splitlines()

    for zombie_pid in zombie_pids:

        cmdline = f'''ps -eo pid,etime | grep {zombie_pid}'''
        ps_output, ps_stderr, ps_return_code = cmdprocess(cmdline)

        if ps_return_code != 0:
            continue

        for line in ps_output.splitlines():
            
            pid, etime = line.split()

            status = is_timeout(etime)
            logging.debug(f"PID: {pid:<8} ETIME: {etime:<15} TIMEOUT: {status}")

            if not status: 
                continue

            kill_cmdline = f"kill -9 {pid}"
            logging.debug(f"call kill : [{kill_cmdline}]")

            cmdprocess(kill_cmdline)

if __name__ == "__main__":
    main()

