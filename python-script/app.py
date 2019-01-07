import ConfigParser
import time
import json
import hashlib
import ipfsapi
import requests
from flask import Flask, request
from iota_cache import IotaCache
from tag_generator import TagGenerator


app = Flask(__name__)

cf = ConfigParser.ConfigParser()
cf.read("conf")

iota_addr = cf.get("iota", "addr")
iota_seed = cf.get("iota", "seed")
cache = IotaCache(iota_addr, iota_seed)

ipfs_ip = cf.get("ipfs", "ip")
ipfs_port = cf.get("ipfs", "port")
ipfs_client = ipfsapi.connect(ipfs_ip, ipfs_port)

tvm_addr = cf.get("tvm", "addr")


@app.route('/')
def hello_world():
    return 'Hello World!'


@app.route('/upload_transaction', methods=['POST'])
def upload_transaction():
    req_json = request.get_json()

    if req_json is None:
        return 'request error'

    req_json["timestamp"] = str(time.time())

    result_flag, result_msg = cache_json_in_tangle(req_json)

    if result_flag is True:
        return 'ok'
    else:
        return result_msg


@app.route('/upload_contract', methods=['POST'])
def upload_contract():

    # # multi-file part
    #
    # if len(request.files) == 0:
    #     return 'request error'
    # hash_list = []
    # for file_key in request.files:
    #     ipfs_hash = ipfs_client.add_str(request.files[file_key].read())
    #     if ipfs_hash[:2] != "Qm":
    #         return 'ipfs error'
    #     hash_list.append(ipfs_hash)
    #
    # for hash in hash_list:
    #     print hash

    result_dict = dict()
    result_dict['type'] = 'contract'
    contract_file = request.files.get('contract', None)
    if contract_file is None:
        return 'no contract'
    ipfs_hash = ipfs_client.add_str(contract_file.read())

    if ipfs_hash[:2] != 'Qm':
        return 'ipfs error'
    result_dict['address'] = ipfs_hash

    md5 = request.form.get('md5', None)
    if md5 is None:
        return 'no md5'
    result_dict['checkMD5'] = md5

    result_dict['command'] = ''
    result_dict['contractName'] = 'contract_' + str(int(round(time.time())))
    result_dict['contractType'] = 'fabric'
    result_dict["contractVersion"] = "v1.0"
    result_dict['vmVersion'] = '1.0'
    result_dict['sequence'] = '0'
    result_dict['timestamp'] = int(round(time.time()))
    result_dict['user'] = 'user1'
    result_dict['signature'] = ''
    result_dict['operation'] = 'install'

    print result_dict

    # cache in tangle part
    result_flag, result_msg = cache_json_in_tangle(result_dict)

    if result_flag is True:
        return "".join([
            "SUCCESS \n\n",
            "[INIT-CMD]:\n",
            "curl -X POST \\ \n",
            "  http://127.0.0.1:5000/upload_action \\ \n",
            "  -H 'cache-control: no-cache' \\ \n",
            "  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \\ \n",
            "  -F md5=", result_dict['checkMD5'], " \\ \n",
            "  -F contractName=", result_dict['contractName'], " \\ \n",
            """  -F 'command={"Args":["init","a","7","b","8"]}' \\ \n""",
            "  -F operation=instantiate \\ \n",
            "  -F address=", result_dict['address'], " \n"
        ])
    else:
        return result_msg


@app.route('/upload_action', methods=['POST'])
def upload_action():

    result_dict = dict()
    result_dict['type'] = 'contract'

    address = request.form.get('address', None)
    if address is None:
        return 'no address'
    result_dict['address'] = address

    md5 = request.form.get('md5', None)
    if md5 is None:
        return 'no md5'
    result_dict['checkMD5'] = md5

    contract_name = request.form.get("contractName", None)
    if contract_name is None:
        return 'no contractName'
    result_dict['contractName'] = contract_name

    command = request.form.get("command", None)
    if command is None:
        return 'no command'
    result_dict['command'] = json.dumps(json.loads(command))

    result_dict['contractType'] = 'fabric'
    result_dict["contractVersion"] = "v1.0"
    result_dict['vmVersion'] = '1.0'
    result_dict['sequence'] = '0'
    result_dict['timestamp'] = int(round(time.time()))
    result_dict['user'] = 'user1'
    result_dict['signature'] = ''
    operation = request.form.get('operation', None)
    if operation is None:
        return 'no operation'
    elif operation not in ['instantiate', 'query']:
        return 'wrong operation'
    result_dict['operation'] = operation

    print result_dict

    # cache in tangle part
    result_flag, result_msg = cache_json_in_tangle(result_dict)

    if result_flag is True:
        return "".join([
            "SUCCESS \n\n",
            "[QUERY-CMD]:\n",
            "curl -X POST \\ \n",
            "  http://127.0.0.1:5000/query_action \\ \n",
            "  -H 'cache-control: no-cache' \\ \n",
            "  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \\ \n",
            "  -F md5=", result_dict['checkMD5'], " \\ \n",
            "  -F contractName=", result_dict['contractName'], " \\ \n",
            """  -F 'command={"Args":["query","a"]}' \\ \n""",
            "  -F address=", result_dict['address'], " \n"
        ])
    else:
        return result_msg


@app.route('/query_action', methods=['POST'])
def query_action():

    result_dict = dict()
    result_dict['type'] = 'contract'

    address = request.form.get('address', None)
    if address is None:
        return 'no address'
    result_dict['address'] = address

    md5 = request.form.get('md5', None)
    if md5 is None:
        return 'no md5'
    result_dict['checkMD5'] = md5

    contract_name = request.form.get("contractName", None)
    if contract_name is None:
        return 'no contractName'
    result_dict['contractName'] = contract_name

    command = request.form.get("command", None)
    if command is None:
        return 'no command'
    result_dict['command'] = json.dumps(json.loads(command))

    result_dict['contractType'] = 'fabric'
    result_dict["contractVersion"] = "v1.0"
    result_dict['vmVersion'] = '1.0'
    result_dict['sequence'] = '0'
    result_dict['timestamp'] = int(round(time.time()))
    result_dict['user'] = 'user1'
    result_dict['signature'] = ''
    result_dict['operation'] = 'query'

    print result_dict

    headers = {
        'Content-Type': "application/json",
        'cache-control': "no-cache"
    }
    result = requests.get(tvm_addr + '/executeContract', data=json.dumps(result_dict), headers=headers)
    if result.status_code == 200:
        return result.content
    else:
        return "tvm error"


def cache_json_in_tangle(result_json):

    ipfs_hash = ipfs_client.add_str(json.dumps(result_json, sort_keys=True))

    if ipfs_hash[:2] == "Qm":
        print("[INFO]Cache json %s in ipfs, the hash is %s." % (json.dumps(result_json, sort_keys=True), ipfs_hash))

        cache.cache_txn_in_tangle_sdk(ipfs_hash, TagGenerator.get_current_tag())
        print("[INFO]Cache hash %s in tangle, the tangle tag is %s." % (ipfs_hash, TagGenerator.get_current_tag()))

        return True, 'ok'

    else:
        print("[ERROR]Cache ipfs error, error message is %s." % ipfs_hash)
        return False, 'ipfs error'


def get_md5(file_content):

    hash_md5 = hashlib.md5(file_content)

    return hash_md5.hexdigest()


if __name__ == '__main__':
    app.run()


