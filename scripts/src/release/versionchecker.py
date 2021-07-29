import re
import argparse
import os
import requests

version_file = "cmd/version.go"

def check_if_version_file_is_modified(api_url):
    # api_url https://api.github.com/repos/<organization-name>/<repository-name>/pulls/<pr_number>

    files_api_url = f'{api_url}/files'
    headers = {'Accept': 'application/vnd.github.v3+json'}
    pattern_versionfile = re.compile(version_file)
    page_number = 1
    max_page_size,page_size = 100,100


    while (page_size == max_page_size):

        files_api_query = f'{files_api_url}?per_page={page_size}&page={page_number}'
        r = requests.get(files_api_query,headers=headers)
        files = r.json()
        page_size = len(files)
        page_number += 1


        for f in files:
            filename = f["filename"]
            if pattern_versionfile.match(filename):
                return True

    return False

def get_version_from_file():
    version_go = open(version_file, "r")
    for line in version_go:
        if "var Version =" in line:
            name,version = line.split('=',1)
            version = ''.join(version.split())
            version = version.strip('\"')
            return version



def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-u", "--api-url", dest="api_url", type=str, required=False,
                        help="API URL for the pull request")
    parser.add_argument("-u", "--version", dest="version", type=str, required=False,
                        help="Version to compare")
    args = parser.parse_args()
    if not args.api_url or check_if_version_file_is_modified(args.api_url):
        new_version = get_version_from_file()
        if not args.version:
            print(f"::set-output name=version_in_PR::{new_version}")
        else if new_version > args.version:
            print("::set-output name=version_in_PR::true")


