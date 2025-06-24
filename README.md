# protohackers

[Protohackers](https://protohackers.com/) server programming challenge answers.

## Deployment

- install ['flyctl'](https://fly.io/docs/flyctl/);
- authenticate;
- navigate to the root of the repository;
- allocate dedicated IPv4 address (costs $2/mo): `flyctl ips allocate-v4 --app protohkr`;
  - without this it is not possible to connect over just TCP, even if an IPv6 address is allocated;
  - more information [here](https://community.fly.io/t/tcp-and-udp-service-ports-dont-work/9746) and [here](https://community.fly.io/t/announcement-shared-anycast-ipv4/9384/25);
- deploy: `flyctl deploy --dockerfile Dockerfile.<problem name> --ha=false`;
  - see problem names under `cmd/` directory and [here](https://protohackers.com/problems).

> [!important]
> Any uncommitted changes will also be deployed.
