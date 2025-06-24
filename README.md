# protohackers

[Protohackers](https://protohackers.com/) server programming challenge answers.

## Deployment

- install ['flyctl'](https://fly.io/docs/flyctl/);
- authenticate;
- navigate to the root of the repository;
- allocate dedicated IPv4 address (costs $2/mo): `flyctl ips allocate-v4 --app <app name>`;
  - without this it is not possible to connect over just TCP, even if an IPv6 address is allocated;
  - [here's](https://www.tigrisdata.com/blog/docker-registry-at-home/) a blog post over at [Tigris](https://tigrisdata.com) that gives away $50 in credits;
  - more information about this requirement [here](https://community.fly.io/t/tcp-and-udp-service-ports-dont-work/9746) and [here](https://community.fly.io/t/announcement-shared-anycast-ipv4/9384/25);
- deploy: `flyctl deploy --dockerfile Dockerfile.<problem name> --ha=false`;
  - see problem names under `cmd/` directory and [here](https://protohackers.com/problems).

> [!important]
> Any uncommitted changes will also be deployed.

## Testing

### Locally

- move to problem directory: `pushd cmd/<problem name>/`
- run tests:

  ```console
  go test -v
  === RUN   TestProblem
  --- PASS: TestProblem (0.00s)
  PASS
  ok      github.com/usrme/protohackers/cmd/<problem name>    3.138s
  ```

- move back to the root of the repository: `popd`;
- build image: `docker build --tag <problem name> -f Dockerfile.<problem name> .`;
- run container: `docker run -p 5000:5000 <problem name>`;
- pass command:

  ```console
  $ echo "hello world" | nc -w1 localhost 5000
  hello world
  ```

### Remotely

- [deploy](#deployment);
- pass command:

  ```console
  $ echo "hello world" | nc -w1 <app name>.fly.dev 5000
  hello world
  ```
