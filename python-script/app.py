import ConfigParser
import time
import json
import ipfsapi
from flask import Flask, request
from iota_cache import IotaCache
from tag_generator import TagGenerator

import requests

app = Flask(__name__)

cf = ConfigParser.ConfigParser()
cf.read("conf")
iota_addr = cf.get("iota", "addr")
iota_seed = cf.get("iota", "seed")
cache = IotaCache(iota_addr, iota_seed)
ipfs_ip = cf.get("ipfs", "ip")
ipfs_port = cf.get("ipfs", "port")
ipfs_client = ipfsapi.connect(ipfs_ip, ipfs_port)


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
    #     result = ipfs_client.block_put(request.files[file_key])
    #     ipfs_hash = result["Key"]
    #     if ipfs_hash[:2] != "Qm":
    #         return 'ipfs error'
    #     hash_list.append(ipfs_hash)
    #
    # for hash in hash_list:
    #     print hash

    result_dict = {}
    contract_file = request.files.get('contract', None)
    if contract_file is None:
        return 'no contract'

    ipfs_result = ipfs_client.block_put(contract_file)
    ipfs_hash = ipfs_result['Key']
    if ipfs_hash[:2] != 'Qm':
        return 'ipfs error'
    result_dict['address'] = ipfs_hash

    command = request.form.get("command", None)
    if command is None:
        return 'no command'
    result_dict['command'] = command

    result_dict['contractName'] = 'contract_' + str(time.time())
    result_dict['contractType'] = 'fabric'
    result_dict["contractVersion"] = "v1.0"
    result_dict['vmVersion'] = '1.0'
    result_dict['sequence'] = '10'
    result_dict['timestamp'] = str(time.time())
    result_dict['user'] = 'user1'
    result_dict['signature'] = ''
    result_dict['operation'] = 'install'

    # temporary post data to tvm
    r = requests.post('http://192.168.199.129:8080/executeContract', json=json.dumps(result_dict))
    if r.status_code == 200:
        return 'ok'
    else:
        return 'error'

    # # cache in tangle part
    #
    # result_flag, result_msg = cache_json_in_tangle(result_dict)
    #
    # if result_flag is True:
    #     return 'ok'
    # else:
    #     return result_msg


@app.route('/upload_action', methods=['POST'])
def upload_action():
    pass


def cache_json_in_tangle(result_json):

    ipfs_hash = ipfs_client.add_str(json.dumps(result_json, sort_keys=True))

    if ipfs_hash[:2] == "Qm":
        print("[INFO]Cache json %s in ipfs, the hash is %s." % (json.dumps(result_json, sort_keys=True), ipfs_hash))

        cache.cache_txn_in_tangle(ipfs_hash, TagGenerator.get_current_tag())
        print("[INFO]Cache hash %s in tangle, the tangle tag is %s." % (ipfs_hash, TagGenerator.get_current_tag()))

        return True, 'ok'

    else:
        print("[ERROR]Cache ipfs error, error message is %s." % ipfs_hash)
        return False, 'ipfs error'


if __name__ == '__main__':
    app.run()


