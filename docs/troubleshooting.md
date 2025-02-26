# Troubleshooting guidelines and known solutions

> Help us to help you to resolve issue you are having.

## Tools to help with troubleshooting issues

- run command `rcc configuration diagnostics` and see if there are warnings
  or errors in output
- if failure is with specific robot, then try running command
  `rcc configuration diagnostics --robot path/to/robot.yaml` and see if
  those robot diagnostics have something that identifies a problem
- run command `rcc configuration speedtest` to see, if problem is actually
  performance related (like slow disk or network access)
- run rcc commands with `--debug` and `--timeline` flags, and see if anything
  there adds more information on why failure is happening

## How to troubleshoot issue you are having?

- are you using latest versions of tools and libraries, and if not, then first
  thing to do is update them and then retry
- if you have only encoutered the problem just once, try to repeat it, and
  if you cannot repeat it, you are done
- has this ever before worked in this same user/machine/network combination,
  and if not, then are you using correct profile and settings?
- if this worked previously, and then stopped working, what has changed or
  what have you changed? any updates? new IT policies? new network location?
- gather evidence, that is all logs, console outputs, stack traces, screenshots,
  and look them thru
- what is first error your see in console output, and what is last error your
  see, then look between

## Reporting an issue

- describe what were you trying to achieve
- describe what did you actually do when trying to achieve it
- describe what did actually happen
- describe what were you expecting to happen
- describe what did happen that you think indicates that there is an issue
- describe what error messages did you see
- describe steps that are needed to be able to reproduce this issue
- describe what have you already tried to resolve this issue
- describe what has changed since this was not present and everything worked ok
- you should share your `conda.yaml` used with robot or environment
- you should share your `robot.yaml` that defines your robot
- you should share your code, or minimal sample code, that can reproduce
  problem you are having

## Known solutions

### Access denied while building holotree environment (Windows)

If file is .dll or .exe file, then there is probably some process running, that
has actually locked that file, and tooling cannot complete its operation while
that other process is running. Other process might be virus scanner, some other
tool (Assistant, Workforce Agent, Automation Studio, VS Code, rcc) using same
environment, or even open Explorer view.

To resolve this, close other applications, or wait them to finish before trying
same operation again.

### Message "Serialized environment creation" repeats

There can be few reasons for this. Here are some ways to resolve it.

If multiple robots in same machine are trying to create new environment or
refresh existing one at exactly same time, then only one of them can continue.
This is there to protect integrity and security of holotree and hololib, and
also conserve resources for doing duplicate work. In this case, best thing to
resolve this is just to wait processes to complete.

Other case is where there are multiple rcc processes running, but none of them
seems to be progressing. This might be indication that there is one "zombie"
process, which is holding on to a lock, and wont go away since some of its
child processes is still running (like python, web browser, or Excel). In this
case, best way is to close those "hanging" processes, and let OS to finish
that pending (and lock holding) process.

Third case is where there seems to be only one rcc, and it is just waiting and
repeating that message. In this case it is probably a permission issue, and
for some reason .lck file is not accessible/lockable by rcc process. In this
case, you should go and look if current user has rights to actually modify
those .lck files, and if not, you have to grant them those. This might require
administrator privileges to actually change those file permissions.
