from pprint import pprint

from pyzil.crypto import zilkey
from pyzil.zilliqa import chain
from pyzil.zilliqa.units import Zil, Qa
from pyzil.account import Account
from pyzil.contract import Contract


chain.set_active_chain(chain.TestNet)

contract_addr="zil15e20r8mz6zwqxa7mvg2a72pazvdevcuguafxfp"
contract = Contract.load_from_address(contract_addr, load_state=True)
print(contract.status)
pprint(contract.state)
owner = Account(private_key="3375F915F3F9AE35E6B301B7670F53AD1A5BE15D8221EC7FD5E503F21D3450C8")
contract.account = owner
pprint(owner.address0x)
pprint(owner)
# response = contract.call(
    # method="AuthorizeOperator", 
    # params=[
        # Contract.value_dict(
            # "operator",
            # "ByStr20",
            # "zil1sf2t9jdvmuvp6ht8jmtrxg8mkgx5ahgj6h833r")
        # ])

# Try to mint new tokens
response = contract.call(method="Mint", 
        params=[
            Contract.value_dict(
                "recipient",
                "ByStr20",
                "zil1sf2t9jdvmuvp6ht8jmtrxg8mkgx5ahgj6h"), 
            Contract.value_dict(
                "amount",
                "Uint128",
                100000)
            ])

pprint(response)

# Now try to send tokens to user with no ZIL
poor_user = Account()

    response = contract.call(method="Transfer", 
