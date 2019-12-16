from pprint import pprint
from pyzil.crypto import zilkey
from pyzil.zilliqa import chain
from pyzil.zilliqa.units import Zil, Qa
from pyzil.account import Account, BatchTransfer
chain.set_active_chain(chain.TestNet)  
Account.from_keystore("zilliqa_keystore.json")
Account.from_keystore(keystore_file="zilliqa_keystore.json", password="password")
Acct= Account.from_keystore(keystore_file="zilliqa_keystore.json", password="password")
print(Acct.get_balance)
print(Acct.get_balance())
code = open("FungibleToken.scilla").read()

from pyzil.contract import Contract
contract = Contract.new_
contract = Contract.new_from_code(code)
contract.account = Acct
contract.deploy(timeout=300, sleep=10, gas_limit=(10000*100))

