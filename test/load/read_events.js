import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

export let options = {
  summaryTrendStats: ["avg", "min", "max", "med", "p(75)", "p(99)"],
};

const client = new grpc.Client();
client.load(['../../internal/event/api'], 'event.proto');

export default () => {
  client.connect('localhost:5050', {
    plaintext: true
  });

  const data = {id: 1};
  const response = client.invoke('event.EventService/GetEvent', data);

  check(response, {
    'status is OK': (r) => r && r.status === grpc.StatusOK,
  });

  client.close();
  sleep(1);
};
