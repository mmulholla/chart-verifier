import argparse
import docker
import os
import sys
import yaml
try:
    from yaml import CLoader as Loader, CDumper as Dumper
except ImportError:
    from yaml import Loader, Dumper

sys.path.append('./scripts/src/')
from report import report_info

def build_image(image_id):
    print(f"Build Image : {image_id}")

    cwd = os.getcwd()

    # Print the current working directory
    print(f"Current working directory: {cwd}")

    client = docker.from_env()

    try:
        image = client.images.build(path="./",tag=image_id)
        print("logs:",image)
    except docker.errors.BuildError:
        print("docker build error")
        sys.exit(1)
    except  docker.errors.APIError:
        print ("docker API error")
        sys.exit(1)

def test_image(image_id,chart):

    docker_command = "verify " + chart["url"]

    set_values = ""
    vendor_type = ""
    profile_version = ""
    if "vendorType" in chart["metadata"]:
        vendor_type = chart["metadata"]["vendorType"]
        set_values = "profile.vendortype=%s" % vendor_type
    if "profileVersion" in chart["metadata"]:
        profile_version = chart["metadata"]["profileVersion"]
        if set_values:
            set_values = "%s,profile.version=%s" % (set_values,profile_version)
        else:
            set_values = "profile.version=%s" % profile_version

    if set_values:
        docker_command = "%s --set %s" % (docker_command, set_values)

    client = docker.from_env()
    out = client.containers.run(image_id,docker_command,stdin_open=True,tty=True,stderr=True)
    report = yaml.load(out, Loader=Loader)
    report_path = "banddreport.yaml"
    print("[INFO] report:\n", report)
    with open(report_path, "w") as fd:
        yaml.dump(report,fd)

    results = report_info.get_report_results(report_path,vendor_type,profile_version)

    expectedPassed = int(chart["results"]["passed"])
    expectedFailed = int(chart["results"]["failed"])

    if expectedFailed != results["failed"] or expectedPassed != results["passed"]:
        print("[ERROR] Chart verifier report includes unexpected results:")
        print(f'- Number of checks passed expected : {expectedPassed}, got {results["passed"]}')
        print(f'- Number of checks failed expected : {expectedFailed}, got {results["failed"]}')
        sys.exit(1)
    else:
        print(f'[PASS] Chart result validated : {chart["url"]}')





image_id = "quay.io/redhat-certification/chart-verifier:4a5b6cd"

build_image(image_id)

chart = {"url" : "https://github.com/redhat-certification/chart-verifier/blob/main/pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz?raw=true",
       "results":{"passed":"10","failed":"1"},
        "metadata":{"vendorType":"partner","profileVersion":"v1.0"}}

os.environ["VERIFIER_IMAGE"] = image_id
test_image(image_id,chart)
