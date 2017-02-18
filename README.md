`lowed` is a small service for generating a never-ending stream of random metric samples for a configured list of service names and metric types. It specifically generates [DogStatsD](http://docs.datadoghq.com/guides/dogstatsd/#datagram-format) metrics suitable for consumption by [Veneur](https://github.com/stripe/veneur).

It generates an instance of each `metrics` for each of the configured `services` every `delay`. Non-counter metrics support a `range` with a `max` and `min`, from which a random value will be generated.

# TODO

Tags!

# Usage

```
./lowed -f example
```

# Example Config

```yaml
delay: "100ms"
stats_address: 127.0.0.1:8200
services:
  - stallion_srv
  - okapi_srv
  - ermine_srv
  - wombat_srv
  - warthog_srv
  - burro_srv
metrics:
  counters:
    - name: request_errors_total
    - name: request_successes_total
  timers:
    - name: request_duration_ms
      range:
        min: 10
        max: 300
  histograms:
    - name: request_body_bytes
      range:
        min: 1024
        max: 307200
  gauges:
    - name: request_queue_depth
      range:
        min: 1
        max: 20
```

# Name

Lowed is named for [Paddy Lowe], codesigner of the [Mercedes F1 W05 Hybrid](https://en.wikipedia.org/wiki/Mercedes_F1_W05_Hybrid). Since `lowed` is an engine for generating metrics, it made sense to choose someone who designed and engine! Also, it is a clever homonym of "load", since the generated metrics may provide load testing fodder.
