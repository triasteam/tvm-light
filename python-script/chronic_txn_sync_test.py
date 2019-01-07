import json
from mock import MagicMock, call
import unittest
from tag_generator import TagGenerator
from chronic_txn_sync import push_and_sync, cache, tm
from chronic_txn_sync import ipfs_client


class TestChronicTxnSync(unittest.TestCase):

    def setUp(self):
        pass

    def mock_get_non_consumed_txns(self, tag):
        return self.ipfs_addrs

    def mock_get_json(self, ipfs_addr):
        index = self.ipfs_addrs.index(ipfs_addr)
        return self.transactions[index]

    def mock_set_txn_as_synced(self, ipfs_addr, tag):
        pass

    def mock_broadcast_tx_async(self, tx):
        pass

    def mock_broadcast_txs_async(self, txs):
        pass

    def test_push_and_sync_for_one_tx(self):

        self.ipfs_addrs = [
            "QmSRjsq7wg7njG2HwyHsRT2qUT9teWccjv65Ud9TcWG3Sz"
        ]
        self.transactions = [
            {"to": "b", "from": "a", "timestamp": "1543392247.23", "transfer": "1"}
        ]

        cache.get_non_consumed_txns = MagicMock(side_effect=self.mock_get_non_consumed_txns)
        cache.set_txn_as_synced = MagicMock(side_effect=self.mock_set_txn_as_synced)
        tm.broadcast_tx_async = MagicMock(side_effect=self.mock_broadcast_tx_async)
        tm.broadcast_txs_async = MagicMock(side_effect=self.mock_broadcast_txs_async)

        current_tag = TagGenerator.get_current_tag()
        push_and_sync(current_tag)

        for ipfs_addr in self.ipfs_addrs:
            cache.set_txn_as_synced.assert_any_call(ipfs_addr, current_tag)

        tm.broadcast_tx_async.assert_any_call(json.dumps(self.transactions[0], sort_keys=True))

    # def test_push_and_sync_for_more_tx(self):
    #
    #     self.ipfs_addrs = [
    #         "QmSRjsq7wg7njG2HwyHsRT2qUT9teWccjv65Ud9TcWG3Sz",
    #         "QmURmxLAmL5UGjXA8Vt4kak2hjXHnjZGjphScujY26SRN9",
    #         "QmS15xXJsbDqqhVLXor2Gb48izGEbWsxrg3F6ERtXQTCvM"
    #     ]
    #     self.transactions = [
    #         {"to": "b", "from": "a", "timestamp": "1543392247.23", "transfer": "1"},
    #         {"to": "d", "from": "c", "timestamp": "1543392284.83", "transfer": "2"},
    #         {"to": "f", "from": "e", "timestamp": "1543392301.58", "transfer": "3"},
    #     ]
    #
    #     cache.get_non_consumed_txns = MagicMock(side_effect=self.mock_get_non_consumed_txns)
    #     cache.set_txn_as_synced = MagicMock(side_effect=self.mock_set_txn_as_synced)
    #     tm.broadcast_tx_async = MagicMock(side_effect=self.mock_broadcast_tx_async)
    #     tm.broadcast_txs_async = MagicMock(side_effect=self.mock_broadcast_txs_async)
    #
    #     current_tag = TagGenerator.get_current_tag()
    #     push_and_sync(current_tag)
    #
    #     for ipfs_addr in self.ipfs_addrs:
    #         cache.set_txn_as_synced.assert_any_call(ipfs_addr, current_tag)
    #
    #     transaction_list = []
    #     for transaction in self.transactions:
    #         transaction_list.append(transaction)
    #     tm.broadcast_txs_async.assert_called_once_with(json.dumps({"txs": transaction_list}, sort_keys=True))

    def test_push_and_sync_for_one_contract(self):

        self.ipfs_addrs = [
            "QmURGGVodbJv1LnkNqMLo3uvpCPo67KcDk1jYkXJh8gmSU"
        ]
        self.transactions = [
            {
                "address": "QmStXVUdqAbKDTecFQkawuvSNqPs6su5KJB94Uvd9MiCny",
                "checkMD5": "7a3f59dd79140c6ce5de2d6a6ef5e352",
                "command": "",
                "contractName": "contract_1545643033",
                "contractType": "fabric",
                "contractVersion": "v1.0",
                "operation": "install",
                "sequence": "0",
                "signature": "",
                "timestamp": 1545643033,
                "user": "user1",
                "vmVersion": "1.0"
            }
        ]

        cache.get_non_consumed_txns = MagicMock(side_effect=self.mock_get_non_consumed_txns)
        cache.set_txn_as_synced = MagicMock(side_effect=self.mock_set_txn_as_synced)
        tm.broadcast_tx_async = MagicMock(side_effect=self.mock_broadcast_tx_async)
        tm.broadcast_txs_async = MagicMock(side_effect=self.mock_broadcast_txs_async)

        current_tag = TagGenerator.get_current_tag()
        push_and_sync(current_tag)

        for ipfs_addr in self.ipfs_addrs:
            cache.set_txn_as_synced.assert_any_call(ipfs_addr, current_tag)

        tm.broadcast_tx_async.assert_any_call(json.dumps(self.transactions[0], sort_keys=True))

    def test_push_and_sync_for_more_contract(self):

        self.ipfs_addrs = [
            "QmT2fg3n4dZPcRismgRRp65icuTch21NCRRjYyPh5W3Kwp",
            "Qmd8o5YjC4cGqEG8sPAPvyNSeCWEFGNyAeQJnSfUcKyTiJ",
            "QmR3s8rsct34xXbq8PNK8PqVaCV55KHd53s9kmfKigocBh"
        ]
        self.transactions = [
            {"address": "QmStXVUdqAbKDTecFQkawuvSNqPs6su5KJB94Uvd9MiCny",
             "checkMD5": "7a3f59dd79140c6ce5de2d6a6ef5e352", "command": "", "contractName": "contract_1545726302",
             "contractType": "fabric", "contractVersion": "v1.0", "operation": "install", "sequence": "0",
             "signature": "", "timestamp": 1545726302, "type": "contract", "user": "user1", "vmVersion": "1.0"},
            {"address": "QmStXVUdqAbKDTecFQkawuvSNqPs6su5KJB94Uvd9MiCny",
             "checkMD5": "7a3f59dd79140c6ce5de2d6a6ef5e352", "command": "", "contractName": "contract_1545724345",
             "contractType": "fabric", "contractVersion": "v1.0", "operation": "install", "sequence": "0",
             "signature": "", "timestamp": 1545724345, "type": "contract", "user": "user1", "vmVersion": "1.0"},
            {"address": "QmStXVUdqAbKDTecFQkawuvSNqPs6su5KJB94Uvd9MiCny",
             "checkMD5": "7a3f59dd79140c6ce5de2d6a6ef5e352", "command": "", "contractName": "contract_1545724782",
             "contractType": "fabric", "contractVersion": "v1.0", "operation": "install", "sequence": "0",
             "signature": "", "timestamp": 1545724782, "type": "contract", "user": "user1", "vmVersion": "1.0"},
        ]

        cache.get_non_consumed_txns = MagicMock(side_effect=self.mock_get_non_consumed_txns)
        cache.set_txn_as_synced = MagicMock(side_effect=self.mock_set_txn_as_synced)
        tm.broadcast_tx_async = MagicMock(side_effect=self.mock_broadcast_tx_async)
        ipfs_client.get_json = MagicMock(side_effect=self.mock_get_json)

        current_tag = TagGenerator.get_current_tag()
        push_and_sync(current_tag)

        tm.broadcast_tx_async.assert_has_calls([
            call(json.dumps(self.transactions[1], sort_keys=True)),
            call(json.dumps(self.transactions[2], sort_keys=True)),
            call(json.dumps(self.transactions[0], sort_keys=True))
        ])

        cache.set_txn_as_synced.assert_has_calls([
            call(self.ipfs_addrs[1], current_tag),
            call(self.ipfs_addrs[2], current_tag),
            call(self.ipfs_addrs[0], current_tag)
        ])


