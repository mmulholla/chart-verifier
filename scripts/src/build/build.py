"""
Used by a github action to build a test a chart verifier docker image based on a pull request

parameters:
    ---image-name : Name to be used for the chart verifier image.
    --sha-value : Sha value to be used for the image.
    --verifier-version : Version of the chart-verifier to test.

results:
    exit code 1 if the image fails to build or a test fails.
"""

import argparse
import docker
import os
import sys
import tarfile


def build_image(image_id):
    print(f"Build Image : {image_id}")

    client = docker.from_env()

    try:
        image = client.images.build(path="./",tag=image_id)
        print("images:",image)
    except docker.errors.BuildError:
        print("[ERROR] docker build error")
        return False
    except  docker.errors.APIError:
        print ("[ERROR] docker API error")
        return False

    return True

tar_content_files = [ {"name": "config", "arc_name": "config"},
                      {"name": "out/chart-verifier", "arc_name": "chart-verifier"} ]

def create_tarfile(version):

    tgz_name = f"chart-verifier-{version}.tgz"
    if os.path.exists(tgz_name):
        os.remove(tgz_name)

    try:
        with tarfile.open(tgz_name, "x:gz") as tar:
            for tar_content_file in tar_content_files:
                tar.add(os.path.join(os.getcwd(),tar_content_file["name"]),arcname=tar_content_file["arc_name"])
    except Exception as err:
        print(f"[ERROR] Exception creating tarfile {err}")
        return False,""

    return True,os.path.join(os.getcwd(),tgz_name)


def main():

    parser = argparse.ArgumentParser()
    parser.add_argument("-i", "--image-name", dest="image_name", type=str, required=True,
                        help="Name of the chart verifier image")
    parser.add_argument("-s", "--sha-value", dest="sha_value", type=str, required=True,
                        help="Image sha value to test")
    parser.add_argument("-v", "--verifier-version", dest="verifier_version", type=str, required=True,
                        help="New version of chart verifier")


    args = parser.parse_args()

    image_id = f"{args.image_name}:{args.sha_value}"

    if build_image(image_id):
        print(f'::set-output name=verifier-image-tag::{args.sha_value}')
    else:
        sys.exit(1)

    outcome,tarfile_name = create_tarfile(args.verifier_version)
    if outcome:
        print(f'[INFO] Release asset created : {tarfile_name}.')
        print(f'::set-output name=PR_tarball_name::{tarfile_name}')
    else:
        sys.exit(1)
