from kit import Repository
import argparse
import json
import sys
import os

def init_repo(path):
    if path is None:
        print("Error: path is required")
        sys.exit(1)
    repo = Repository(path)
    os.environ["MENACE_REPO_PATH"] = path
    print(f"Repository initialized at {path}")
    return repo

def get_file_tree():
    path = os.environ["MENACE_REPO_PATH"]
    repo = Repository(path)
    if repo is None:
        print("Error: repository not initialized")
        sys.exit(1)
    return repo.get_file_tree()

def find_symbols(symbol_name, symbol_type):
    path = os.environ["MENACE_REPO_PATH"]
    repo = Repository(path)
    if repo is None:
        print("Error: repository not initialized")
        sys.exit(1)
    result = repo.find_symbol_usages(symbol_name=symbol_name, symbol_type=symbol_type)
    if result is None:
        print("No results found")
        sys.exit(1)
    return result

def get_file_content(path):
    repo_path = os.environ["MENACE_REPO_PATH"]
    repo = Repository(repo_path)
    if repo is None:
        print("Error: repository not initialized")
        sys.exit(1)
    result = repo.get_file_content(path)
    if result is None:
        print("No results found")
        sys.exit(1)
    return result

def main():
    parser = argparse.ArgumentParser(description='Repository CLI Tool')
    parser.add_argument('function', choices=['init', 'file_tree', 'find_symbols', 'get_file_content'],
                      help='Function to execute')
    parser.add_argument('--path', help='Path to repository')
    parser.add_argument('--symbol', help='Symbol name for find_symbols')
    parser.add_argument('--symbol-type', help='Symbol type for find_symbols')
    parser.add_argument('--file-path', help='File path for get_file_content')

    args = parser.parse_args()

    if args.function == 'init':
        if not args.path:
            print("Error: --path is required for init")
            sys.exit(1)
        repo = init_repo(args.path)
        print(json.dumps({"Result": {"path": args.path}, "Status": 200}))
    
    elif args.function == 'file_tree':
        if not os.getenv("MENACE_REPO_PATH"):
            print("Error: repository not initialized")
            sys.exit(1)
        result = get_file_tree()
        print(json.dumps({"Result": result, "Status": 200}))
    
    elif args.function == 'find_symbols':
        if not all([args.symbol, args.symbol_type]):
            print("Error:--symbol, and --symbol-type are required for find_symbols")
            sys.exit(1)
        result = find_symbols(args.symbol, args.symbol_type)
        print(json.dumps({"Result": result, "Status": 200}))
    
    elif args.function == 'get_file_content':
        if not all([args.path]):
            print("Error: --path is required for get_file_content")
            sys.exit(1)
        result = get_file_content(args.path)
        print(json.dumps({"Result": result, "Status": 200}))

if __name__ == "__main__":
    main()