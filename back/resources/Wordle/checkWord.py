import twl, sys
if(len(sys.argv)==2):
    if isinstance(sys.argv[1], str):
        print(twl.check(sys.argv[1]))
