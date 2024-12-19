import twl, sys
if(len(sys.argv)==1):
    if isinstance(sys.argv[0], str):
        sys.stdout.write(twl.check(sys.argv[0]))
        sys.stdout.flush()