# Loppu: the scalable scientific robotics framework for getting research done

**L**aboratory **O**perations and **P**rocess **P**latform for
**U**ser-friendly robotics
---
Loppu is an open source laboratory robotics framework designed to simplify
the development and deployment of autonomous laboratory systems. Loppu
provides a robust core library with a simple API for quickly automating
everyday lab tasks, and is written to scale as your needs.

## Features

* Easy installation and use on any platform
* Example robots you can download to start automating your lab today, and
  reference when you want to learn more
* Designed for reuse -- build robots without learning how to code
* Designed to teach and scale -- learn how to code and build unique robots
* Robust core library
    * Data streaming
    * Data storage
    * Plotting
    * Watchdogs

### Designed around the needs of modern laboratories

The common theme of every lab is generating data, so Loppu provides tools to
efficiently and securely generate and handle data.

### Designed to write as little code as possible

Loppu is designed to allow scientists to automate their labwork without
needing to become robotics experts. This is done by focusing on a robust
core library of `nodes` and allowing users to share nodes.

Loppu provides a feature-rich cli, `scaffold`, to allow users to interact
with the framework.

## Installation

The framework can be installed using `go get`:

```zsh
go get github.com/jt05610/loppu
```

## Usage

The following code snippet demonstrates how to use the framework to create a
syringe pump:

```go
package main

import (
	"github.com/jt05610/loppu"
)

func main() {
}
```

## Contributing

Contributions are welcome! Please refer to
the [contributing guidelines](CONTRIBUTING.md) for more information.

## Roadmap

| Feature                     | Target Date |
|-----------------------------|-------------|
| Arduino implementation      | May 2023    |
| Syringe pump example        | May 2023    |
| PipBot example              | May 2023    |
| Python, R, Julia clients    | May 2023    |
| Raspberry Pi implementation | Jul 2023    |
| ROS integration             | Jul 2023    |

### Feature Details

#### Arduino implementation

##### Description

Implement [loppu-firmware](https://github.com/jt05610/loppu-firmware) on
Arduino as this is the most popular open-source microcontroller.

##### Goals

* Implement all relevant libraries
* Get syringe pump example working with Arduino

#### Syringe pump example

##### Description

Create clear examples of how to use the framework through implementation of
syringe pumps useful for microfluidic experiments.

##### Goals

* Syringe pump example w/gui and selectable syringes
* Scale syringe pumps and include calibration

#### PipBot example

##### Description

Code needed to hack a 3D printer and operate as a fluid handler.

##### Goals

* Make PipBot easy to download and use
* Include Excel template for designing and analyzing experiments
* Also have web interface to help users learn to transition from Excel for
  data analysis :)

#### Python, R, Julia clients

##### Description

Generate clients in popular languages to interact with robots

##### Goals

* Python client
* R client
* Julia client

#### Raspberry Pi implementation

##### Description

Since Raspberry Pi is also popular and would be very useful for this framework,
the framework needs to be implemented on it.

##### Goals

* Full Hardware and Software implementation on same device
    * USB port to RS485 for extending
* Support PiCamera

#### ROS integration

##### Description

As ROS is the most popular open-source robotics framework, a ROS node should
be published to make it easy to integrate existing robots.

##### Goals

* ROS node
* Implement one syringe pump with ROS and another with Loppu, then make them
  work together

## FAQs

### What does "Loppu" mean?

Loppu is the Finnish word for "end" or "finish" which reflects our
primary goal of providing tools for finishing research projects faster.
Chat-GPT was used to generate the backronym.

### Why do we need another robotics framework?

The needs of a robot operating in a predictable laboratory environment are
simply different from the needs of one in traditional robotics settings.
