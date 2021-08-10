<img src="http://www.cineoutsider.com/reviews/pix/z/za/zardoz1.jpg" width="100%" alt="I am ZARDOZ">

# ZARDOZ
(A Golang testing framework)[https://github.com/APiercey/zardoz] for testing asyncronous dependencies.

Table of Contents
=================

* [Installation](#installation)
* [Usage](#usage)
   * [Test](#test)
   * [Assertions](#assertions)
      * [t.Assert](#tassert)
      * [t.AssertSync](#tassertsync)
      * [t.AssertAsync](#tassertasync)
   * [Working Example](#working-example)
* [Reading results](#reading-results)
* [Setup and Cleanup](#setup-and-cleanup)

## Installation

Pull the module into your GOPATH using
```golang
go get github.com/apiercey/zardoz
```

## Usage
ZARDOZ can be imported like so:
```golang
import z "github.com/APiercey/zardoz"
```
ZARDOZ comes with [three assertions](https://github.com/APiercey/zardoz/tree/main#assertions) for testing asyncronous dependencies. Test blocks begin with a `Describe` block, which provides a test suite.

```golang
func main() {
  z.Describe("Example External System", func(s *z.Suite) {
      s.Test("can turn lights on", test_light_turn_on)
  })
}
```

### Test

A test can be written using anonymous functions or by passing in the function name to the describe block. The function will receive a Test struct, which provides assertions.

Examples:

```golang
func test_lego_car_turns_left(t *z.Test) {
    execute_turn_left_command()

    t.AssertSync(func() bool {
        return evaluate_lego_car_turned_left()
    }, 500)
}

s.Test("can turn left", test_lego_car_turns_left)
```

or simply

```golang
s.Test("can turn left", func (t *z.Test) {
    execute_turn_left_command()

    t.AssertSync(func() bool {
        return evaluate_lego_car_turned_left()
    }, 500)
})

```

### Assertions
ZARDOZ provides three assertions:

- `t.Assert` asserts a value is true immediatly.
- `t.AssertSync` assert a condition will become true.
- `t.AssertAsync` asserts a condition will become true and can be parallelized with other Assert calls.

#### t.Assert
Simple assertions like in other unit testing frameworks.

Examples
```golang
t.Assert(true) // passes
t.Assert(false) // fails
t.Assert(some_func_that_returns_boolean()) // depends on the returned value
```

#### t.AssertSync
Asserts something will eventually true. It expects two values: a predict function that compares for a condition and a timeout of how long it should attempt to assert this condition. Timeout values are in milliseconds.

If the predict should return _false_ it will try again.

Multiple `t.AssertSync` assertions result will assert their conditions syncronously.

Examples
```golang
t.AssertSync(func() {
  return true
}, 1_000) 

t.AssertSync(func() {
  return maybe_returns_true()
}, 1_000) 

t.AssertSync(maybe_returns_true, 1_000) 
```

In the example above, the test will take a maximum time of 3000ms to complete if all assertions are false.

#### t.AssertAsync
Asserts something will eventually be true. It expects two values: a predict function that compares for a condition and a timeout of how long it should attempt to assert this condition. Timeout values are in milliseconds.

If the predict should return _false_ it will try again.

Multiple `t.AssertSync` assertions result will assert their conditions asyncronously and can be used to observe conditions which take place within the same time frame.

Examples
```golang
t.AssertAsync(func() {
  return true
}, 1_000) 

t.AssertAsync(func() {
  return maybe_returns_true()
}, 1_000) 

t.AssertAsync(maybe_returns_true, 1_000) 
```

In the example above, the test will take a maximum time of roughly 1000ms to complete if all assertions are false.


### Working Example
Below is an example of testing IoT device using the MQTT protocol. The suite will send messaeges to a device and expects an asynchronous response.

```golang
package main

import "fmt"
import z "github.com/APiercey/zardoz"
import mqtt "github.com/eclipse/paho.mqtt.golang"

// Setup code for completness
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
    fmt.Printf("Unexpected connection lost: %v", err)
}

// Setup code for completness
func setupMqttClient(clientID string) mqtt.Client {
    var broker = "localhost"
    var port = 1883

    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
    opts.SetClientID("test_runner")
    opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
    })
    opts.OnConnectionLost = connectLostHandler
    client := mqtt.NewClient(opts)

    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    return client
}

func test_falls_asleep_on_message(t *z.Test) {
    expect_to_fall_asleep := false

    client.Subscribe("iot-device/output", 1, func(c mqtt.Client, m mqtt.Message) {
        if string(m.Payload()) == "ASLEEP" {
          expect_to_fall_asleep = true
        }
    }).Wait()

    client.Publish("iot-device/input", 0, false, "SLEEP")

    t.AssertSync(func() bool { return expect_to_fall_asleep }, 500)
}

func test_wakes_up_on_messages(t *z.Test) {
    expect_to_be_awake := false
    expect_to_download_config := false

    client.Subscribe("iot-device/output", 1, func(c mqtt.Client, m mqtt.Message) {
        if string(m.Payload()) == "WOKEN_UP" {
          expect_to_be_awake = true
        }
        
        if string(m.Payload()) == "CONFIG_DOWNLOADED" {
          expect_to_download_config = true
        }
    }).Wait()

    client.Publish("iot-device/input", 0, false, "WAKEUP")

    // Both assertions should be true within 500ms
    t.AssertAsync(func() bool { return expect_to_be_awake }, 500)
    t.AssertAsync(func() bool { return expect_to_download_config }, 500)
}

func test_receives_ping(t *z.Test) {
    expected_response_received := false

    client.Subscribe("iot-device/output", 1, func(c mqtt.Client, m mqtt.Message) {
        expected_response_received = string(m.Payload()) == "PING"
    }).Wait()

    t.AssertSync(func() bool {
        return expected_response_received
    }, 1_000)
}

func main() {
    z.Describe("IoT Device Commands", func(s *z.Suite) {
        s.Test("wakes up on message", test_wakes_up_on_messages)
        s.Test("falls asleep on message", test_falls_asleep_on_message)
    })

    z.Describe("Diagnostics", func(s *z.Suite) {
        s.Test("receives PING", test_receives_ping)
    })
}
```

## Reading results
After every test ZARDOZ will provide a summary of the number of tests ran, passes, and failes.
When tests fail, ZARDOZ will provide a preview of the offending assertion and as well as the path and line number under the summary.

Example below:
```golang
$ go run *.go

Running...
FF

Failures:

  1) Example External System can turn on

    Never returned true after 500ms when evaluating:
    t.AssertSync(func() bool {
        return expected_response_received
    }, 500)
}
  2) Example External System can turn off

    Never returned true after 500ms when evaluating:
    t.AssertSync(func() bool {
        return expected_response_received
    }, 500)
}

Assertion failed /Users/example-user/.go/src/example-tests/v1/main.go:18
Assertion failed /Users/example-user/.go/src/example-tests/v1/main.go:32

2 tests ran with 0 passes and 2 failures.

Running...
F

Failures:

  1) Information about external system

    Never returned true after 15000ms when evaluating:
    t.AssertSync(func() bool {
        return expected_response_received
    }, 15_000)
}

Assertion failed /Users/example-user/.go/src/example-tests/v1/main.go:44

1 tests ran with 0 passes and 1 failures.
```

## Setup and Cleanup
ZARDOZ allows you to run setup code or cleanup code, if you so wish. These are executed once per test.

- `s.Setup` executed before the test.
- `s.Cleanup` executed after the test.

Examples
```golang
func main() {
  z.Describe("Example External System", func(s *z.Suite) {
      s.Setup(func () {
          setup_my_dependency()
      })
      
      // or
      
      s.Setup(setup_my_dependency)
      
      s.Test("can handle commands", test_some_command)
      
      s.Cleanup(func () {
          cleanup_my_dependency()
      })
      
      // or
      
      s.Cleanup(cleanup_my_dependency)
      
  })
}
```
