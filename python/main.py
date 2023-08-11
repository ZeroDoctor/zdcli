import importlib
import argparse
import json

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('-s', '--script')
    parser.add_argument('-a', '--args')
    parser.add_argument('-f', '--function', default='main')
    args = parser.parse_args()
    
    print(f"[INFO] Starting base python script with [fucntion={args.function}] and [args={args.args}]...\n")
    
    mod = importlib.import_module(args.script)
    func = getattr(mod, args.function)

    result = func(**json.loads(args.args))
    if result != None:
        print(result)

if __name__ == "__main__":
    main()
