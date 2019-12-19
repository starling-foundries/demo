const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { StatusType, MessageType } = require('@zilliqa-js/subscriptions');
const { Contracts, ContractStatus } = require('@zilliqa-js/contract')
async function test() {
  // first run `python websocket.py` to start websocket server locally
  const zilliqa = new Zilliqa('https://dev-api.zilliqa.com');
  const contract = Contract.at("zil15e20r8mz6zwqxa7mvg2a72pazvdevcuguafxfp").addresses()
  console.log(contract);
   const subscriber = zilliqa.subscriptionBuilder.buildEventLogSubscriptions(
    'ws://',
    {
      addresses: [
        contract,
      ],
    },
  );
  subscriber.emitter.on(StatusType.SUBSCRIBE_EVENT_LOG, (event) => {
    console.log('get SubscribeEventLog echo: ', event);
  });
  subscriber.emitter.on(MessageType.EVENT_LOG, (event) => {
    console.log('get new event log: ', JSON.stringify(event));
  });
  subscriber.emitter.on(MessageType.UNSUBSCRIBE, (event) => {
    console.log('get unsubscribe event: ', event);
  });

  await subscriber.start();
  // await subscriber.stop();
}

test();